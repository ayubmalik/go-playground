package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello")
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func Reverse2(s string) string {
	i := 0
	r := []rune(s)
	rr := make([]rune, len(r))
	for j := len(r) - 1; j >= 0; j-- {
		rr[i] = r[j]
		i++
	}
	return string(rr)
}
