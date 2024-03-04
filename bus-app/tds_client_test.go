package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMama(t *testing.T) {
	left := "hello"
	right := "world!"

	msg := left + " " + right
	assert.Equal(t, "hello world!", msg)
}

func TestOrigins(t *testing.T) {
	var headers map[string][]string

	jsonResponse := `[]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers = r.Header
		w.Write([]byte(jsonResponse))
	}))
	defer server.Close()

	tds := TdsRestApi{
		client:    http.Client{},
		key:       "some api key",
		carrierId: 777,
		url:       server.URL,
	}

	_, err := tds.Origins()
	assert.Nil(t, err, "error was not nil")
	assert.Equal(t, "some api key", headers["Tds-Api-Key"][0])
}
