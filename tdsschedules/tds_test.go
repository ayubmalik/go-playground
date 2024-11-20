package tdsschedules_test

import (
	"testing"
	"time"
)

func TestNewScheduleRequest(t *testing.T) {
	//t.Run("new schedule request", func(t *testing.T) {
	//	req := newScheduleQuery("orig1", "dest1", "2025-12-25")
	//
	//	if req.PurchaseType != "SCHEDULE_BOOK" {
	//		t.Errorf("req.purchaseType = %s, want %s", req.PurchaseType, "SCHEDULE_BOOK")
	//	}
	//
	//	j, _ := json.Marshal(req)
	//	t.Log(string(j))
	//
	//})
}

// Only example how to write a range
func doRange(days int) func(yield func(date string) bool) {
	return func(yield func(date string) bool) {
		for i := range days {
			dt := time.Now().Add(24 * time.Hour * time.Duration(i+1)).Format("2006-01-02")
			if !yield(dt) {
				return
			}
		}
	}
}

func TestDoRange(t *testing.T) {
	t.Skip()
	doRange(1)
}
