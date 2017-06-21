package lib

import (
	"math"
)

// lib/ directory gets ignored automatically when go packages a build

// Shape ...
type Shape interface {
	Area() float64
}

// Point ...
type Point struct {
	X, Y float64
}

// Triangle ...
type Triangle struct {
	P1, P2, P3 Point
}

// Circle ...
type Circle struct {
	X, Y, R float64
}

// Rectangle ...
type Rectangle struct {
	P1, P2 Point
}

func distance(p1, p2 Point) float64 {
	a := p1.X - p2.X
	b := p1.Y - p2.Y
	return math.Sqrt(a*a + b*b)
}

// Area ...
func (r *Rectangle) Area() float64 {
	l := distance(r.P1, r.P2)
	w := distance(r.P1, r.P2)
	return l * w
}

// Area ...
func (c *Circle) Area() float64 {
	return math.Pi * c.R * c.R
}

// Area ...
func (t *Triangle) Area() float64 {
	d1 := distance(t.P1, t.P2)
	d2 := distance(t.P2, t.P3)
	return d1 * d2 / 2.0
}

// Square ...
func Square(i int) int {
	return i * i
}