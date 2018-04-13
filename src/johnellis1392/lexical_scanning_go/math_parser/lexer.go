package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	digits        = "0123456789"
	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetFull  = alphabetLower + alphabetUpper
	whitespace    = " \n\r\t"
	operators     = "+*"
)

const eof rune = -1

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF

	tokenNumber
	tokenPlus
	tokenMul
)

func (t tokenType) String() string {
	switch t {
	case tokenError:
		return "tokenError"
	case tokenEOF:
		return "tokenEOF"
	case tokenNumber:
		return "tokenNumber"
	case tokenPlus:
		return "tokenPlus"
	case tokenMul:
		return "tokenMul"
	default:
		return "(unknown)"
	}
}

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	return fmt.Sprintf("token{%q, \"%v\"}", t.typ, t.val)
}

type stateFn func(*lexer) stateFn

// Grammar:
// file := expr eof
// expr := expr * number
// expr := expr + number
// expr := number
// number := digit+

// State Machine Functions
func lexNumber(l *lexer) stateFn {
	fmt.Printf("lexNumber(): pos=%v, start=%v, r=%q\n", l.pos, l.start, l.peek())
	if !l.acceptRun(digits) {
		return l.errorf("illegal start of number: %q", l.peek())
	}
	l.emit(tokenNumber)
	return lexExpr
}

func lexExpr(l *lexer) stateFn {
	fmt.Printf("lexExpr(): pos=%v, start=%v, r=%q\n", l.pos, l.start, l.peek())
	for {
		switch {
		case l.accept(digits):
			l.backup()
			return lexNumber
		case l.acceptRun(whitespace):
			l.ignore()
			continue
		case l.accept("+"):
			l.emit(tokenPlus)
			continue
		case l.accept("*"):
			l.emit(tokenMul)
			continue
		case l.peek() == eof:
			l.backup()
			return lexFile
		default:
			return l.errorf("illegal char in expr: %q", l.peek())
		}
	}
}

func lexFile(l *lexer) stateFn {
	fmt.Printf("lexFile(): pos=%v, start=%v, r=%q\n", l.pos, l.start, l.peek())
	for {
		switch {
		case l.acceptRun(whitespace):
			l.ignore()
		case l.accept(digits):
			l.backup()
			return lexExpr
		case l.peek() == eof:
			// Done
			return nil
		default:
			return l.errorf("illegal character in file: %q", l.peek())
		}
	}
}

// Lexer Functions
type lexer struct {
	input  []rune
	output chan token
	state  stateFn

	pos   int
	start int
	r     string
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	err := fmt.Sprintf(format, args...)
	l.output <- token{tokenError, err}
	return nil
}

func (l *lexer) peek() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	return l.input[l.pos]
}

func (l *lexer) emit(t tokenType) {
	l.output <- token{t, string(l.input[l.start:l.pos])}
	l.start = l.pos
}

func (l *lexer) acceptRun(valid string) bool {
	if !l.accept(valid) {
		return false
	}

	for l.accept(valid) {
	}

	return true
}

func (l *lexer) accept(valid string) bool {
	if r := l.next(); contains(valid, r) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos--
	if l.pos < l.start {
		l.pos = l.start
	}
}

func (l *lexer) next() rune {
	l.pos++
	if l.pos > len(l.input) {
		return eof
	}
	r := l.input[l.pos-1]
	l.r = fmt.Sprintf("%q", r)
	return r
}

func (l *lexer) run() {
	for l.state != nil {
		l.state = l.state(l)
	}
	close(l.output)
}

func runes(input string) []rune {
	var rs []rune
	for pos := 0; pos < len(input); {
		r, width := utf8.DecodeRuneInString(input[pos:])
		rs = append(rs, r)
		pos += width
	}
	return rs
}

func newLexer(input string) *lexer {
	l := lexer{
		input:  runes(input),
		output: make(chan token),
		state:  lexFile,
		pos:    0,
		start:  0,
		r:      "",
	}
	return &l
}

func lex(input string) chan token {
	l := newLexer(input)
	go l.run()
	return l.output
}

func contains(valid string, r rune) bool {
	return strings.ContainsRune(valid, r)
}

func isDigit(r rune) bool {
	return contains(digits, r)
}

func isWhitespace(r rune) bool {
	return contains(whitespace, r)
}

func isOperator(r rune) bool {
	return contains(operators, r)
}
