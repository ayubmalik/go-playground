package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Result struct {
	Postcode Postcode    `json:"postcode"`
	Coords   Coordinates `json:"coords"`
}

type Postcode string

func (p Postcode) Normalise() Postcode {
	return Postcode(strings.ReplaceAll(string(p), " ", ""))
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func main() {
	dbUser := os.Getenv("DB_USER")
	if err := run(dbUser); err != nil {
		log.Fatal(err)
	}
}

func run(dbUser string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /latlng/{postcode}", handleGetLatLng())
	return http.ListenAndServe(":8080", mux)
}

func handleGetLatLng() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		postcode, err := url.QueryUnescape(r.PathValue("postcode"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := Result{
			Postcode: Postcode(postcode).Normalise(),
			Coords:   Coordinates{1.0, 2.0},
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
