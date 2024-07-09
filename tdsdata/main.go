package main

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"net/http"
	"os"
	"strings"
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
	payload := `{"carrierId":221,"type":"ORIGIN"}`
	stops, err := getStops(hc, apiKey, payload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("count:", len(stops))
	var pairs []ODPair
	for _, orig := range stops {
		for _, dest := range stops {
			var key string
			if strings.Compare(orig.StationCode, dest.StationName) < 0 {
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

	fmt.Println("pairs:", len(pairs))

	pairs = lo.UniqBy(pairs, func(item ODPair) string {
		return item.key
	})

	for _, pair := range pairs {
		fmt.Println(pair.key)
	}

	fmt.Println("pairs2:", len(pairs))

}

func getStops(hc http.Client, apiKey, payload string) ([]Stop, error) {
	req, err := http.NewRequest("POST", "https://ride-api.bustickets.com/tickets/stop", strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Tds-Api-Key", apiKey)
	req.Header.Add("TDS-Carrier-Code", "PPB")
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
