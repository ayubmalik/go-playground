package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type City struct {
	CityId int    `json:"cityId"`
	Name   string `json:"name"`
}

type State struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Country      string `json:"country"`
}

type Stop struct {
	StopUuid    string `json:"stopUuid"`
	StationName string `json:"stationName"`
	City        City   `json:"city"`
	State       State  `json:"state"`
}

type ODPair struct {
	origin      Stop
	destination Stop
}

func main() {

	hc := http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			ReadBufferSize:      16 * 1024 * 1024,
		},
	}

	apiKey := os.Getenv("API_KEY")
	payload := `{"carrierId":304,"type":"ORIGIN"}`
	stops, err := getStops(hc, apiKey, payload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("count:", len(stops))
	//var pairs []ODPair
	for _, stop := range stops {
		fmt.Printf("origin %s, %s\n", stop.StationName, stop.City.Name)
		payload = fmt.Sprintf(`{"carrierId": 304 ,"type": "DESTINATION", "origin": {"stopUuid": "%s"}}`, stop.StopUuid)
		fmt.Println(payload)
		time.Sleep(100 * time.Millisecond)
		dests, err := getStops(hc, apiKey, payload)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("dest count %d\n", len(dests))

	}
}

func getStops(hc http.Client, apiKey, payload string) ([]Stop, error) {
	req, err := http.NewRequest("POST", "https://ride-api.bustickets.com/tickets/stop", strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Tds-Api-Key", apiKey)
	req.Header.Add("TDS-Carrier-Code", "BTC")
	resp, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d, %s", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	var stops []Stop
	decErr := json.NewDecoder(resp.Body).Decode(&stops)
	if decErr != nil {
		return nil, fmt.Errorf("error decoding response: %w", decErr)
	}
	return stops, nil
}
