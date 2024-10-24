package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Stop struct {
	StopUuid string `json:"stopUuid"`
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

type ScheduleProduct struct {
	ScheduleRun ScheduleRun
}

type ScheduleRun struct {
	ScheduleUuid string
}

type TdsClient struct {
	client  http.Client
	baseUrl string
	apiKey  string
	carrier string
}

func (t TdsClient) FindSchedules(qry ScheduleQuery) (ScheduleResult, error) {
	var result ScheduleResult

	payload, err := json.Marshal(qry)
	if err != nil {
		return result, err
	}

	log.Printf("PAYLOAD: \n%s\n", string(payload))
	req, err := http.NewRequest("POST", t.baseUrl+"/schedule", bytes.NewBuffer(payload))
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
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(buf, &result)

	return result, err
}
