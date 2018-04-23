package main

type eventType uint

const (
	tick eventType = iota
	quit
	keyPress
	keyRepeat
	keyRelease
)

// event represents an event that can be created by this system
type event struct {
	typ eventType
	val rune
}
