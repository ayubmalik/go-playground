package tdsschedules_test

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"tdsschedules"
	"testing"
)

func TestOrigDestinationDB(t *testing.T) {
	ctx := context.Background()

	conn, pgContainer, err := startPostgresWithMigrations(ctx)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	t.Cleanup(func() {
		err := conn.Close(ctx)
		if err != nil {
			t.Logf("Error closing connection: %v", err)
		}

		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	stops := []tdsschedules.StopSummary{
		{ID: uuid.New().String(), Name: "name1", Code: "c1de", City: "city", State: "SS"},
		{ID: uuid.New().String(), Name: "name2", Code: "c2de", City: "city", State: "SS"},
		{ID: uuid.New().String(), Name: "name3", Code: "c3de", City: "city", State: "SS"},
	}

	insertStops(t, ctx, conn, stops)

	db := tdsschedules.NewOriginDestinationDB(conn)

	t.Run("put and get all origin destinations", func(t *testing.T) {
		ods := []tdsschedules.OriginDestination{
			{Origin: stops[0], Destination: stops[1]},
			{Origin: stops[0], Destination: stops[2]},
			{Origin: stops[1], Destination: stops[2]},
			{Origin: stops[0], Destination: stops[1]},
			{Origin: stops[0], Destination: stops[2]},
			{Origin: stops[1], Destination: stops[2]},
		}

		for _, od := range ods {
			if err := db.Put(ctx, od); err != nil {
				t.Fatalf("failed to put origin: %s", err)
			}
		}

		gotODs, err := db.GetAll(ctx)
		if err != nil {
			t.Fatalf("failed to get all origin: %s", err)
		}
		if len(gotODs) != 3 {
			t.Errorf("len(GetAll()) = %v, want %v", len(gotODs), 3)
		}

		for i, od := range gotODs {
			t.Logf("od[%d]: %v", i, od)
			if !cmp.Equal(od, ods[i]) {
				t.Errorf("od = %v, want = %v", od, stops[i])
			}
		}
	})
}

func insertStops(t *testing.T, ctx context.Context, conn *pgx.Conn, stops []tdsschedules.StopSummary) {
	t.Helper()
	stopsDB := tdsschedules.NewStopSummaryDB(conn)
	for i, stop := range stops {
		if err := stopsDB.Put(ctx, stop); err != nil {
			t.Fatalf("failed to put stop[%d]: %s", i, err)
		}
	}
}
