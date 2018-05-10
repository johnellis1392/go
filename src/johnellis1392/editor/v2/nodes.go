package main

import (
	"fmt"
	"io"
	"strings"
)

// node represents a node in the program's parse tree.
type node interface {
	handle(event) (node, error)
	render(io.Writer) error
	String() string
}

// rootNode ...
type rootNode struct {
	val node
}

var _ node = (*rootNode)(nil)

func (n *rootNode) handle(e event) (node, error) {
	if n.val == nil {
		n.val = &jsonNode{n, nil}
	}

	switch e.val {
	// case '\n', '\r':
	// 	return n.val.handle(e)
	// case 'i':
	// 	return n.val.handle(e)
	// case 'd':
	// 	n.val = nil
	// 	return nil, nil
	// case 'u':
	// 	return nil, nil
	// case 'r':
	// 	return nil, nil
	// case 'a':
	// 	return nil, nil
	// case 'n':
	// 	return nil, nil
	// case 'p':
	// 	return nil, nil
	// case 'o':
	// 	return nil, nil
	default:
		return n.val.handle(e)
	}
}

func (n *rootNode) render(w io.Writer) error {
	return nil
}

func (n *rootNode) String() string {
	return fmt.Sprintf("rootNode{%v}", n.val)
}

// jsonNode ...
type jsonNode struct {
	parent node
	val    node
}

var _ node = (*jsonNode)(nil)

func (n *jsonNode) handle(e event) (node, error) {
	if n.val != nil {
		return n.val.handle(e)
	}

	switch {
	case e.val == '\n', e.val == '\r':
		return n.parent, nil
	case isWhitespace(e.val): // Ignore
		return n, nil
	case e.val == '{': // Object
		n.val = &objectNode{n, nil}
		return n.val, nil
	case e.val == '[': // Array
		n.val = &arrayNode{n, nil}
		return n.val, nil
	case isAlphaNumeric(e.val), isDigit(e.val), e.val == '"': // Value Node
		n.val = &valueNode{n, nil}
		return n.val.handle(e)
	default:
		// TODO: Return Formatted 'Invalid Key / Action' Error Messsage
		return n, nil
	}
}

func (n *jsonNode) render(w io.Writer) error {
	return nil
}

func (n *jsonNode) String() string {
	return fmt.Sprintf("jsonNode{%v}", n.val)
}

// objectNode ...
type objectNode struct {
	parent node
	pairs  []*pairNode
}

var _ node = (*objectNode)(nil)

func (n *objectNode) handle(e event) (node, error) {
	switch {
	case isWhitespace(e.val): // Ignore
		return n, nil
	case isAlphaNumeric(e.val): // Value
		p := &pairNode{n, nil, nil}
		n.pairs = append(n.pairs, p)
		return p.handle(e)
	case e.val == '}': // End Object
		return n.parent, nil
	default:
		return n, nil
	}
}

func (n *objectNode) render(w io.Writer) error {
	return nil
}

func (n *objectNode) String() string {
	return fmt.Sprintf("objectNode{%v}", n.pairs)
}

// pairNode ...
type pairNode struct {
	parent node
	key    *valueNode
	val    *jsonNode
}

var _ node = (*pairNode)(nil)

func (n *pairNode) handle(e event) (node, error) {
	if n.key == nil {
		n.key = &valueNode{n, nil}
	}

	if n.val == nil {
		n.val = &jsonNode{n, nil}
	}

	switch {
	case isAlphaNumeric(e.val):
		return n.key.handle(e)
	case e.val == ':':
		return n.val, nil
	default:
		return n, nil
	}
}

func (n *pairNode) render(w io.Writer) error {
	return nil
}

func (n *pairNode) String() string {
	return fmt.Sprintf("pairNode{%v, %v}", n.key, n.val)
}

// arrayNode ...
type arrayNode struct {
	parent node
	vals   []*jsonNode
}

var _ node = (*arrayNode)(nil)

