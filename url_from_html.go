package main

import (
	"net"
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

	baseDomain, err := extractMainDomain(parsedURL.Host)
	if err != nil {
		return nil, err
	}

	links := traverseNode(node)
	linksToReturn := make([]string, 0, len(links))

	for _, link := range links {
		// Skip empty links and fragments
		if len(link) == 0 || link[0] == '#' {
			continue
		}

		var finalURL string

		// Handle protocol-relative URLs (//)
		if len(link) > 1 && link[:2] == "//" {
			finalURL = parsedURL.Scheme + ":" + link
		} else if link[0] == '/' {
			finalURL = cfg.baseUrl.String() + link
		} else {
			finalURL = link
		}

		// Parse to verify and check domain
		parsedLink, err := url.Parse(finalURL)
		if err != nil {
			continue // Skip invalid URLs
		}

		// Skip if no host (invalid URL)
		if parsedLink.Host == "" {
			continue
		}

		// Extract and compare main domains
		linkDomain, err := extractMainDomain(parsedLink.Host)
		if err != nil || linkDomain != baseDomain {
			continue // Skip different domains/subdomains
		}

		linksToReturn = append(linksToReturn, finalURL)
	}

	return linksToReturn, nil
}

func extractMainDomain(hostname string) (string, error) {
	// Handle IP addresses
	if ip := net.ParseIP(hostname); ip != nil {
		return hostname, nil
	}

	// Remove port if present
	hostname = strings.Split(hostname, ":")[0]

	// Get main domain using publicsuffix
	domain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		// Fallback for non-standard domains
		parts := strings.Split(hostname, ".")
		if len(parts) > 1 {
			return parts[len(parts)-2], nil
		}
		return hostname, nil
	}

	// Return just the main part (e.g., "example" from "example.com")
	return strings.Split(domain, ".")[0], nil
}

func traverseNode(node *html.Node) []string {
	links := make([]string, 0)
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, traverseNode(c)...)
	}

	return links
}
