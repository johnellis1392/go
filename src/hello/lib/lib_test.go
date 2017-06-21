package lib

import (
  "testing"
)

func TestSquare (t *testing.T) {
  if Square(2) != 4 {
    t.Error("Expected: 4")
  }
}


