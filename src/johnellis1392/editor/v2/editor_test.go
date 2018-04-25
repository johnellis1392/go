package main

import (
	"testing"
)

func Test_Editor_CreatesNewElement_WhenInsertPressed(t *testing.T) {
	e := newEditor()
	ev := event{keyPress, 'i'}

	if err := e.handle(ev); err != nil {
		t.Error(err)
	}

	// Verify New State
}
