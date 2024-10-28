package main

import (
	"context"
	"testing"
)

type stubClient struct {
	count       int
	origin      string
	destination string
}

func (s *stubClient) FindSchedules(ctx context.Context, qry ScheduleQuery) (ScheduleResult, error) {
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
