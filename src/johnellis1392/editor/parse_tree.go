package main

import (
	"io"
	"strings"
)

const (
	digits        = "0123456789"
	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetFull  = alphabetLower + alphabetUpper

	underscore    = "_"
	alphanumerics = digits + underscore + alphabetFull
	whitespace    = " \t\r\n"
)

// parseTree is a tree structure represented the parsed form of the program
type parseTree struct {
	node      // Focused Node
	root node // Tree Root
}

// node represents a parse node
type node interface {
	handle(*event) error
	render(io.Writer) error
}

// rootNode represents the root of the parse tree
type rootNode struct {
	exprs []node
}

var _ node = (*rootNode)(nil)

func (n *rootNode) handle(e *event) error {
	if e.typ != keyPress {
		return nil
	}

	switch {
	case isChar(e.val):
		decl := &declNode{n, nil}
		n.exprs = append(n.exprs, decl)
		return decl.handle(e)
	case isDigit(e.val):
		expr := &exprNode{n, nil}
		n.exprs = append(n.exprs, expr)
		return expr.handle(e)
	case isWhitespace(e.val):
		return nil
	default:
		// NOP
		return nil
	}
}

func (n *rootNode) render(w io.Writer) error {
	for _, i := range n.exprs {
		if err := i.render(w); err != nil {
			return err
		}
	}
	return nil
}

// declNode
type declNode struct {
	parent node
	expr   node
}

var _ node = (*declNode)(nil)

func (n *declNode) handle(e *event) error {
	if e.typ != keyPress {
		return nil
	}

	switch {
	case isChar(e.val):
		// Character Pressed: Default to Identifier
		return n.expr.handle(e)
	case isDigit(e.val):
		// Digit Pressed: Default to Value
		return n.expr.handle(e)
	case isWhitespace(e.val):
		return nil
	default:
		return nil
	}
}

func (n *declNode) render(w io.Writer) error {
	return nil
}

// exprNode represents a type of expression in the language
type exprNode struct {
	parent node
	val    node
}

var _ node = (*exprNode)(nil)

func (n *exprNode) handle(e *event) error {
	if e.typ != keyPress {
		return nil
	}

	switch {
	case isChar(e.val), isDigit(e.val):
		val := &valueNode{n, nil}
		n.val = val
		return val.handle(e)
	default:
		return nil
	}
}

func (n *exprNode) render(w io.Writer) error {
	return nil
}

type assignNode struct {
	parent node
	ident  identNode
	expr   node
}

var _ node = (*assignNode)(nil)

func (n *assignNode) handle(e *event) error {
	return nil
}

func (n *assignNode) render(w io.Writer) error {
	return nil
}

// valueNode
type valueNode struct {
	parent node
	val    node
}

var _ node = (*valueNode)(nil)

func (n *valueNode) handle(e *event) error {
	if e.typ != keyPress {
		return nil
	}

	switch {
	case isChar(e.val):
		n.val = &identNode{n, ""}
		return n.val.handle(e)
	case isDigit(e.val):
		n.val = &numberNode{n, ""}
		return n.val.handle(e)
	default:
		return nil
	}
}

func (n *valueNode) render(w io.Writer) error {
	return nil
}

// identNode
type identNode struct {
	parent node
	val    string
}

var _ node = (*identNode)(nil)

func (n *identNode) handle(e *event) error {
	switch {
	case isDigit(e.val), isChar(e.val):
		// TODO: Incorporate Cursor Position
		n.val += string(e.val)
		return nil
	default:
		return nil
	}
}

func (n *identNode) render(w io.Writer) error {
	return nil
}

// numberNode
type numberNode struct {
	parent node
	val    string
}

var _ node = (*numberNode)(nil)

func (n *numberNode) handle(e *event) error {
	switch {
	case isDigit(e.val):
		n.val += string(e.val)
		return nil
	default:
		return n.parent.handle(e)
	}
}

func (n *numberNode) render(w io.Writer) error {
	return nil
}

// blockNode
type blockNode struct {
	parent node
	exprs  []node
}

var _ node = (*blockNode)(nil)

func (n *blockNode) handle(e *event) error {
	if e.typ != keyPress {
		return nil
	}

	switch {
	case e.val == '}', e.val == ';', e.val == '\n':
		// End Block
		return n.parent.handle(e)
	case isChar(e.val), isDigit(e.val), isWhitespace(e.val):
		expr := &exprNode{n, nil}
		n.exprs = append(n.exprs, expr)
		return expr.handle(e)
	default:
		return nil
	}
}

func (n *blockNode) render(w io.Writer) error {
	return nil
}

// // parseFn represents a state transition function for the editor
// type parseFn func(*parseTree, *cursor, *event) parseFn
//
// // parseRoot is the initial parse state for the parse tree
// func parseRoot(p *parseTree, c *cursor, e *event) parseFn {
// 	if e.typ == tick {
// 		return parseRoot
// 	}
//
// 	if e.typ == quit {
// 		return nil
// 	}
//
// 	switch {
// 	case isChar(e.val):
// 		// Parse Expression
// 		return parseExpr
// 	case isWhitespace(e.val):
// 		// Ignore
// 		return parseRoot
// 	case isDigit(e.val):
// 		// Value
// 		return parseExpr
// 	default:
// 		// Unknown
// 		return nil
// 	}
// }
//
// func parseExpr(p *parseTree, c *cursor, e *event) parseFn {
// 	return nil
// }
//
// func parseValue(p *parseTree, c *cursor, e *event) parseFn {
// 	return nil
// }

// * Javascript Grammar:
// file := { expr }
// expr := value | decl | assign | arith
// value := fndecl | fncall | ident | number | string | bool

func contains(s string, r rune) bool {
	return strings.ContainsRune(s, r)
}

func isWhitespace(r rune) bool {
	return contains(whitespace, r)
}

func isDigit(r rune) bool {
	return contains(digits, r)
}

func isChar(r rune) bool {
	return contains(alphabetFull, r)
}
