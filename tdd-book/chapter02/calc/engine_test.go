package calc_test

import (
	"testing"

	"github.com/tdd-book/chapter02/calc"
)

func TestAdd(t *testing.T) {
	// Arrange
	e := calc.Engine{}
	x, y := 2.5, 3.5
	want := 6.0

	// Act
	got := e.Add(x, y)

	// Assert
	if got != want {
		t.Errorf("Add(%.2f,%.2f) incorrect, got: %.2f, want:%.2f; ", x, y, got, want)
	}
}
