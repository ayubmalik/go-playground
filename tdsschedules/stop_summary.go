package tdsschedules

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type StopSummary struct {
	ID    string
	Name  string
	Code  string
	City  string
	State string
}

func NewStopSummaryDB(conn *pgx.Conn) *StopSummaryDB {
	return &StopSummaryDB{conn}
}

type StopSummaryDB struct {
	conn *pgx.Conn
}

func (db *StopSummaryDB) GetAll(ctx context.Context) ([]StopSummary, error) {
	query := `SELECT id, station_name, station_code, city_name, state_code FROM stop_summary`

	rows, err := db.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query stop summaries %w", err)
	}
	defer rows.Close()

	// see also pgx.CollectRows())
	stops := make([]StopSummary, 0)
	for rows.Next() {
		stop := StopSummary{}
		err := rows.Scan(&stop.ID, &stop.Name, &stop.Code, &stop.City, &stop.State)
		if err != nil {
			return nil, fmt.Errorf("could not scan stop summary %w", err)
		}
		stops = append(stops, stop)
	}

	return stops, nil
}

func (db *StopSummaryDB) Put(ctx context.Context, stop StopSummary) error {
	query := `INSERT INTO stop_summary(id, station_name, station_code, city_name, state_code)
			  VALUES (@id, @station_name, @station_code, @city_name, @state_code)`
	t, err := db.conn.Exec(ctx, query, pgx.NamedArgs{
		"id":           stop.ID,
		"station_name": stop.Name,
		"station_code": stop.Code,
		"city_name":    stop.City,
		"state_code":   stop.State,
	})
	if err != nil {
		return fmt.Errorf("could not insert stop summary: %w", err)
	}
	slog.Debug("put stop", "count", t.RowsAffected())
	return nil
}

func (db *StopSummaryDB) Get(ctx context.Context, uuid string) (StopSummary, error) {
	query := `SELECT id, station_name, station_code, city_name, state_code 
			  FROM stop_summary
			  WHERE id = @id
			  `
	row := db.conn.QueryRow(ctx, query, pgx.NamedArgs{
		"id": uuid,
	})

	var stop StopSummary
	err := row.Scan(&stop.ID, &stop.Name, &stop.Code, &stop.City, &stop.State)

	return stop, err
}

func (db *StopSummaryDB) Delete(ctx context.Context, uuid string) error {
	query := `DELETE FROM stop_summary WHERE id = $1`

	t, err := db.conn.Exec(ctx, query, uuid)
	if err != nil {
		return fmt.Errorf("could not delete stop summary %w", err)
	}
	slog.Debug("delete stop", "count", t.RowsAffected())

	return nil
}
