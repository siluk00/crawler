package main

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// ainda falta testar o caso em que duas queries são iguais, tipo item=a&item=b
func NormalizeURL(urlUnparsed string) (string, error) {
	urlParsed, err := url.Parse(urlUnparsed)
	if err != nil {
		return "", err
	}

	hostname := urlParsed.Hostname()
	if urlParsed.IsAbs() && urlParsed.Scheme != "http" && urlParsed.Scheme != "https" {
		return "", fmt.Errorf("scheme is not http")
	}
	if !strings.Contains(urlUnparsed, "://") || !strings.Contains(hostname, ".") {
		return "", fmt.Errorf("invalid url")
	}
	if urlParsed.RawQuery == "" {
		return strings.TrimRight(urlParsed.Hostname()+urlParsed.Path, "/"), nil
	}

	queries := urlParsed.Query()
	if len(queries) == 1 {
		return hostname + urlParsed.Path + "?" + urlParsed.RawQuery, nil
	}

	keys := make([]string, 0, len(queries))
	for k := range queries {
		keys = append(keys, k)
	}

	finalQuery := ""
	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			finalQuery += "&"
		}

		finalQuery += k + "=" + queries[k][0]
	}

	return hostname + urlParsed.Path + "?" + finalQuery, nil
}
