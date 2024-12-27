// Command ods finds origin destination pairs
package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strings"
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
	dbUrl := os.Getenv("LOCAL_DATABASE_URL")

	ctx, cancel := context.WithCancel(context.Background())
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil || conn.Ping(ctx) != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		cancel()
		_ = conn.Close(ctx)
	}()

	slog.Info("creating client with", "apiKey", apiKey, "carrierCode", carrierCode)
	tdsClient := tdsschedules.NewTDSClient(apiKey, carrierCode)

	stopSummaryDB := tdsschedules.NewStopSummaryDB(conn)

	candidateODs := getOriginDestinationCandidates(ctx, tdsClient, stopSummaryDB)

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

	//// otherwise try other count of week
	//count := 6
	//wg := sync.WaitGroup{}
	//wg.Add(count)
	//
	//for i := 0; i < count; i++ {
	//	go func() {
	//		defer wg.Done()
	//		if scheduleExists(ctx, client, departDate.Add(24*time.Hour), candidate) {
	//			// save OD
	//		}
	//	}()
	//}
	//wg.Wait()
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

	slog.Info("found schedule",
		"day", departDate.Weekday(),
		"date", departDate.Format(time.DateOnly),
		"o", candidate.Origin.StationCode,
		"d", candidate.Destination.StationCode,
		"id", schedules.FirstID(),
	)
	return true
}

func getOriginDestinationCandidates(ctx context.Context, tdsClient tdsschedules.TdsClient, db *tdsschedules.StopSummaryDB) <-chan ODPair {
	stops, err := tdsClient.FindStops()
	if err != nil {
		slog.Error("Error finding stops", "error", err)
	}
	slog.Info("saving found stops", "count", len(stops))
	for _, stop := range stops {
		slog.Info("delete stop", "stop.id", stop.StopUuid, "stop.name", stop.Name, "stop.code", stop.StationCode)
		err = db.Delete(ctx, strings.TrimSpace(stop.StopUuid))
		if err != nil {
			slog.Error("Error deleting stop", "error", err)
			os.Exit(1)
		}

		err = db.Put(ctx, tdsschedules.StopSummary{
			ID:    stop.StopUuid,
			Name:  stop.Name,
			Code:  stop.StationCode,
			City:  stop.City.Name,
			State: stop.State.Abbreviation,
		})
		if err != nil {
			slog.Error("Error saving stop summary", "error", err, "stop.id", stop.StopUuid, "stop.name", stop.Name, "stop.city", stop.City.Name)
		}
	}

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

func createPairs(stops []tdsschedules.Stop) []ODPair {
	pairs := make([]ODPair, 0, len(stops)*(len(stops)-1)/2)
	for i := 0; i < len(stops); i++ {
		for j := i + 1; j < len(stops); j++ {
			pairs = append(pairs, ODPair{stops[i], stops[j]})
		}
	}
	return pairs
}
