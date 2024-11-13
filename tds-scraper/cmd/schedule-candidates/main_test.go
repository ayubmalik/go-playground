package main

import (
	"context"
	"testing"
	"time"
)

type stubClient struct {
	count       int
	origin      string
	destination string
}

func (s *stubClient) FindSchedules(_ context.Context, qry ScheduleQuery) (ScheduleResult, error) {
	s.count++
	s.origin = qry.Origin.StopUuid
	s.destination = qry.Destination.StopUuid
	return ScheduleResult{}, nil
}

func TestTrySchedule(t *testing.T) {
	tds := &stubClient{}
	t.Run("days range", func(t *testing.T) {
		_ = trySchedule(tds, 7, "origin", "dest")
		if tds.count != 7 {
			t.Errorf("client calls = %d, want 2", tds.count)
		}

		if tds.origin != "origin" {
			t.Errorf("origin = %s, want origin", tds.origin)
		}

		if tds.destination != "dest" {
			t.Errorf("destination = %s, want dest", tds.destination)
		}
	})
}

func TestDays(t *testing.T) {
	tests := []struct {
		name string
		dt   string
		want string
	}{
		{"Tue", "2024-10-01", "2024-10-07"},
		{"Wed", "2024-10-02", "2024-10-07"},
		{"Thu", "2024-10-03", "2024-10-07"},
		{"Fri", "2024-10-04", "2024-10-07"},
		{"Sat", "2024-10-05", "2024-10-07"},
		{"Sun", "2024-10-06", "2024-10-07"},
		{"Mon", "2024-10-07", "2024-10-14"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			date, _ := time.Parse(time.DateOnly, test.dt)
			want, _ := time.Parse(time.DateOnly, test.want)

			got := nextMonday(date)
			if got != want {
				t.Errorf("nextMonday = %s, want %s", got, test.want)
			}
		})
	}
}

func nextMonday(date time.Time) time.Time {
	offset := (8 - date.Weekday()) % 7
	if offset == 0 {
		offset = 7
	}
	return date.Add(24 * time.Hour * time.Duration(offset))
}

func TestWithDayOfWeek(t *testing.T) {
	now := time.Now()

	t.Logf("now = %v", now.Format(time.RFC3339))
	t.Logf("day = %v", now.Day())
	t.Logf("week = %v", now.Weekday())

	day := 24 * time.Hour
	mon := now.Add(day)
	t.Logf("now = %v", mon.Format(time.RFC3339))
	t.Logf("day = %v", mon.Day())
	t.Logf("plus = %v", mon.Weekday())
	diff := 8 - mon.Weekday()

	inc := day * time.Duration(diff)
	mon = mon.Add(inc)
	t.Logf("diff = %v", diff)
	t.Logf("mon = %v", mon)
}
