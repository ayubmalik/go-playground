package tdsschedules_test

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"tdsschedules"
	"testing"
	"time"
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
			{ID: "00000000-0001-f000-0000-000000000000", Name: "name1", Code: "c1de", City: "city", State: "SS"},
			{ID: "00000000-0002-f000-0000-000000000000", Name: "name2", Code: "c2de", City: "city", State: "SS"},
			{ID: "00000000-0003-f000-0000-000000000000", Name: "name3", Code: "c3de", City: "city3", State: "S3"},
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
			if !cmp.Equal(stop, stops[i]) {
				t.Errorf("stop = %v, want = %v", stop, stops[i])
			}
		}
	})
}

func startPostgresWithMigrations(ctx context.Context) (*pgx.Conn, *postgres.PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("tds"),
		postgres.WithUsername("tds"),
		postgres.WithPassword("tds"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	var connString string
	if pgContainer != nil {
		connString, err = pgContainer.ConnectionString(ctx, "sslmode=disable")
	}

	var conn *pgx.Conn
	if connString != "" {
		conn, err = pgx.Connect(ctx, connString)
	}

	if conn != nil {
		err = conn.Ping(ctx)
	}

	var m *migrate.Migrate
	if err == nil {
		m, err = migrate.New("file://migrations", connString)
	}

	if m != nil {
		err = m.Steps(2)
	}

	return conn, pgContainer, err
}
