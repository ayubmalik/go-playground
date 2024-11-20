package tdsschedules_test

import (
	"tdsschedules"
	"testing"
	"time"
)

func TestNextMonday(t *testing.T) {
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
		{"Tue", "2024-10-08", "2024-10-14"},
		{"Wed", "2024-10-09", "2024-10-14"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			date, _ := time.Parse(time.DateOnly, test.dt)
			want, _ := time.Parse(time.DateOnly, test.want)

			got := tdsschedules.NextMonday(date)
			if got != want {
				t.Errorf("nextMonday for %s = %s, want %s", date, got, test.want)
			}
		})
	}
}
