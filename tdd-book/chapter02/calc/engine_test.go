package calc_test

import (
	"log"
	"os"
	"testing"

	"github.com/tdd-book/chapter02/calc"
)

func TestMain(m *testing.M) {
	setup()

	err := m.Run()

	teardown()

	os.Exit(err)
}

func init() {
	log.Println("init")
}

func setup() {
	log.Println("setup")
}

func teardown() {
	log.Println("teardown")
}

func TestAdd(t *testing.T) {
	defer func() {
		log.Println("deferred teardown")
	}()

	// Arrange
	e := calc.Engine{}
	x, y := 2.5, 3.5
	want := 6.0

	// Act
	got := e.Add(x, y)

	// Assert
	if got != want {
		t.Errorf("Add(%.2f, %.2f) incorrect, got: %.2f, want:%.2f; ", x, y, got, want)
	}
}
