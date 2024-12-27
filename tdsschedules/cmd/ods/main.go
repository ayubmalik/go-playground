// Command ods finds origin destination pairs
package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"slices"
	"strings"
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

	setLogLevel()

	apiKey := os.Getenv("TDS_API_KEY")
	carrierCode := os.Getenv("TDS_CARRIER_CODE")
	dbUrl := os.Getenv("DATABASE_URL")

	ctx, cancel := context.WithCancel(context.Background())
	conn, err := pgxpool.New(ctx, dbUrl)
	if err != nil || conn.Ping(ctx) != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		cancel()
		conn.Close()
	}()

	slog.Info("creating client with", "apiKey", apiKey, "carrierCode", carrierCode)
	tdsClient := tdsschedules.NewTDSClient(apiKey, carrierCode)

	stopSummaryDB := tdsschedules.NewStopSummaryDB(conn)
	originDestinationDB := tdsschedules.NewOriginDestinationDB(conn)

	slog.Info("deleting OD pairs")
	err = originDestinationDB.DeleteAll(ctx)
	if err != nil {
		slog.Error("Error deleting all OD pairs", "error", err)
		os.Exit(1)
	}

	slog.Info("deleting all stop summaries")
	err = stopSummaryDB.DeleteAll(ctx)
	if err != nil {
		slog.Error("Error deleting all stop summaries", "error", err)
		os.Exit(1)
	}

	candidateODs := getOriginDestinationCandidates(ctx, tdsClient, stopSummaryDB)
	findODPairs(ctx, tdsClient, candidateODs, originDestinationDB)
}

func setLogLevel() {
	levelName := "INFO"
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		levelName = level
	}

	var level slog.Level
	err := level.UnmarshalText([]byte(levelName))
	if err != nil {
		slog.Error("Error unmarshalling LOG_LEVEL", "error", err)
		os.Exit(1)
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})))
}

func findODPairs(ctx context.Context, client tdsschedules.TdsClient, candidates <-chan ODPair, db *tdsschedules.OrigDestinationDB) {
	for candidate := range candidates {
		select {
		case <-ctx.Done():
			return
		default:
			tryODPair(ctx, client, candidate, db)
		}
	}
}

// TODO error handling
func tryODPair(ctx context.Context, client tdsschedules.TdsClient, candidate ODPair, db *tdsschedules.OrigDestinationDB) {
	// find a schedule starting on a monday
	departDate := tdsschedules.NextMonday(time.Now())
	slog.Debug("trying non monday OD pair", "departDate", departDate, "origin", candidate.Origin.StationCode, "destination", candidate.Destination.StationCode)
	if scheduleExists(ctx, client, departDate, candidate) {
		slog.Info("found OD pair", "departDate", departDate, "origin", candidate.Origin.StationCode, "destination", candidate.Destination.StationCode)

		od := createOD(candidate)
		err := db.Put(ctx, od)
		if err != nil {
			slog.Error("Error putting OD to DB", "error", err)
		}
		return
	}

	// otherwise try other count of week
	count := 6
	wg := sync.WaitGroup{}
	wg.Add(count)

	ctx2, cancel := context.WithCancel(context.Background())
	for i := 1; i <= count; i++ {
		go func() {
			defer wg.Done()
			slog.Debug("trying non monday OD pair", "departDate", departDate, "origin", candidate.Origin.StationCode, "destination", candidate.Destination.StationCode)
			if scheduleExists(ctx2, client, departDate.Add(24*time.Duration(i)*time.Hour), candidate) {
				slog.Info("found OD pair", "departDate", departDate, "origin", candidate.Origin.StationCode, "destination", candidate.Destination.StationCode)
				od := createOD(candidate)
				err := db.Put(ctx, od)
				if err != nil {
					slog.Error("Error putting OD to DB", "error", err)
				}
				cancel()
			}
		}()
	}
	wg.Wait()
	cancel()
}

func createOD(candidate ODPair) tdsschedules.OriginDestination {
	od := tdsschedules.OriginDestination{
		Origin: tdsschedules.StopSummary{
			ID:    candidate.Origin.StopUuid,
			Name:  candidate.Origin.Name,
			Code:  candidate.Origin.StationCode,
			City:  candidate.Origin.City.Name,
			State: candidate.Origin.State.Abbreviation,
		},
		Destination: tdsschedules.StopSummary{
			ID:    candidate.Destination.StopUuid,
			Name:  candidate.Destination.Name,
			Code:  candidate.Destination.StationCode,
			City:  candidate.Destination.City.Name,
			State: candidate.Destination.State.Abbreviation,
		},
	}
	return od
}

func scheduleExists(ctx context.Context, client tdsschedules.TdsClient, departDate time.Time, candidate ODPair) bool {
	slog.Debug("searching for schedules", "departDate", departDate, "o", candidate.Origin.StationCode, "d", candidate.Destination.StationCode)
	schedules, err := client.SearchSchedules(ctx, candidate.Origin, candidate.Destination, departDate)
	if err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("searching for schedules", "o", candidate.Origin.StationCode, "d", candidate.Destination.StationCode, "error", err)
		return false
	}

	if schedules.IsEmpty() {
		return false
	}

	slog.Debug("found schedule",
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
		pairs := createPairs(stops)
		for _, pair := range pairs {
			candidates <- pair
		}
	}()

	return candidates
}

func createPairs(stops []tdsschedules.Stop) []ODPair {
	slices.SortFunc(stops, func(l, r tdsschedules.Stop) int {
		return strings.Compare(l.StationCode, r.StationCode)
	})

	pairs := make([]ODPair, 0, len(stops)*(len(stops)-1)/2)
	for i := 0; i < len(stops); i++ {
		for j := i + 1; j < len(stops); j++ {
			pairs = append(pairs, ODPair{stops[i], stops[j]})
		}
	}
	return pairs
}
