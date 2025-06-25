package main

import (
	"context"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]bool
	baseUrl            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	counter            int
	maxPages           int
	ctx                context.Context
	cancel             context.CancelFunc
}
