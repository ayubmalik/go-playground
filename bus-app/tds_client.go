package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TdsClient struct {
	url       string
	key       string
	carrierId int
}

type StopsQuery struct {
	Type      string
	CarrierId int
}

func (c TdsClient) Origins() ([]StopCity, error) {
	url := c.url + "/stop"
	qry := StopsQuery{
		CarrierId: c.carrierId,
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
	req.Header.Add("User-Agent", "curl/7.88.1")
	req.Header.Add("tds-api-key", c.key)

	hc := &http.Client{
		Timeout: 60 * time.Second,
	}

	res, err := hc.Do(req)
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

	for _, s := range stopCities {
		fmt.Println("Got stop", s)
	}
	return stopCities, nil
}
