package main

import (
	"encoding/json"
	"fmt"
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
}

func (r Repo) Find(p Postcode) *Coords {
	return nil
}

func main() {
	repo := Repo{}
	if err := run(repo); err != nil {
		log.Fatal(err)
	}
}

func run(finder CoordsFinder) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /latlng/{postcode}", handleGetLatLng(finder))
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
		if coords == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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
