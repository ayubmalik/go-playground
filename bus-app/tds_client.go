package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type TdsRestApi struct {
	client    http.Client
	url       string
	key       string
	carrierId int
}

type StopsQuery struct {
	Type      string `json:"type"`
	CarrierId int    `json:"carrierId"`
}

const (
	ContentType     = "content-type"
	Accept          = "accept"
	ApplicationJson = "application/json"
	TdsApiKey       = "tds-api-key"
)

func (tds TdsRestApi) Origins() ([]StopCity, error) {
	url := tds.url + "/stop"
	qry := StopsQuery{
		Type:      "ORIGIN",
		CarrierId: tds.carrierId,
	}

	buf, err := json.Marshal(qry)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Add(ContentType, ApplicationJson)
	req.Header.Add(Accept, ApplicationJson)
	req.Header.Add("tds-api-key", tds.key)

	res, err := tds.client.Do(req)
	if err != nil {
		log.Println("STATUS", res.Status)
		return nil, err
	}

	defer res.Body.Close()

	var stopCities []StopCity

	err = json.NewDecoder(res.Body).Decode(&stopCities)
	if err != nil {
		return nil, err
	}

	return stopCities, nil
}
