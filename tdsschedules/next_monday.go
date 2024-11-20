package tdsschedules

import (
	"time"
)

func NextMonday(date time.Time) time.Time {
	offset := (8 - date.Weekday()) % 7
	if offset == 0 {
		offset = 7
	}
	return date.Add(24 * time.Hour * time.Duration(offset))
}
