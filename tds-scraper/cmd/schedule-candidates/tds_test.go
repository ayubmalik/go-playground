package main

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestNewScheduleRequest(t *testing.T) {
	t.Run("new schedule request", func(t *testing.T) {
		req := newScheduleQuery("orig1", "dest1", "2025-12-25")

		if req.PurchaseType != "SCHEDULE_BOOK" {
			t.Errorf("req.purchaseType = %s, want %s", req.PurchaseType, "SCHEDULE_BOOK")
		}

		j, _ := json.Marshal(req)
		t.Log(string(j))

	})
}

func TestTdsClient(t *testing.T) {
	tds := createTDSClient()

	newYork := "83be15f2-118b-45d9-839c-c92e841f10fd"
	newPaltz := "bbd3cdc1-0e9e-4869-b337-abcb6868bf41"
	departure := "2024-10-31"
	qry := newScheduleQuery(newYork, newPaltz, departure)

	result, err := tds.FindSchedules(qry)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
	t.Logf("IsEmpty %v", result.IsEmpty())
}

// Only example how to write a range
func doRange(days int) func(yield func(date string) bool) {
	return func(yield func(date string) bool) {
		for i := range days {
			dt := time.Now().Add(24 * time.Hour * time.Duration(i+1)).Format("2006-01-02")
			if !yield(dt) {
				return
			}
		}
	}
}

func TestDoRange(t *testing.T) {
	doRange(1)
	t.Skip()
}

func createTDSClient() TdsClient {
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

	tds := TdsClient{
		client:  client,
		baseUrl: "https://ride-api.bustickets.com/tickets/v2",
		apiKey:  "E54589A3-85E1-43D5-90C4-E0F33645CF1A",
		carrier: "BTC",
	}
	return tds
}
