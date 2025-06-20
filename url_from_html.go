package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

func (cfg *config) GetURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	node, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	links := traverseNode(node, 0)
	linksToReturn := make([]string, 0, 1)

	for _, link := range links {
		//fmt.Println("**", link)
		if len(link) > 1 && link[:2] == "//" {
			linksToReturn = append(linksToReturn, parsedURL.Scheme+":"+link)
		} else if link[0] == '/' {
			linksToReturn = append(linksToReturn, cfg.baseUrl.String()+link)
		} else if link[0] == '#' {
			continue
		}
	}

	return linksToReturn, nil
}

func traverseNode(node *html.Node, depth int) []string {
	links := make([]string, 0)
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, traverseNode(c, depth+1)...)
	}

	return links
}

func extractMainDomain(hostname string) (string, error) {
	mainhost, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return "", err
	}

	return strings.Split(mainhost, ".")[0], nil
}
