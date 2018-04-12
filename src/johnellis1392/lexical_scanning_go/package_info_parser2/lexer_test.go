package main

import (
	"testing"
)

func TestLexStringRemovesQuotes(t *testing.T) {
	// t.SkipNow()

	input := "\"something\""
	l := newLexer(input)
	l.state = lexString

	go func() {
		_ = l.state(l)
	}()

	tok := <-l.output
	if tok.typ != tokenString {
		t.Errorf("unexpected type: expected tokenString, found '%q'", tok)
	}

	// Check Quotes are Gone
	if tok.val != "something" {
		t.Errorf("improperly escaped string literal: %q", tok.val)
	}
}
