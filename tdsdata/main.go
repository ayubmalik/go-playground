package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	hc := http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConns:        8,
			MaxIdleConnsPerHost: 8,
			ReadBufferSize:      8 * 1024 * 1024,
		},
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("API_KEY environment variable is required")
		return
	}

	scraper := ScheduleScraper{
		client:  hc,
		apiKey:  apiKey,
		baseUrl: "https://ride-api.bustickets.com/tickets",
	}

	schedules, err := scraper.Scrape()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, s := range schedules {
		fmt.Println(s)
	}
}
