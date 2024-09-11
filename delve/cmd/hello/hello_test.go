package main

import (
	"fmt"
	"testing"
)

func TestHello(t *testing.T) {
	t.Run("print hello", func(t *testing.T) {
		fmt.Println("hello")
	})

	t.Run("test reverse", func(t *testing.T) {
		months := []string{"jan", "feb", "mar"}
		got := months[0]
		want := "naj"

		if got != want {
			t.Errorf("months[0] = %v, want %v", want, got)
		}
	})
}
