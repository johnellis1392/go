package main

import (
  "fmt"
)

type signal float64
type midi struct {
  bend uint16 // Pitch Bend
  notes []uint16
}

type transform func(signal) signal
type midiTransform func(midi) midi

type audioSource struct{}
type midiSource struct{}

type app struct {
  sources []audioSource
}

func main() {
  fmt.Println("Starting application...")
}
