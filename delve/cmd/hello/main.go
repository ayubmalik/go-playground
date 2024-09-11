package main

import (
	"fmt"
	"os"
)

func main() {
	ARUBA := "Hello World!"

	fmt.Println(ARUBA)
	type Test struct {
		x int
		y int
	}
	Good()

	os.Exit(0)
}

func Good() {
	msg := "Goodbye cruel world!"

	fmt.Println(msg)
}
