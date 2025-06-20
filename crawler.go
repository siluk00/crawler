package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func (cfg *config) crawlPageRecursively(rawCurrentURL string, level int) {
	fmt.Printf("trying to crawl %s\n", rawCurrentURL)
	currentUrlParsed, err := url.Parse(rawCurrentURL)
	//fmt.Printf("*trying to crawl %s\n", currentUrlParsed.String())
	if err != nil {
		return
	}

	domainBaseUrl, _ := extractMainDomain(cfg.baseUrl.Host)
	domainCurrentUrl, _ := extractMainDomain(currentUrlParsed.Host)

	if domainCurrentUrl != domainBaseUrl && domainBaseUrl != "" {
		return
	}

	normalizedUrl, err := NormalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error normalizing %v", err)
		return
	}

	//cfg.mu.Lock()
	//defer cfg.mu.Unlock()

	cfg.pages[normalizedUrl]++
	if cfg.pages[normalizedUrl] > 1 {
		return
	}

	htmlFromPage, err := GetHTML(currentUrlParsed.Scheme + "://" + normalizedUrl)
	if err != nil {
		fmt.Printf("Error getting html: %v", err)
		return
	}
	//fmt.Println(htmlFromPage)

	urlsCrawled, err := cfg.GetURLsFromHTML(htmlFromPage, currentUrlParsed.Scheme+"://"+normalizedUrl)
	if err != nil {
		fmt.Printf("Error crawling url: %v", err)
		return
	}
	//fmt.Println(urlsCrawled)

	fmt.Printf("crawling %s in level %d\n", rawCurrentURL, level)
	level++
	for _, urlCrawled := range urlsCrawled {
		fmt.Printf("--%s\n", urlCrawled)
		cfg.crawlPageRecursively(urlCrawled, level)
	}
}

func (cfg *config) crawlPage() error {

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, syscall.SIGINT)
	go func() {
		<-osSigChan
		fmt.Println()
		printMap(cfg.pages)
		os.Exit(0)
	}()

	defer printMap(cfg.pages)

	//fmt.Printf("baseurl:%s\n", cfg.baseUrl.String())

	cfg.crawlPageRecursively(cfg.baseUrl.String(), 0)
	return nil
}

func printMap(pagesMap map[string]int) {
	for k, v := range pagesMap {
		fmt.Printf("%s: %d\n", k, v)
	}
}
