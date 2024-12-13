package main

import (
	"context"
	"gotest.tools/v3/assert"
	"os"
	"tdsschedules"
	"testing"
	"time"
)

func TestSchedules(t *testing.T) {
	apiKey := os.Getenv("TDS_API_KEY")
	carrierCode := os.Getenv("TDS_CARRIER_CODE")

	t.Logf("creating client with apiKey = %s, carrierCode=%s", apiKey, carrierCode)
	tdsClient := tdsschedules.NewTDSClient(apiKey, carrierCode)

	ctx := context.Background()
	origin := tdsschedules.Stop{StopUuid: "35e5b11d-8b14-44a7-8112-cbe297c4005e"}
	dest := tdsschedules.Stop{StopUuid: "83be15f2-118b-45d9-839c-c92e841f10fd"}
	depart := tdsschedules.NextMonday(time.Now())

	schedules, err := tdsClient.SearchSchedules(ctx, origin, dest, depart)
	if err != nil {
		t.Fatalf("failed to search schedules: %v", err)
	}
	t.Logf("schedules = %v", schedules.IsEmpty())
}

func TestCreatePairs(t *testing.T) {
	t.Run("create pairs", func(t *testing.T) {
		stops := []tdsschedules.Stop{
			{StopUuid: "a"}, {StopUuid: "b"}, {StopUuid: "c"}, {StopUuid: "d"}, {StopUuid: "e"},
		}

		pairs := createPairs(stops)
		for _, pair := range pairs {
			t.Logf("pair = %v", pair)
		}
		assert.Equal(t, len(pairs), 10)
	})
}
