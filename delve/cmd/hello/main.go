package main

import (
	"fmt"
	"os"
)

type Stop struct {
	Origin      string
	Destination string
}

func main() {
	ARUBA := "Hello World!"

	fmt.Println(ARUBA)
	type Test struct {
		x int
		y int
	}
	Good()

	n := 1
	stop := Stop{
		Origin:      "Manchester",
		Destination: "Liverpool",
	}

	fmt.Println("Stop", stop, n)
	os.Exit(0)
}

func Good() {
	msg := "Goodbye cruel world!"

	fmt.Println(msg)
}
