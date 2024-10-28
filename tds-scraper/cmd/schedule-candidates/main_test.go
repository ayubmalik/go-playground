package main

import "testing"

func TestTrySchedule(t *testing.T) {
	tds := createTDSClient()

	t.Run("days range", func(t *testing.T) {
		if err := trySchedule(tds, 2, "origin", "dest"); err != nil {
			t.Errorf("not trySchedule error %v", err)
		}
	})
}
