package calculator

import "fmt"

type Adder interface {
	Add(x, y float64) float64
}

func NewCalculator(a Adder) *Calculator {
	return &Calculator{adder: a}
}

type Calculator struct {
	adder Adder
}

func (c *Calculator) PrintAdd(x, y float64) {
	fmt.Println("result = ", c.adder.Add(x, y))
}
