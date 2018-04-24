package main

import (
	"fmt"
	"strings"
)

const (
	digits        = "0123456789"
	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetFull  = alphabetLower + alphabetUpper

	underscore    = "_"
	alphaNumerics = underscore + digits + alphabetFull

	whitespace = " \t\r\n"
)

// Simple Parse Tree for Json Grammar:

// root := json
// json := object | array | value
// object := '{' pair { ',' pair } '}' | '{' '}'
// pair := value ':' json
// array := '[' json { ',' json } ']' | '[' ']'

// value := ident | string | number | booelan
// ident := [a-zA-Z_][a-zA-Z0-9_]*
// string := '\"' ([^\"\\\n\r]* | '\\' . ) '\"'
// number := [1-9][0-9]*
// boolean := "true" | "false"

type eventType uint32

const (
	keyPress eventType = iota
	keyRepeat
	keyRelease
)

type event struct {
	typ eventType
	val rune
}

// node represents a node in the program's parse tree.
type node interface {
	init(event) error
	handle(event) error
}

// rootNode ...
type rootNode struct {
	val node
}

var _ node = (*rootNode)(nil)

func (n *rootNode) init(e event) error {
	n.val = &jsonNode{n, nil}
	return n.val.init(e)
}

func (n *rootNode) handle(e event) error {
	if n.val == nil {
		n.val = &jsonNode{n, nil}
	}
	return n.val.handle(e)
}

// jsonNode ...
type jsonNode struct {
	parent node
	val    node
}

var _ node = (*jsonNode)(nil)

func (n *jsonNode) init(e event) error {
	switch {
	case e.val == '{':
		n.val = &objectNode{n, nil}
		return n.val.init(e)
	case e.val == '[':
		n.val = &arrayNode{n, nil}
		return n.val.init(e)
	case isAlphaNumeric(e.val):
		n.val = &valueNode{n, nil}
		return n.val.init(e)
	default:
		return nil
	}
}

func (n *jsonNode) handle(e event) error {
	// TODO: Handle Case of Child Existing
	if n.val != nil {
		return nil
	}

	switch {
	case isWhitespace(e.val): // Ignore
		return nil
	case e.val == '{': // Object
		n.val = &objectNode{n, nil}
		return n.val.handle(e)
	case e.val == '[': // Array
		n.val = &arrayNode{n, nil}
		return n.val.handle(e)
	case isAlphaNumeric(e.val), e.val == '"': // Value Node
		n.val = &valueNode{n, nil}
		return n.val.handle(e)
	default:
		// TODO: Return Formatted 'Invalid Key / Action' Error Messsage
		return nil
	}
}

// objectNode ...
type objectNode struct {
	parent node
	pairs  []*pairNode
}

var _ node = (*objectNode)(nil)

func (n *objectNode) init(e event) error {
	if e.val != '{' {
		return fmt.Errorf("illegal state: unexpected character '%q' in object initializer", e.val)
	}

	n.pairs = []*pairNode{}
	return nil
}

func (n *objectNode) handle(e event) error {
	switch {
	case isWhitespace(e.val): // Ignore
		return nil
	case isAlphaNumeric(e.val): // Value
		p := &pairNode{n, nil, nil}
		n.pairs = append(n.pairs, p)
		return p.handle(e)
	default:
		return nil
	}
}

// pairNode ...
type pairNode struct {
	parent node
	key    *valueNode
	val    *jsonNode
}

var _ node = (*pairNode)(nil)

func (n *pairNode) init(e event) error {
	n.key = &valueNode{n, nil}
	if err := n.key.init(e); err != nil {
		return err
	}
	return nil
}

func (n *pairNode) handle(e event) error {
	switch {
	case e.val == 'k':
		return n.key.handle(e)
	case isAlphaNumeric(e.val):
		if n.val == nil {
			n.val = &jsonNode{n, nil}
			return n.val.init(e)
		}
		return n.val.handle(e)
	default:
		return nil
	}
}

// arrayNode ...
type arrayNode struct {
	parent node
	vals   []*jsonNode
}

var _ node = (*arrayNode)(nil)

func (n *arrayNode) init(e event) error {
	if e.val != '[' {
		return fmt.Errorf("unexpected character in array initialization: '%q'", e.val)
	}
	n.vals = []*jsonNode{}
	return nil
}

func (n *arrayNode) handle(e event) error {
	switch {
	case e.val == ']': // End of Array
		return n.parent.handle(e)
	case e.val == '{', e.val == '[', e.val == '"', isAlphaNumeric(e.val):
		v := &jsonNode{n, nil}
		n.vals = append(n.vals, v)
		return v.init(e)
	case isWhitespace(e.val):
		return nil
	default:
		return nil
	}
}

