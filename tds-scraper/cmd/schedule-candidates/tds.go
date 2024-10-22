package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Stop struct {
	StopUuid string `json:"stopUuid"`
}

type ScheduleRequest struct {
	PassengerCounts map[string]int `json:"passengerCounts"`
	PurchaseType    string         `json:"purchaseType"`
	Origin          Stop           `json:"origin"`
	Destination     Stop           `json:"destination"`
	DepartDate      string         `json:"departDate"`
}

type ScheduleResponse struct {
	ScheduleProducts []ScheduleProduct `json:"scheduleProducts"`
}

type ScheduleProduct struct {
	ScheduleRun ScheduleRun `json:"scheduleRun"`
}

type ScheduleRun struct {
	ScheduleUuid string
}

type ScheduleAPI struct {
	client  http.Client
	baseUrl string
	apiKey  string
	carrier string
}

func (s ScheduleAPI) Get(scheduleRequest ScheduleRequest) (ScheduleResponse, error) {
	var scheduleResponse ScheduleResponse

	payload, err := json.Marshal(scheduleRequest)
	if err != nil {
		return ScheduleResponse{}, err
	}

	log.Printf("PAYLOAD: \n%s\n", string(payload))
	req, err := http.NewRequest("POST", s.baseUrl+"/schedule", bytes.NewBuffer(payload))
	if err != nil {
		return scheduleResponse, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("TdsApiKey", s.apiKey)
	req.Header.Add("TDS-Carrier-Code", s.carrier)

	resp, err := s.client.Do(req)
	if err != nil {
		return scheduleResponse, fmt.Errorf("error executing request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return scheduleResponse, fmt.Errorf("bad status code: %d, %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return scheduleResponse, err
	}

	err = json.Unmarshal(buf, &scheduleResponse)

	return scheduleResponse, err
}

func trySchedule(schedules ScheduleAPI, days int, origin, dest string) error {
	var dates []string
	for i := range days {
		dt := time.Now().Add(time.Duration(i+1) * 24 * time.Hour)
		dates = append(dates, dt.Format("2006-01-02"))
	}
	log.Printf("dates: %v", dates)
	return nil
}
