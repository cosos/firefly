package main

import (
	"log"

	"golang.org/x/exp/constraints"
)

func GMin[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func main() {
	fmin := GMin[int]
	log.Println(fmin(2, 3))
	log.Println(fmin(5, 20))
}
