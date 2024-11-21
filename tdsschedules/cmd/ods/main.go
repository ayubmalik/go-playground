// Command ods finds origin destination pairs
package main

import (
	"context"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"sync"
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

	slog.Info("creating client with", "apiKey", apiKey, "carrierCode", carrierCode)
	tdsClient := tdsschedules.NewTDSClient(apiKey, carrierCode)

	candidateODs := getOriginDestinationCandidates(tdsClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	findODPairs(ctx, tdsClient, candidateODs)
}

func findODPairs(ctx context.Context, client tdsschedules.TdsClient, candidates <-chan ODPair) {
	for candidate := range candidates {
		select {
		case <-ctx.Done():
			return
		default:
			tryODPair(ctx, client, candidate)
		}
	}
}

// TODO error handling
func tryODPair(ctx context.Context, client tdsschedules.TdsClient, candidate ODPair) {
	// find a schedule starting on a monday
	departDate := tdsschedules.NextMonday(time.Now())
	if scheduleExists(ctx, client, departDate, candidate) {
		// TODO save OD
		return
	}

	// otherwise try other count of week
	count := 6
	wg := sync.WaitGroup{}
	wg.Add(count)

}

func scheduleExists(ctx context.Context, client tdsschedules.TdsClient, departDate time.Time, candidate ODPair) bool {
	slog.Debug("searching for schedules", "departDate", departDate, "o", candidate.Origin.StationCode, "d", candidate.Destination.StationCode)
	schedules, err := client.SearchSchedules(ctx, candidate.Origin, candidate.Destination, departDate)
	if err != nil {
		slog.Error("searching for schedules", "o", candidate.Origin.StationCode, "d", candidate.Destination.StationCode, "error", err)
		return false
	}

	if schedules.IsEmpty() {
		return false
	}

	slog.Info("found schedule", "id", schedules.FirstID())
	return true
}

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
