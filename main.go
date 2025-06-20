package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	website := os.Args[1]
	websiteURL, err := url.Parse(website)
	if err != nil {
		log.Fatalf("provide a url")
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	cfg := config{
		baseUrl:            websiteURL,
		pages:              make(map[string]int),
		mu:                 &mu,
		concurrencyControl: make(chan struct{}),
		wg:                 &wg,
	}
	fmt.Printf("starting crawl of %s\n", website)
	go cfg.crawlPage()
	time.Sleep(5 * time.Second)
	fmt.Printf("Ending after %d seconds\n", 5)
	os.Exit(0)
}
