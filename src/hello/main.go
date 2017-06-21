package main

import (
	"fmt"
	"hello/lib"
	"strings"
)

func main() {
	r := lib.Rectangle{P1: lib.Point{X: 0, Y: 0}, P2: lib.Point{X: 10, Y: 10}}
	c := lib.Circle{X: 0, Y: 0, R: 5}

	fmt.Println(r.Area())
	fmt.Println(c.Area())

	s := strings.Split("123;456", ";")
	x1, x2 := s[0], s[1]
	fmt.Println(x1)
	fmt.Println(x2)
}
