package main

import (
	"fmt"
	"hello/lib"
	"strings"
)

func main() {
	r := lib.Rectangle{lib.Point{0, 0}, lib.Point{10, 10}}
	c := lib.Circle{0, 0, 5}

	fmt.Println(r.Area())
	fmt.Println(c.Area())

	s := strings.Split("123;456", ";")
	x1, x2 := s[0], s[1]
	fmt.Println(x1)
	fmt.Println(x2)
}
