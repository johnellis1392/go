package main

import (
	"testing"
)

func arrayEquals(a1, a2 []item) bool {
	if a1 == nil && a2 == nil {
		return true
	} else if a1 == nil || a2 == nil {
		return false
	} else if len(a1) != len(a2) {
		return false
	}

	for i := range a1 {
		item1 := a1[i]
		item2 := a2[i]
		if !item1.equals(item2) {
			return false
		}
	}

	return true
}

func Test_LexerBasic(t *testing.T) {
	t.Skip()

	// Start Lexing
	l, itemchan := lex("ExampleLexer", `Something {{1}} Else`)
	if l == nil || itemchan == nil {
		t.Errorf("Unexpected return values from call to lex(a, b string):\n * lexer: %v\n * chan item: %v", l, itemchan)
	}

	// Collect Items
	var items []item
outer:
	for {
		select {
		case i, ok := <-itemchan:
			items = append(items, i)
			if ok {
				continue
			} else if i.typ == itemEOF {
				// Success
				break outer
			} else {
				// Failure
				t.Errorf("Error: Reached end of item channel without EOF")
			}
		}
	}

	// Verify
	expected := []item{}
	if !arrayEquals(expected, items) {
		t.Errorf("Unexpected Token Stream Values\n * Expected: %v\n * Actual: %v", expected, items)
	}
}
