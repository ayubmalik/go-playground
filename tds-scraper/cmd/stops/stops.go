package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	carrierID   = 221
	carrierCode = "PPB"
	baseUrl     = "https://ride-api.bustickets.com/tickets"
)

type City struct {
	Name   string `json:"name"`
	CityId int    `json:"cityId"`
}

type Stop struct {
	StopUuid    string `json:"stopUuid"`
	StationName string `json:"stationName"`
	StationCode string `json:"stationCode"`
	City        City   `json:"city"`
}

type ODPair struct {
	Origin      Stop
	Destination Stop
}

func (od ODPair) String() string {
	return fmt.Sprintf("%s -> %s", od.Origin.StationName, od.Destination.StationName)
}

type StopQuery struct {
	Type      string `json:"type"`
	CarrierId int    `json:"carrierId"`
	CityMode  bool   `json:"cityMode"`
}

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("API_KEY env must be set")
		os.Exit(1)
	}

	client := http.Client{}
	start := time.Now()
	qry := StopQuery{
		CarrierId: carrierID,
		Type:      "ORIGIN",
	}

	origins, err := getStops(qry, apiKey, client)
	if err != nil {
		fmt.Printf("error getting origin stops: %v\n", err)
		os.Exit(1)
	}

	qry = StopQuery{
		CarrierId: carrierID,
		Type:      "DESTINATION",
	}

	destinations, err := getStops(qry, apiKey, client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var pairs []ODPair
	for _, o := range origins {
		for _, d := range destinations {
			if o.StopUuid == d.StopUuid {
				continue
			}

			od := ODPair{Origin: o, Destination: d}
			pairs = append(pairs, od)
		}
	}

	fmt.Printf("count pairs: %d\n", len(pairs))
	took := time.Since(start)
	fmt.Printf("took %s\n", took)

	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Printf("nats conn err: %s\nDef url %s\n", err, nats.DefaultURL)
		os.Exit(1)
	}

	for _, p := range pairs {
		data, err := json.Marshal(p)
		if err != nil {
			fmt.Printf("could not marshal pair: %v", err)
			os.Exit(1)
		}

		msg := &nats.Msg{
			Subject: "tds.q.schedules.candidates",
			Data:    data,
		}

		err = conn.PublishMsg(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	defer conn.Close()

}

func getStops(qry StopQuery, apiKey string, client http.Client) ([]Stop, error) {
	body, err := json.Marshal(qry)
	if err != nil {
		return nil, fmt.Errorf("could not marshal query: %w", err)
	}

	url := baseUrl + "/stop"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Tds-Api-Key", apiKey)
	req.Header.Add("TDS-Carrier-Code", carrierCode)
	resp, err := client.Do(req)
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
