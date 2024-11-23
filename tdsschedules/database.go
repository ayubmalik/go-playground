package tdsschedules

import (
	"github.com/jackc/pgx/v5"
)

type StopSummary struct {
	ID          string
	StationName string
	StationCode string
	CityName    string
	StateCode   string
}

func NewStopSummaryDB(conn *pgx.Conn) *StopSummaryDB {
	return &StopSummaryDB{conn}
}

type StopSummaryDB struct {
	conn *pgx.Conn
}

func (db *StopSummaryDB) GetAll() ([]StopSummary, error) {

	return []StopSummary{{
		ID:          "",
		StationName: "",
		StationCode: "",
		CityName:    "",
		StateCode:   "",
	}}, nil
}

func (db *StopSummaryDB) Put(summary StopSummary) error {
	
}
