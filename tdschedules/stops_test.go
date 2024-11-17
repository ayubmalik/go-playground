package tdschedules_test

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestStops(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file")
	}

	ctx := context.Background()
	connStr := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			t.Fatalf("Error closing connection: %v", err)
		}
	}(conn, ctx)

	if err := conn.Ping(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

}
