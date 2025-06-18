package main

import (
	"fmt"
	"os"
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
	fmt.Printf("starting crawl of %s\n", website)
	go CrawlPage(website)
	time.Sleep(5 * time.Second)
	os.Exit(0)
}
