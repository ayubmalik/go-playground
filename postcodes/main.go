package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Postcode string

func (p Postcode) Normalise() Postcode {
	return Postcode(strings.ReplaceAll(string(p), " ", ""))
}

type Coords struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type CoordsFinder interface {
	Find(p Postcode) *Coords
}

type CoordsResult struct {
	Postcode Postcode `json:"postcode"`
	Coords   Coords   `json:"coords"`
}

type Repo struct {
	db *pgx.Conn
}

func (r Repo) Find(p Postcode) *Coords {
	qry := "SELECT lat, lng FROM postcode_geo WHERE postcode = $1"
	var coords Coords
	row := r.db.QueryRow(context.Background(), qry, p)
	err := row.Scan(&coords.Lat, &coords.Lng)
	if err != nil {
		return nil
	}
	return &coords
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	url := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	err = conn.Ping(ctx)
	if err != nil {
		return err
	}
	repo := Repo{db: conn}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /latlng/{postcode}", handleGetLatLng(repo))
	return http.ListenAndServe(":8080", mux)
}

func handleGetLatLng(finder CoordsFinder) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		param, err := url.QueryUnescape(r.PathValue("postcode"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		postcode := Postcode(param)
		coords := finder.Find(postcode.Normalise())
		fmt.Printf("coords: %+v\n", coords)
		if coords == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		result := CoordsResult{
			Postcode: postcode,
			Coords:   *coords,
		}

		fmt.Printf("result: %+v\n", result)
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
