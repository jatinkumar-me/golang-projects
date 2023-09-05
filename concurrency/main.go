package main

import (
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
func get(url string, ch chan<- result) {
	start := time.Now()

	if resp, err := http.Get(url); err != nil {
		ch <- result{url, err, 0}
	} else {
		t := time.Since(start).Round(time.Millisecond)
        time.Sleep(100*time.Millisecond)
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

	for _, url := range list {
		go get(url, resultChan)
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
