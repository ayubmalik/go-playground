package tdsschedules_test

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	db := tdsschedules.NewOriginDestinationDB(conn)

	t.Run("put and get all origin destinations", func(t *testing.T) {
		t.Fatalf("TODO! %v", db)
	})
}
