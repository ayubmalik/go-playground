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

	t.Run("positive input", func(t *testing.T) {
		x, y := 2.5, 3.5
		want := 6.0

		got := e.Add(x, y)

		if got != want {
			t.Errorf("Add(%.2f, %.2f) incorrect, got: %.2f, want:%.2f; ", x, y, got, want)
		}
	})

	t.Run("negative input", func(t *testing.T) {
		x, y := -2.5, -3.5
		want := -6.0

		got := e.Add(x, y)

		if got != want {
			t.Errorf("Add(%.2f, %.2f) incorrect, got: %.2f, want:%.2f; ", x, y, got, want)
		}
	})
}

func BenchmarkEngine_Add(b *testing.B) {
	e := calc.Engine{}
	for i := 0; i < b.N; i++ {
		e.Add(2, 3)
	}
}

func BenchmarkEngine_Double(b *testing.B) {
	e := calc.Engine{}

	b.Run("double", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			e.Double(2.5)
		}
	})

	b.Run("double2", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			e.Double(2.5)
		}
	})

}
