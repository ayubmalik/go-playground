package main

import (
	"log"
	"time"
)

type stop struct {
	StopUuid string
}

type ScheduleRequest struct {
	PassengerCounts map[string]int
	PurchaseType    string
	Origin          stop
	Destination     stop
	DepartureDate   string
}

type ScheduleResponse struct{}

type ScheduleAPI interface {
	Get(req ScheduleRequest) (ScheduleResponse, error)
}

func trySchedule(schedules ScheduleAPI, days int, origin, dest string) error {
	var dates []string
	for i := range days {
		dt := time.Now().Add(time.Duration(i+1) * 24 * time.Hour)
		dates = append(dates, dt.Format("2006-01-02"))
	}
	log.Printf("dates: %v", dates)
	return nil
}
