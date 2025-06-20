package main

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {

	cfg := config{

		pages:              nil,
		mu:                 nil,
		concurrencyControl: nil,
		wg:                 nil,
	}

	tests := []struct {
		url      string
		htmlBody string
		expected []string
	}{
		// Test 1: Basic relative and absolute links
		{
			url: "https://blog.boot.dev",
			htmlBody: `
            <html>
                <body>
                    <a href="/path/one">Relative</a>
                    <a href="https://other.com/path/one">Absolute</a>
                    <a href="path/two">Relative no slash</a>
                </body>
            </html>`,
			expected: []string{
				"https://blog.boot.dev/path/one",
				"https://other.com/path/one",
				"https://blog.boot.dev/path/two",
			},
		},

		// Test 2: Edge cases with malformed URLs and fragments
		{
			url: "https://example.com/base/",
			htmlBody: `
            <html>
                <body>
                    <a href="#section">Fragment</a>
                    <a href="?query=param">Query only</a>
                    <a href="mailto:test@example.com">Mailto</a>
                    <a href="javascript:void(0)">JS link</a>
                    <a href="//protocol-relative.com">Protocol relative</a>
                </body>
            </html>`,
			expected: []string{
				"https://example.com/base/?query=param",
				"https://protocol-relative.com/",
			},
		},

		// Test 3: Nested elements and whitespace
		{
			url: "https://test.org",
			htmlBody: `
            <html>
                <body>
                    <div>
                        <a href="/nested">
                            <span>With nested</span>
                            <img src="test.png">
                        </a>
                    </div>
                    <a href="  /trim  ">With whitespace</a>
                </body>
            </html>`,
			expected: []string{
				"https://test.org/nested",
				"https://test.org/trim",
			},
		},

		// Test 4: Base tag and relative URLs
		{
			url: "https://base.com/original/",
			htmlBody: `
            <html>
                <head>
                    <base href="https://base.com/override/">
                </head>
                <body>
                    <a href="path/three">Base relative</a>
                    <a href="/path/four">Root relative</a>
                </body>
            </html>`,
			expected: []string{
				"https://base.com/override/path/three",
				"https://base.com/path/four",
			},
		},

		// Test 5: Multiple links with duplicates
		{
			url: "https://dupe.check",
			htmlBody: `
            <html>
                <body>
                    <a href="/same">Link 1</a>
                    <a href="/same">Link 2</a>
                    <a href="/other">Link 3</a>
                    <a href="/other#frag">Link 4</a>
                </body>
            </html>`,
			expected: []string{
				"https://dupe.check/same",
				"https://dupe.check/other",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			var err error
			cfg.baseUrl, err = url.Parse(test.url)
			if err != nil {
				t.Errorf("error parsing baseurl: %v", err)
			}
			actual, err := cfg.GetURLsFromHTML(test.htmlBody, test.url)
			if err != nil {
				t.Errorf("error getting urls: %v", err)
			}
			if reflect.DeepEqual(actual, test.expected) {
				t.Errorf("Test %d failed: Expected: %s, Actual: %s", i, test.expected[0], actual[0])
			}
		})
	}
}
