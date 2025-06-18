package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetHTML(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error getting html response: %v", err)
	}
	defer response.Body.Close()

	contentType := response.Header.Get("content-type")
	if !strings.Contains(contentType, "text/html") && !strings.Contains(contentType, "application/xhtml+xml") {
		return "", fmt.Errorf("response is not of the html type")
	}

	if response.StatusCode >= 400 {
		return "", fmt.Errorf("response status is an error message")
	}

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	return string(respBody), nil
}