func (n *arrayNode) handle(e event) (node, error) {
	switch {
	case e.val == ']': // End of Array
		return n.parent, nil
	case e.val == '{', e.val == '[', e.val == '"', isAlphaNumeric(e.val):
		v := &jsonNode{n, nil}
		n.vals = append(n.vals, v)
		return v.handle(e)
	case isWhitespace(e.val):
		return n, nil
	default:
		return n, nil
	}
}

func (n *arrayNode) render(w io.Writer) error {
	return nil
}

func (n *arrayNode) String() string {
	return fmt.Sprintf("arrayNode{%v}", n.vals)
}

// valueNode ...
type valueNode struct {
	parent node
	val    node
}

var _ node = (*valueNode)(nil)

func (n *valueNode) handle(e event) (node, error) {
	switch {
	case isDigit(e.val): // Number
		n.val = &numberNode{n, nil}
		return n.val.handle(e)
	case isChar(e.val): // Ident
		n.val = &identNode{n, nil}
		return n.val.handle(e)
	case e.val == '"': // String
		n.val = &stringNode{n, nil}
		return n.val, nil
	case e.val == '\n', e.val == '\r', e.val == ':':
		return n.parent.handle(e)
	case isWhitespace(e.val): // Ignore
		return n, nil
	default:
		return n, nil
	}
}

func (n *valueNode) render(w io.Writer) error {
	return nil
}

func (n *valueNode) String() string {
	return fmt.Sprintf("valueNode{%v}", n.val)
}

// identNode ...
type identNode struct {
	parent node
	val    *textBuffer
}

var _ node = (*identNode)(nil)

func (n *identNode) handle(e event) (node, error) {
	if n.val == nil {
		n.val = &textBuffer{0, nil}
	}

	switch {
	case isAlphaNumeric(e.val):
		n.val.insert(e.val)
		return n, nil
	case e.val == '\n', e.val == '\r', e.val == ':':
		return n.parent.handle(e)
	case isWhitespace(e.val): // Ignore
		return n, nil
	default:
		return n, nil
	}
}

func (n *identNode) render(w io.Writer) error {
	return nil
}

func (n *identNode) String() string {
	return fmt.Sprintf("identNode{%v}", n.val)
}

// stringNode ...
type stringNode struct {
	parent node
	val    *textBuffer
}

var _ node = (*stringNode)(nil)

func (n *stringNode) handle(e event) (node, error) {
	if n.val == nil {
		n.val = &textBuffer{0, nil}
	}

	switch {
	case isAlphaNumeric(e.val), isWhitespace(e.val):
		n.val.insert(e.val)
		return n, nil
	case e.val == '"': // End of String
		return n.parent, nil
	default:
		return n, nil
	}
}

func (n *stringNode) render(w io.Writer) error {
	return nil
}

func (n *stringNode) String() string {
	return fmt.Sprintf("stringNode{%v}", n.val)
}

// numberNode ...
type numberNode struct {
	parent node
	val    *textBuffer
}

var _ node = (*numberNode)(nil)

func (n *numberNode) handle(e event) (node, error) {
	if n.val == nil {
		n.val = &textBuffer{0, nil}
	}

	switch {
	case isDigit(e.val):
		n.val.insert(e.val)
		return n, nil
	case e.val == '\n', e.val == '\r':
		return n.parent, nil
	default:
		return n, nil
	}
}

func (n *numberNode) render(w io.Writer) error {
	return nil
}

func (n *numberNode) String() string {
	return fmt.Sprintf("numberNode{%v}", n.val)
}

// booleanNode ...
type booleanNode struct {
	parent node
	val    bool
}

var _ node = (*booleanNode)(nil)

func (n *booleanNode) handle(e event) (node, error) {
	switch {
	case e.val == 't':
		n.val = true
		return n, nil
	case e.val == 'f':
		n.val = false
		return n, nil
	case e.val == '\n', e.val == '\r':
		return n.parent, nil
	default:
		return n, nil
	}
}

func (n *booleanNode) render(w io.Writer) error {
	return nil
}

func (n *booleanNode) String() string {
	return fmt.Sprintf("booleanNode{%v}", n.val)
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
