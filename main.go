package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	const timeOfExecution = 30
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	websiteURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatalf("provide a url")
	}
	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("provide a number as concurrency")
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("provide a number for max pages")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	cfg := config{
		baseUrl:            websiteURL,
		pages:              make(map[string]bool),
		mu:                 &mu,
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &wg,
		counter:            0,
		maxPages:           maxPages,
	}
	cfg.ctx, cfg.cancel = context.WithCancel(context.Background())
	fmt.Printf("starting crawl of %s\n", os.Args[1])
	go cfg.crawlPage()
	time.Sleep(timeOfExecution * time.Second)
	fmt.Printf("Ended after %d seconds\n", timeOfExecution)
}
