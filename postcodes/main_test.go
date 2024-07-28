package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockFinder func(p Postcode) *Coords

func (mf mockFinder) Find(p Postcode) *Coords {
	return mf(p)
}

func Test_handleGetLatLng(t *testing.T) {
	mf := mockFinder(func(p Postcode) *Coords {
		return &Coords{
			Lat: 123,
			Lng: 456,
		}
	})

	handler := handleGetLatLng(mf)

	t.Run("valid postcode result", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetPathValue("postcode", "BD72AN")

		want := CoordsResult{
			Postcode: "BD72AN",
			Coords: Coords{
				Lat: 123.00,
				Lng: 456.00,
			},
		}

		handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
		}

		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("handler returned wrong content type: got %v want %v", w.Header().Get("Content-Type"), "application/json")
		}

		got := CoordsResult{}
		_ = json.NewDecoder(w.Body).Decode(&got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("handler returned wrong result: got %v want %v", got, want)
		}

	})
}
