package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func crawlPage(domain, rawCurrentURL string, pages map[string]int) {
	urlParse, _ := url.Parse(rawCurrentURL)
	if urlParse.Host != domain {
		return
	}
	normalizedUrl, err := NormalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("Error normalizing %v", err)
		return
	}

	pages[normalizedUrl]++
	if pages[normalizedUrl] > 1 {
		return
	}
	htmlFromPage, err := GetHTML(rawCurrentURL)
	if err != nil {
		log.Printf("Error getting html: %v", err)
		return
	}
	fmt.Println(htmlFromPage)
	urlsCrawled, err := GetURLsFromHTML(htmlFromPage, rawCurrentURL)
	if err != nil {
		log.Printf("Error crawling url: %v", err)
		return
	}

	for _, urlCrawled := range urlsCrawled {
		crawlPage(domain, urlCrawled, pages)
	}
}

func CrawlPage(baseUrl string) error {
	urlParsed, err := url.Parse(baseUrl)
	if err != nil {
		return fmt.Errorf("error parsing url: %v", err)
	}
	pagesMap := make(map[string]int)

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, syscall.SIGINT)
	go func() {
		<-osSigChan
		fmt.Println()
		printMap(pagesMap)
		os.Exit(0)
	}()

	defer printMap(pagesMap)

	crawlPage(urlParsed.Host, baseUrl, pagesMap)
	return nil
}

func printMap(pagesMap map[string]int) {
	for k, v := range pagesMap {
		fmt.Printf("%s: %d\n", k, v)
	}
}
