package tdsschedules_test

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"tdsschedules"
	"testing"
)

func TestStopSummaryDB(t *testing.T) {
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

	db := tdsschedules.NewStopSummaryDB(conn)

	t.Run("put and get all stops", func(t *testing.T) {
		stops := []tdsschedules.StopSummary{
			{ID: uuid.New(), Name: "name1", Code: "c1de", City: "city", State: "SS"},
			{ID: uuid.New(), Name: "name2", Code: "c2de", City: "city", State: "SS"},
			{ID: uuid.New(), Name: "name3", Code: "c3de", City: "city3", State: "S3"},
		}

		for i, stop := range stops {
			if err := db.Put(ctx, stop); err != nil {
				t.Fatalf("failed to put stop[%d]: %s", i, err)
			}
		}

		gotStops, _ := db.GetAll(ctx)
		if len(gotStops) != len(stops) {
			t.Errorf("len(GetAll()) = %v, want %v", len(gotStops), len(stops))
		}

		for i, stop := range gotStops {
			t.Logf("stop[%d]: %v", i, stop)
			if !cmp.Equal(stop, stops[i]) {
				t.Errorf("stop = %v, want = %v", stop, stops[i])
			}
		}
	})
}
