package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func (cfg *config) crawlPageRecursively(rawCurrentURL string, level int) {
	defer cfg.wg.Done()

	cfg.concurrencyControl <- struct{}{}
	defer func() { <-cfg.concurrencyControl }()

	fmt.Printf("beggining crawl of %s at level %d\n", rawCurrentURL, level)
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
		log.Printf("error normalizing: %v\n", err)
		return
	}

	htmlFromPage, err := GetHTML(currentUrlParsed.Scheme + "://" + normalizedUrl)
	if err != nil {
		log.Printf("Error getting html: %v", err)
		return
	}
	//fmt.Println(htmlFromPage)

	urlsCrawled, err := cfg.GetURLsFromHTML(htmlFromPage, currentUrlParsed.Scheme+"://"+normalizedUrl)
	if err != nil {
		log.Printf("Error crawling url: %v", err)
		return
	}
	//fmt.Println(urlsCrawled)

	fmt.Printf("crawling %s in level %d\n", rawCurrentURL, level)
	level++
	for _, urlForCrawl := range urlsCrawled {
		//fmt.Printf("--%s\n", urlCrawled)

		norm, err := NormalizeURL(urlForCrawl)
		if err != nil {
			log.Printf("error normalizing url: %v\n", err)
			continue
		}
		//fmt.Printf("%d:%d ", cfg.counter, cfg.maxPages)

		cfg.mu.Lock()

		cfg.counter++
		if cfg.counter > cfg.maxPages {
			cfg.cancel()
		}

		if !cfg.pages[norm] {
			cfg.pages[norm] = true
			cfg.mu.Unlock()

			cfg.wg.Add(1)

			go cfg.crawlPageRecursively(urlForCrawl, level)
		} else {
			cfg.mu.Unlock()
		}

		select {
		case <-cfg.ctx.Done():
			return
		default:
		}

	}

}

func (cfg *config) crawlPage() {
	defer cfg.cancel()

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, syscall.SIGINT)
	done := make(chan struct{})

	go func() {
		<-osSigChan
		fmt.Println("interrupting...")
		close(done)
		cfg.wg.Wait()
		cfg.printMap()
		os.Exit(0)
	}()

	//fmt.Printf("baseurl:%s\n", cfg.baseUrl.String())
	cfg.wg.Add(1)

	go cfg.crawlPageRecursively(cfg.baseUrl.String(), 0)
	select {
	case <-done:
	case <-cfg.ctx.Done():
		fmt.Println("maxPages reached")
		return
	default:
		cfg.wg.Wait()
		cfg.printMap()
	}

	if err := cfg.ctx.Err(); err != nil {
		fmt.Printf("Crawler finished due to context cancellation: %v\n", err)
	} else {
		fmt.Println("Crawler finished naturally.")
	}
}

func (cfg *config) printMap() {
	//for k, v := range cfg.pages {
	//	fmt.Printf("%s: %d\n", k, v)
	//}
	fmt.Printf("%d\n", cfg.counter)
}
