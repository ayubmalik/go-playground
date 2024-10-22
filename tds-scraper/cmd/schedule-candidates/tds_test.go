package main

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestTrySchedule(t *testing.T) {

	t.Run("days range", func(t *testing.T) {

		var api ScheduleAPI
		days := 7

		err := trySchedule(api, days, "origin", "dest")

		if err != nil {
			t.Errorf("not trySchedule error %v", err)
		}
	})
}

func TestNewScheduleRequest(t *testing.T) {
	t.Run("new schedule request", func(t *testing.T) {
		req := newScheduleRequest("orig1", "dest1", "2025-12-25")

		if req.PurchaseType != "SCHEDULE_BOOK" {
			t.Errorf("req.purchaseType = %s, want %s", req.PurchaseType, "SCHEDULE_BOOK")
		}

		j, _ := json.Marshal(req)
		t.Log(string(j))

	})
}

func newScheduleRequest(origin, dest, departure string) ScheduleRequest {
	return ScheduleRequest{
		PurchaseType:    "SCHEDULE_BOOK",
		Origin:          Stop{origin},
		Destination:     Stop{dest},
		DepartDate:      departure,
		PassengerCounts: map[string]int{"Adult": 1},
	}
}

func TestScheduleAPI(t *testing.T) {
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   10 * time.Second,
	}

	api := ScheduleAPI{
		client:  client,
		baseUrl: "https://ride-api.bustickets.com/tickets/v2",
		apiKey:  "E54589A3-85E1-43D5-90C4-E0F33645CF1A",
		carrier: "BTC",
	}

	newYork := "0e75ad0a-eff6-491b-9597-dbb262509d40"
	newPaltz := "bbd3cdc1-0e9e-4869-b337-abcb6868bf41"
	departure := "2024-10-31"
	sr := newScheduleRequest(newYork, newPaltz, departure)

	response, err := api.Get(sr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(response)

}
