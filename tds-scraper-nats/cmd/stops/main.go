package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	return fmt.Sprintf("%s %s", od.Origin.StationName, od.Destination.StationName)
}

type StopQuery struct {
	Type      string `json:"type"`
	Origin    Stop   `json:"origin"`
	CarrierId int    `json:"carrierId"`
}

func main() {
	var credentials string
	flag.StringVar(&credentials, "credentials", "", "credentials file location")
	flag.Parse()

	if credentials == "" {
		log.Fatalln("-credentials flag needs to be set")
	}

	apiKey := "E54589A3-85E1-43D5-90C4-E0F33645CF1A"

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()

	origins, err := GetOrigins(client, apiKey)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("got %d origins", len(origins))
	destinations, err := GetDestinations(client, apiKey, origins)
	if err != nil {
		log.Fatalln(err)
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

	took := time.Since(start)
	log.Printf("count pairs: %d\n", len(pairs))
	log.Printf("took %s\n", took)

	nc, err := nats.Connect("tls://connect.ngs.global", nats.UserCredentials(credentials), nats.Name("pub.schedules.candidates"))
	if err != nil {
		log.Fatalf("could not connect to nats %s\n", err)
	}

	defer func() {
		_ = nc.Drain()
	}()

	for id, p := range pairs {
		data := p.Origin.StopUuid + " " + p.Destination.StopUuid

		hdr := nats.Header{}
		hdr.Add(nats.MsgIdHdr, strconv.Itoa(id))
		msg := &nats.Msg{
			Subject: "tds.schedules.candidates",
			Header:  hdr,
			Data:    []byte(data),
		}

		err = nc.PublishMsg(msg)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func GetOrigins(client http.Client, apiKey string) ([]Stop, error) {
	qry := StopQuery{
		CarrierId: carrierID,
		Type:      "ORIGIN",
	}

	origins, err := getStops(client, apiKey, qry)
	if err != nil {
		return nil, fmt.Errorf("error getting origin stops: %w", err)
	}

	return origins, nil
}

func GetDestinations(client http.Client, apiKey string, origins []Stop) ([]Stop, error) {
	qry := StopQuery{
		CarrierId: carrierID,
		Type:      "DESTINATION",
	}

	destinations, err := getStops(client, apiKey, qry)
	if err != nil {
		return nil, err
	}

	return destinations, nil
}

func getStops(client http.Client, apiKey string, qry StopQuery) ([]Stop, error) {
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
