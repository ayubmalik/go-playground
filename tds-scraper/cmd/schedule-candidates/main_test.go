package main

import "testing"

type stubClient struct {
	count int
}

func (s stubClient) FindSchedules(qry ScheduleQuery) (ScheduleResult, error) {
	s.count++
}

func TestTrySchedule(t *testing.T) {

	tds := stubClient{}
	t.Run("days range", func(t *testing.T) {
		if err := trySchedule(tds, 2, "origin", "dest"); err != nil {
			t.Errorf("not trySchedule error %v", err)
		}
	})
}
