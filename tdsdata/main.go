package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	hc := http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxIdleConns:        8,
			MaxIdleConnsPerHost: 8,
			ReadBufferSize:      8 * 1024 * 1024,
		},
	}

	apiKey := os.Getenv("API_KEY")

	scraper := ScheduleScraper{
		client:  hc,
		apiKey:  apiKey,
		baseUrl: "https://ride-api.bustickets.com/tickets",
	}

	schedules, err := scraper.Scrape()
	if err != nil {
		return
	}

	for _, s := range schedules {
		fmt.Println(s)
	}
}
