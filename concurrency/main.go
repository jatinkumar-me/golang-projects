package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Don't communicate by sharing memory; instead share memory by communicating.
type result struct {
	url     string
	err     error
	latency time.Duration
}

// It takes a url as a string and a channel with only write end. Won't be able to read any data
func get(ctx context.Context, url string, ch chan<- result) {
	start := time.Now()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if resp, err := http.DefaultClient.Do(req); err != nil {
		ch <- result{url, err, 0}
	} else {
		t := time.Since(start).Round(time.Millisecond)
		time.Sleep(100 * time.Millisecond)
		ch <- result{url, nil, t}
		resp.Body.Close()
	}
}

func main() {
	resultChan := make(chan result)
	list := []string{
		"https://jsonplaceholder.com",
		"https://youtube.com",
		"https://www.wsj.com",
		"https://www.facebook.com",
		"https://www.google.com",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Millisecond)
	defer cancel()

	for _, url := range list {
		go get(ctx, url, resultChan)
	}

	for range list {
		r := <-resultChan

		if r.err != nil {
			log.Printf("%-20s %s\n", r.url, r.err)
		} else {
			log.Printf("%-20s %s\n", r.url, r.latency)
		}
	}
}
