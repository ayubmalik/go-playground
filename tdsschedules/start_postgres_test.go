package tdsschedules_test

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

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
