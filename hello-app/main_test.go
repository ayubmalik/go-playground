package main

import "testing"

func TestBuildProcess(t *testing.T) {
	t.Run("success 1", func(t *testing.T) {
		t.Log("yay 1")
	})

	t.Run("failure", func(t *testing.T) {
		// t.Error("oh no")
	})

	t.Run("success 2", func(t *testing.T) {
		t.Log("yay 2")
	})
}
