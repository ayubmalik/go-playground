package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TdsRestApi struct {
	client    http.Client
	url       string
	key       string
	carrierId int
}

type StopsQuery struct {
	Type      string
	CarrierId int
}

func (tds TdsRestApi) Origins() ([]StopCity, error) {
	url := tds.url + "/stop"
	qry := StopsQuery{
		CarrierId: tds.carrierId,
		Type:      "ORIGIN",
	}

	buf, err := json.Marshal(qry)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("tds-api-key", tds.key)

	res, err := tds.client.Do(req)
	fmt.Println("STATUS", res.Status)
	if err != nil {
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
