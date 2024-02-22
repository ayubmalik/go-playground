package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type TdsWebClient struct {
	resty   *resty.Client
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
	if err != nil {
		return err
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

func (c TdsWebClient) Origins2() error {
	url := c.url + "/stop"

	resp, err := c.resty.R().
		EnableTrace().
		SetHeader("TDS-Api-Key", c.key).
		SetHeader("content-type", "application/json").
		SetHeader("accept", "application/json").
		SetHeader("User-Agent", "curl/7.88.1").
		SetHeader("tds-api-key", c.key).
		SetBody(`{"carrierId":304,"type":"ORIGIN"}`).
		Post(url)

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())

	return nil
}
