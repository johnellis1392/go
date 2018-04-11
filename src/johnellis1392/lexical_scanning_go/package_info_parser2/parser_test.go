package main

import (
	"testing"
)

const (
	input1 = `base = { };`
)

func TestMatch(t *testing.T) {
	input := make(chan token, 3)
	p := newParser(input)

	go func() {
		input <- token{tokenIdent, "a"}
		input <- token{tokenEquals, "="}
		input <- token{tokenNumber, "3"}
	}()

	p.shift()
	p.shift()
	p.shift()

	if m := p.match(nodeIdent, nodeEquals, nodeNumber); !m {
		t.Errorf("invalid match: %v", m)
	}
}
