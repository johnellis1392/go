package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	digit    = "0123456789"
	hexDigit = "0123456789abcdefABCDEF"

	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetFull  = alphabetLower + alphabetUpper

	whitespace   = " \n\r\t"
	underscore   = "_"
	alphaNumeric = underscore + digit + alphabetFull
	pathchars    = alphaNumeric + ".-"
)

func contains(valid string, r rune) bool {
	return strings.ContainsRune(valid, r)
}

func isWhitespace(r rune) bool {
	return contains(whitespace, r)
}

func isDigit(r rune) bool {
	return contains(digit, r)
}

func isAlphanumeric(r rune) bool {
	return contains(alphaNumeric, r)
}

const eof rune = -1

type tokenType int

const (
	tokenErr tokenType = iota
	tokenEOF

	tokenString
	tokenNumber
	tokenIdent
	tokenPath

	tokenEquals
	tokenSemicolon
	tokenLeftBrace
	tokenRightBrace
)

func (t tokenType) String() string {
	switch t {
	case tokenErr:
		return "tokenErr"
	case tokenEOF:
		return "tokenEOF"
	case tokenString:
		return "tokenString"
	case tokenNumber:
		return "tokenNumber"
	case tokenIdent:
		return "tokenIdent"
	case tokenPath:
		return "tokenPath"
	case tokenEquals:
		return "tokenEquals"
	case tokenSemicolon:
		return "tokenSemicolon"
	case tokenLeftBrace:
		return "tokenLeftBrace"
	case tokenRightBrace:
		return "tokenRightBrace"
	default:
		return "(unknown)"
	}
}

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	return fmt.Sprintf("token{%v, \"%v\"}", t.typ, t.val)
}

// State Functions
type stateFn func(*lexer) stateFn

// file  := { decl }
// decl := ident '=' value
// ident := [a-zA-Z][a-zA-Z0-9_]*
// value := (object | string | ident | path | number)
// object := '{' { decl } '}'
// string := '"' ('\\' . | [^"\n\eof])* '"'
// path := ([\.\w]) { '/' ([\.\w]) }

func lexNumber(l *lexer) stateFn {
	if r := l.next(); contains("123456789", r) {
		return l.errorf("illegal start to number literal: %v, %q", r, r)
	}

	l.acceptRun(digit)
	l.emit(tokenNumber)
	return lexFile
}

func lexPath(l *lexer) stateFn {
	if !l.acceptRun(pathchars) {
		return l.errorf("invalid start of path: %v, %q", l.peek(), l.peek())
	}

	for {
		if !l.accept("/") {
			l.emit(tokenPath)
			return lexFile
		}

		l.acceptRun(pathchars)
	}
}

func lexMajorVersion(l *lexer) stateFn {
	if !l.acceptRun(digit) {
		return l.errorf("illegal start of major version: %v, %q", l.peek(), l.peek())
	}

	if !l.accept(".") {
		return l.errorf("expected '.' in major version")
	}

	if !l.acceptRun(digit) {
		return l.errorf("illegal character in major version: %v, %q", l.peek(), l.peek())
	}

	l.emit(tokenIdent)
	return lexFile
}

func lexIdent(l *lexer) stateFn {
	if !l.accept(alphabetFull) {
		return l.errorf("illegal start of identifier: %v, %q", l.peek(), l.peek())
	}

	l.acceptRun(alphaNumeric)

	// Check if Path
	if l.accept("./") {
		return lexPath
	} else if l.accept("-") {
		return lexMajorVersion
	}

	l.emit(tokenIdent)
	return lexFile
}

func lexString(l *lexer) stateFn {
	if !l.accept("\"") {
		return l.errorf("invalid start of string: %v, %q", l.peek(), l.peek())
	}

	for {
		switch r := l.next(); {
		case r == '\\':
			// Skip Next Character
			l.next()
		case r == '\n' || r == eof:
			return l.errorf("unexpected character in string: %v, %q", r, r)
		case r == '"':
			l.emit(tokenString)
			return lexFile
		default:
			continue
		}
	}
}

func lexFile(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphanumeric(r):
			l.backup()
			return lexIdent
		case isDigit(r):
			l.backup()
			return lexNumber
		case r == '"':
			l.backup()
			return lexString
		case r == '.' || r == '/':
			l.backup()
			return lexPath
		case isWhitespace(r):
			l.ignore()
		case r == '{':
			l.emit(tokenLeftBrace)
		case r == '}':
			l.emit(tokenRightBrace)
		case r == ';':
			l.emit(tokenSemicolon)
		case r == '=':
			l.emit(tokenEquals)
		case r == eof:
			l.emit(tokenEOF)
			return nil
		default:
			return l.errorf("invalid character: %v, %q", r)
		}
	}
}

// Lexer Functions
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	err := fmt.Sprintf(format, args...)
	l.output <- token{tokenErr, err}
	return nil
}

func (l *lexer) emit(t tokenType) {
	l.output <- token{t, l.input[l.start:l.pos]}
	l.start = l.pos
	l.width = 0
}

func (l *lexer) accept(valid string) bool {
	if r := l.next(); contains(valid, r) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) bool {
	// Must accept at least one
	if !l.accept(valid) {
		return false
	}

	// Accept all others
	for r := l.next(); contains(valid, r); {
		r = l.next()
	}

	l.pos -= l.width
	return true
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

type lexer struct {
	input  string
	output chan token
	state  stateFn
	start  int
	width  int
	pos    int
}

func (l *lexer) run() {
	for l.state != nil {
		l.state = l.state(l)
	}
	close(l.output)
}

func newLexer(input string) *lexer {
	l := lexer{
		input:  input,
		output: make(chan token),
		state:  lexFile,
		start:  0,
		width:  0,
		pos:    0,
	}
	return &l
}

func lex(input string) chan token {
	l := newLexer(input)
	go l.run()
	return l.output
}
