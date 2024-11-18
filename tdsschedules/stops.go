package tdschedules

import (
	"github.com/jackc/pgx/v5"
)

func NewStopDB(conn *pgx.Conn) *StopDB {
	return &StopDB{conn}
}

type StopDB struct {
	conn *pgx.Conn
}
