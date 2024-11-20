// Command ods finds origin destination pairs
package main

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"tdsschedules"
	"time"
)

type ODPair struct {
	Origin      tdsschedules.Stop
	Destination tdsschedules.Stop
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	apiKey := os.Getenv("TDS_API_KEY")
	carrierCode := os.Getenv("TDS_CARRIER_CODE")
	slog.Info("got environment variables", "apiKey", apiKey, "carrierCode", carrierCode)

	tdsClient := tdsschedules.NewTDSClient(apiKey, carrierCode)

	candidates := getOriginDestinationCandidates(tdsClient)
	findODSchedules(tdsClient, candidates)

}

func nextMonday(date time.Time) time.Time {
	offset := (8 - date.Weekday()) % 7
	if offset == 0 {
		offset = 7
	}
	return date.Add(24 * time.Hour * time.Duration(offset))
}

func findODSchedules(client tdsschedules.TdsClient, candidates <-chan ODPair) {
	for candidate := range candidates {
		departureDate := nextMonday(time.Now())
		slog.Info("finding schedule for", "departureDate", departureDate, "candidate", candidate)
		slog.Info("finding next schedule for", "departureDate", departureDate.Add(24*time.Hour))
	}
}

// getOriginDestinationCandidates returns a channel of all possible OD combinations
func getOriginDestinationCandidates(tdsClient tdsschedules.TdsClient) <-chan ODPair {
	stops, err := tdsClient.FindStops()
	if err != nil {
		slog.Error("Error finding stops", "error", err)
	}

	slog.Info("found stops", "count", len(stops))
	candidates := make(chan ODPair)

	go func() {
		defer close(candidates)
		for _, origin := range stops {
			for _, destination := range stops {
				if origin.StopUuid != destination.StopUuid {
					candidates <- ODPair{origin, destination}
				}
			}
		}
	}()

	return candidates
}
