package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

type TdsWebClient struct {
	url     string
	key     string
	carrier string
}

func (c TdsWebClient) Origins() error {
	u := c.url + "/stop"
	fmt.Println("URL", u)
	body := []byte(`{"carrierId":304,"type":"ORIGIN"}`)

	buf := bytes.NewBuffer(body)
	req, err := http.NewRequest("POST", "https://ride-api.bustickets.com/tickets/stop", buf)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "curl/7.88.1")
	req.Header.Add("tds-api-key", c.key)

	hc := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http2.Transport{
			ReadIdleTimeout:    0,
			MaxReadFrameSize:   0,
			DisableCompression: true,
			AllowHTTP:          false,
		}}

	res, err := hc.Do(req)
	if err != nil {
		return err
	}

	res, err = hc.Do(req)
	fmt.Println("STATUS", res.Status)
	if err != nil {
		return err
	}
	fmt.Println("key", c.key)
	defer res.Body.Close()

	var stopCities []StopCity

	err = json.NewDecoder(res.Body).Decode(&stopCities)
	if err != nil {
		return err
	}

	for _, s := range stopCities {
		fmt.Println("Got stop", s)
	}
	return err

}
