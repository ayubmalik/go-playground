package tdsschedules

import (
	"context"
	"fmt"
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
	query := `SELECT *  FROM origin_destination`

	rows, err := db.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query stop summaries %w", err)
	}
	defer rows.Close()

	// see also pgx.CollectRows())
	ods, err := pgx.CollectRows(rows, pgx.RowToStructByName[OriginDestination])
	return ods, err

}

func (db *OrigDestinationDB) Put(ctx context.Context, od OriginDestination) error {
	query := `INSERT INTO origin_destination(origin, destination)
			  VALUES (@origin, @destination)
			  ON CONFLICT(origin, destination) 
			  DO NOTHING`

	_, err := db.conn.Exec(ctx, query, pgx.NamedArgs{
		"origin":      od.Origin.ID,
		"destination": od.Destination.ID,
	})
	if err != nil {
		return fmt.Errorf("could not insert origin destination: %w", err)
	}
	return nil

}