// valueNode ...
type valueNode struct {
	parent node
	val    node
}

var _ node = (*valueNode)(nil)

func (n *valueNode) init(e event) error {
	switch {
	case isDigit(e.val):
		n.val = &numberNode{n, nil}
		return n.val.init(e)
	case isChar(e.val):
		n.val = &identNode{n, nil}
		return n.val.init(e)
	case e.val == '"':
		n.val = &stringNode{n, nil}
		return n.val.init(e)
	default:
		return fmt.Errorf("")
	}
}

func (n *valueNode) handle(e event) error {
	switch {
	case e.val == '"': // String
		n.val = &stringNode{n, nil}
		return n.val.handle(e)
	case isDigit(e.val): // Number
		n.val = &numberNode{n, nil}
		return n.val.handle(e)
	case isWhitespace(e.val): // Ignore
		return nil
	case isChar(e.val): // Ident
		n.val = &identNode{n, nil}
		return n.val.handle(e)
	default:
		return nil
	}
}

// identNode ...
type identNode struct {
	parent node
	val    *editNode
}

var _ node = (*identNode)(nil)

func (n *identNode) init(e event) error {
	if !isChar(e.val) && e.val != '_' {
		return fmt.Errorf("unexpected character in identifier: '%q'", e.val)
	}

	n.val = &editNode{n, nil, 0}
	return n.val.init(e)
}

func (n *identNode) handle(e event) error {
	// TODO: Write remainder of identifier editing commands
	switch {
	case isAlphaNumeric(e.val):
		return n.val.handle(e)
	case isWhitespace(e.val):
		return nil
	default:
		return nil
	}
}

// stringNode ...
type stringNode struct {
	parent node
	val    *editNode
}

var _ node = (*stringNode)(nil)

func (n *stringNode) init(e event) error {
	if e.val != '"' {
		return fmt.Errorf("unexpected start of string: '%q'", e.val)
	}
	n.val = &editNode{n, nil, 0}
	return n.val.init(e)
}

func (n *stringNode) handle(e event) error {
	switch {
	case isAlphaNumeric(e.val), isWhitespace(e.val):
		return n.val.handle(e)
	default:
		return nil
	}
}

// numberNode ...
type numberNode struct {
	parent node
	val    *editNode
}

var _ node = (*numberNode)(nil)

func (n *numberNode) init(e event) error {
	if !isDigit(e.val) {
		return fmt.Errorf("unexpected start of number: '%q'", e.val)
	}
	n.val = &editNode{n, nil, 0}
	return n.val.init(e)
}

func (n *numberNode) handle(e event) error {
	switch {
	case isDigit(e.val):
		return n.val.handle(e)
	default:
		return nil
	}
}

// booleanNode ...
type booleanNode struct {
	parent node
	val    *editNode
}

var _ node = (*booleanNode)(nil)

func (n *booleanNode) init(e event) error {
	n.val = &editNode{n, nil, 0}
	return n.val.init(e)
}

func (n *booleanNode) handle(e event) error {
	switch {
	default:
		return nil
	}
}

// editNode is a parse tree node that allows for user editing
type editNode struct {
	parent node
	val    []rune
	pos    uint
}

var _ node = (*editNode)(nil)

func (n *editNode) init(e event) error {
	n.val = []rune{}
	n.pos = 0

	switch {
	case isAlphaNumeric(e.val), isWhitespace(e.val):
		n.val = append(n.val, e.val)
		return nil
	default:
		return nil
	}
}

func (n *editNode) handle(e event) error {
	// Reset Position if in invalid location
	if n.pos >= uint(len(n.val)) {
		n.pos = uint(len(n.val)) - 1
	}

	switch {
	case isAlphaNumeric(e.val), isWhitespace(e.val):
		// Insert Character into Buffer
		start, end := n.val[:n.pos], n.val[n.pos:]
		combined := append(start, e.val)
		n.val = append(combined, end...)
		n.pos++
		return nil
	default:
		return nil
	}
}

// Uitility Functions

func contains(s string, r rune) bool {
	return strings.ContainsRune(s, r)
}

func isWhitespace(r rune) bool {
	return contains(whitespace, r)
}

func isChar(r rune) bool {
	return contains(alphabetFull, r)
}

func isDigit(r rune) bool {
	return contains(digits, r)
}

func isAlphaNumeric(r rune) bool {
	return contains(alphaNumerics, r)
}

// Main
func main() {
	fmt.Println("Hello, World!")
}
