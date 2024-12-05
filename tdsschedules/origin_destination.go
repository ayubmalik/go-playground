package tdsschedules

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type OriginDestination struct {
	Origin      StopSummary
	Destination StopSummary
}

func NewOriginDestinationDB(conn *pgx.Conn) *OrigDestinationDB {
	return &OrigDestinationDB{conn}
}

type OrigDestinationDB struct {
	conn *pgx.Conn
}

func (db *OrigDestinationDB) GetAll(ctx context.Context) ([]OriginDestination, error) {

	return nil, nil
}
