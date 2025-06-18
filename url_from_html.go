package main

import (
	"strings"

	"golang.org/x/net/html"
)

func GetURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	node, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}
	links := traverseNode(node, 0)
	for i, link := range links {
		if link[0] == '/' {
			links[i] = rawBaseURL + link
		}
	}

	return links, nil
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
