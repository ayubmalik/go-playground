package tdsschedules

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type City struct {
	Name string `json:"name"`
}

type State struct {
	Abbreviation string `json:"abbreviation"`
	Country      string `json:"country"`
}

type Stop struct {
	StopUuid    string `json:"stopUuid"`
	Name        string `json:"displayBusinessName,omitempty"`
	StationCode string `json:"stationCode,omitempty"`
	City        *City  `json:"city,omitempty"`
	State       *State `json:"state,omitempty"`
}

type ScheduleQuery struct {
	PassengerCounts map[string]int `json:"passengerCounts"`
	PurchaseType    string         `json:"purchaseType"`
	Origin          Stop           `json:"origin"`
	Destination     Stop           `json:"destination"`
	DepartDate      string         `json:"departDate"`
}

type ScheduleResult struct {
	ScheduleProducts []ScheduleProduct
}

func (sr ScheduleResult) IsEmpty() bool {
	return len(sr.ScheduleProducts) == 0
}

func (sr ScheduleResult) FirstOriginDestination() (Stop, Stop) {
	if sr.IsEmpty() {
		return Stop{}, Stop{}
	}

	return sr.ScheduleProducts[0].ScheduleRun.Origin, sr.ScheduleProducts[0].ScheduleRun.Destination
}

func (sr ScheduleResult) FirstID() string {
	if sr.IsEmpty() {
		return ""
	}

	return sr.ScheduleProducts[0].ScheduleRun.ScheduleUuid
}

type ScheduleProduct struct {
	ScheduleRun ScheduleRun
}

type ScheduleRun struct {
	ScheduleUuid string
	Origin       Stop
	Destination  Stop
}

func NewTDSClient(apiKey, carrierCode string) TdsClient {
	timeout := 30 * time.Second
	transport := http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
		TLSHandshakeTimeout: timeout,
		MaxConnsPerHost:     10,
		MaxIdleConns:        10,
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   timeout,
	}

	tds := TdsClient{
		client:  client,
		baseUrl: "https://ride-api.bustickets.com/tickets",
		apiKey:  apiKey,
		carrier: carrierCode,
	}
	return tds
}

type TdsClient struct {
	client  http.Client
	baseUrl string
	apiKey  string
	carrier string
}

func (t TdsClient) SearchSchedules(ctx context.Context, origin Stop, destination Stop, departDate time.Time) (ScheduleResult, error) {
	var result ScheduleResult
	query := ScheduleQuery{
		PurchaseType:    "SCHEDULE_BOOK",
		Origin:          Stop{StopUuid: origin.StopUuid},
		Destination:     Stop{StopUuid: destination.StopUuid},
		DepartDate:      departDate.Format(time.DateOnly),
		PassengerCounts: map[string]int{"Adult": 1},
	}
	payload, err := json.Marshal(query)
	if err != nil {
		return result, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", t.baseUrl+"/v2/schedule", bytes.NewBuffer(payload))
	if err != nil {
		return result, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("TdsApiKey", t.apiKey)
	req.Header.Add("TDS-Carrier-Code", t.carrier)

	resp, err := t.client.Do(req)
	if err != nil {
		return result, fmt.Errorf("error executing request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("bad status code: %d, %s", resp.StatusCode, resp.Status)
	}
	defer func() { _ = resp.Body.Close() }()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(buf, &result)

	return result, err
}

func (t TdsClient) FindStops() ([]Stop, error) {
	payload := `{"carrierId":221,"type":"ORIGIN"}`

	req, err := http.NewRequest("POST", t.baseUrl+"/stop", bytes.NewBufferString(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating stop request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Tds-Api-Key", t.apiKey)
	req.Header.Add("TDS-Carrier-Code", t.carrier)
	resp, err := t.client.Do(req)
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
		return nil, fmt.Errorf("error closing payload: %w", err)
	}

	if decErr != nil {
		return nil, fmt.Errorf("error decoding response: %w", decErr)
	}
	return stops, nil
}
