package lib

import (
	"math"
	"testing"
)

func TestSquare(t *testing.T) {
	if Square(2) != 4 {
		t.Error("Expected: 4")
	}
}

func TestArea_Rectangle(t *testing.T) {
	rectangle := Rectangle{Point{0, 0}, Point{2, 2}}
	if rectangle.Area() != 4 {
		t.Error("Expected '4', Got: ", rectangle.Area())
	}
}

func TestArea_Circle(t *testing.T) {
	circle := Circle{0.0, 0.0, 1.0}
	if circle.Area() != math.Pi*1.0 {
		t.Error("Incorrect area: ", circle.Area())
	}
}

// func TestArea_Triangle(t *testing.T) {
// 	triangle := Triangle{
// 		Point{0, 0},
// 		Point{0, 1},
// 		Point{1, 0},
// 	}
// 	if triangle.Area() != 0.5 {
// 		t.Error("Incorrect area: ", triangle.Area())
// 	}
// }
