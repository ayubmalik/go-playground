package main

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

type City struct {
	CityId int    `json:"cityId"`
	Name   string `json:"name"`
}

type Stop struct {
	StopUuid    string `json:"stopUuid"`
	StationName string `json:"stationName"`
	StationCode string `json:"stationCode"`
	City        City   `json:"city"`
}

type ODPair struct {
	origin      Stop
	destination Stop
	key         string
}

type ScheduleScraper struct {
	client  http.Client
	apiKey  string
	baseUrl string
}

func (s ScheduleScraper) Scrape() ([]ODCandidate, error) {
	stops, err := s.getStops()
	if err != nil {
		return nil, fmt.Errorf("error getting stops: %w", err)
	}

	var pairs []ODPair
	for _, orig := range stops {
		for _, dest := range stops {
			var key string
			if strings.Compare(orig.StationCode, dest.StationCode) < 0 {
				key = orig.StationCode + dest.StationCode
			} else {
				key = dest.StationCode + orig.StationCode
			}

			od := ODPair{
				origin:      orig,
				destination: dest,
				key:         key,
			}

			pairs = append(pairs, od)
		}
	}

	// filter so A->B === B->A
	pairs = lo.UniqBy(pairs, func(item ODPair) string {
		return item.key
	})

	pairs = lo.Shuffle(pairs)
	pairs = pairs[0:1000]

	odChan := make(chan ODPair)
	go func() {
		defer close(odChan)
		for _, pair := range pairs {
			odChan <- pair
		}
	}()

	n := runtime.GOMAXPROCS(0)
	wg := sync.WaitGroup{}
	wg.Add(n)
	start := time.Now()
	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			for od := range odChan {
				delay := 100 + rand.Intn(150)
				time.Sleep(time.Duration(delay) * time.Millisecond)
				//log.Printf("routine %d, trying %s", id, od.key)
				_, sErr := s.getSchedule(od)
				if sErr != nil {
					log.Printf("error getting schedule for %s: %s", od.key, sErr)
					continue
				}
			}
		}(i)
	}
	wg.Wait()
	end := time.Now()
	took := end.Sub(start)
	log.Printf("took %s", took)
	return nil, nil
}

func (s ScheduleScraper) getStops() ([]Stop, error) {
	payload := `{"carrierId":221,"type":"ORIGIN"}`
	url := s.baseUrl + "/stop"
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Tds-Api-Key", s.apiKey)
	req.Header.Add("TDS-Carrier-Code", "PPB")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d, %s", resp.StatusCode, resp.Status)
	}

	var stops []Stop
	decErr := json.NewDecoder(resp.Body).Decode(&stops)
	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing body: %w", err)
	}

	if decErr != nil {
		return nil, fmt.Errorf("error decoding response: %w", decErr)
	}
	return stops, nil
}

func (s ScheduleScraper) getSchedule(od ODPair) ([]ODCandidate, error) {
	url := s.baseUrl + "/v2/schedule"
	payload := fmt.Sprintf(`{
			  "purchaseType": "SCHEDULE_BOOK",
			  "origin": {
				"stopUuid": "%s"
			  },
			  "destination": {
				"stopUuid": "%s"
			  },
			  "departDate": "2024-07-24",
			  "cityMode": false,
			  "isReturn": false,
			  "passengerCounts": {
				"Adult": 1
			  }
			}`, od.origin.StopUuid, od.destination.StopUuid)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Tds-Api-Key", s.apiKey)
	req.Header.Add("TDS-Carrier-Code", "PPB")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d, %s", resp.StatusCode, resp.Status)
	}

	//var stops []Stop
	//decErr := json.NewDecoder(resp.Body).Decode(&stops)
	//err = resp.Body.Close()
	//if err != nil {
	//	return nil, fmt.Errorf("error closing body: %w", err)
	//}
	//
	//if decErr != nil {
	//	return nil, fmt.Errorf("error decoding response: %w", decErr)
	//}
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	log.Println(string(response))
	return nil, nil
}

type ODCandidate struct {
	OD         ODPair
	ScheduleID string
}
