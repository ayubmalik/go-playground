package tdschedules_test

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

func TestStops(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := startPostgres(ctx)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	connString, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get database connection string: %s", err)
	}

	defer func() {
		if err2 := testcontainers.TerminateContainer(pgContainer); err2 != nil {
			log.Fatalf("failed to terminate container: %s", err2)
		}
	}()

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	t.Cleanup(func() {
		err := conn.Close(ctx)
		if err != nil {
			t.Fatalf("Error closing connection: %v", err)
		}
		err = testcontainers.TerminateContainer(pgContainer)
		if err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	t.Run("ping database", func(t *testing.T) {
		if err := conn.Ping(ctx); err != nil {
			t.Fatalf("failed to ping database: %v", err)
		}

	})
}

func startPostgres(ctx context.Context) (*postgres.PostgresContainer, error) {
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
	return pgContainer, err
}
