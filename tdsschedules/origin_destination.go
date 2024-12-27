package tdsschedules

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type OriginDestination struct {
	Origin      StopSummary
	Destination StopSummary
}

func NewOriginDestinationDB(conn *pgxpool.Pool) *OrigDestinationDB {
	return &OrigDestinationDB{conn}
}

type OrigDestinationDB struct {
	conn *pgxpool.Pool
}

func (db *OrigDestinationDB) GetAll(ctx context.Context) ([]OriginDestination, error) {
	query := `SELECT origin.id, origin.station_name, origin.station_code, origin.city_name, origin.state_code,
                     destination.id, destination.station_name, destination.station_code, destination.city_name, destination.state_code
              FROM origin_destination od
                   JOIN stop_summary origin ON od.origin = origin.id
                   JOIN stop_summary destination ON od.destination = destination.id`

	rows, err := db.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query stop summaries %w", err)
	}
	defer rows.Close()

	var ods []OriginDestination
	// see also pgx.CollectRows
	for rows.Next() {
		var od OriginDestination
		if err := rows.Scan(&od.Origin.ID, &od.Origin.Name, &od.Origin.Code, &od.Origin.City, &od.Origin.State,
			&od.Destination.ID, &od.Destination.Name, &od.Destination.Code, &od.Destination.City, &od.Destination.State); err != nil {
			return nil, fmt.Errorf("could not scan stop summary %w", err)
		}
		ods = append(ods, od)
	}
	return ods, err
}

func (db *OrigDestinationDB) Put(ctx context.Context, od OriginDestination) error {
	query := `INSERT INTO origin_destination(origin, destination)
			  VALUES (@origin, @destination)
			  ON CONFLICT(origin, destination) 
			  DO NOTHING`

	t, err := db.conn.Exec(ctx, query, pgx.NamedArgs{
		"origin":      od.Origin.ID,
		"destination": od.Destination.ID,
	})
	if err != nil {
		return fmt.Errorf("could not insert origin destination: %w", err)
	}
	slog.Debug("put origin_destination", "count", t.RowsAffected())
	return nil
}

// TODO(test)
func (db *OrigDestinationDB) DeleteAll(ctx context.Context) error {
	query := `DELETE FROM origin_destination WHERE TRUE`
	t, err := db.conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("could not delete all origin destinations: %w", err)
	}

	slog.Debug("delete all origin_destinations", "count", t.RowsAffected())
	return nil
}
