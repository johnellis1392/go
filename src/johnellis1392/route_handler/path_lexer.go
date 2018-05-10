package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Path represents a compiled path object for routing.
type Path struct {
	input       string
	uri         URI
	queryParams QueryParams
}

// URI represents a list of URI components.
type URI []string

// QueryParams represents a query parameter map.
type QueryParams map[string]string

const (
	digits        = "0123456789"
	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetFull  = alphabetLower + alphabetUpper
	alphanumeric  = "_" + digits + alphabetFull
)

const eof rune = -1

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenErr

	tokenSep
	tokenPathPart
	tokenPathParam
	tokenQuerySigil
	tokenQuerySep
	tokenQName
	tokenQVal
	tokenQEquals
)

// State Functions

type stateFn func(*lexer) stateFn

func lexQVal(l *lexer) stateFn {
	if !l.accept("=") {
		return l.errorf("illegal character in query param: %q", l.peek())
	}

	l.emit(tokenQEquals)

	for l.accept(alphanumeric) {
	}

	l.emit(tokenQVal)

	if !l.accept("&") {
		return l.errorf("illegal character in query param value: %q", l.peek())
	}

	l.emit(tokenQuerySep)
	return lexQueryParam
}

func lexQName(l *lexer) stateFn {
	if !l.accept(alphabetFull) {
		return l.errorf("illegal start of query param name: %q", l.peek())
	}

	for l.accept(alphanumeric) {
	}

	l.emit(tokenQName)
	return lexQVal
}

func lexQueryParam(l *lexer) stateFn {
	switch {
	case l.accept(alphabetFull):
		l.backup()
		return lexQName
	default:
		return l.errorf("unexpected character in query param: %q", l.peek())
	}
}

func lexPathParam(l *lexer) stateFn {
	if !l.accept("{") {
		return l.errorf("illegal character: expected '{', got '%q'", l.peek())
	}

	l.ignore()

	// Must have at least one character
	if !l.accept(alphabetFull) {
		return l.errorf("illegal start of path parameter: expected alphanumeric, got: '%q'", l.peek())
	}

	for l.accept(alphanumeric) {
	}

	if l.peek() != '}' {
		return l.errorf("illegal character in path param: expected '}', got '%q'", l.peek())
	}

	l.emit(tokenPathParam)
	l.accept("}")
	return lexPath
}

func lexPathName(l *lexer) stateFn {
	if !l.accept(alphabetFull) {
		return l.errorf("unexpected character in path name: %q", l.peek())
	}

	for l.accept(alphanumeric) {
	}

	l.emit(tokenPathPart)
	return lexPath
}

func lexPathPart(l *lexer) stateFn {
	switch {
	case l.accept(alphanumeric):
		// Path Identifier
		l.backup()
		return lexPathName
	case l.accept("{"):
		l.backup()
		return lexPathParam
	case l.accept("/"):
		return l.errorf("invalid empty path part: encountered unexpected '/'")
	default:
		return l.errorf("unexpected character: %q", l.peek())
	}
}

func lexPath(l *lexer) stateFn {
	switch {
	case l.accept("/"):
		// Path Separator
		l.emit(tokenSep)
		return lexPathPart
	case l.accept("{"):
		// Start Path Param
		l.backup()
		return lexPathPart
	case l.accept(alphanumeric):
		// Identifier: Parse Path Part
		l.backup()
		return lexPathPart
	case l.accept("?"):
		l.emit(tokenQuerySigil)
		return lexQueryParam
	default:
		return l.errorf("unexpected character in path: %q", l.peek())
	}
}

// Lexer Functions

type token struct {
	typ tokenType
	val string
}

type lexer struct {
	input  []rune
	start  int
	pos    int
	state  stateFn
	output []token
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.output = append(l.output, token{tokenErr, fmt.Sprintf(format, args...)})
	return nil
}

func (l *lexer) emit(typ tokenType) {
	if l.pos > len(l.input) {
		return
	}
	val := l.input[l.start:l.pos]
	tok := token{typ, string(val)}
	l.output = append(l.output, tok)
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	r := l.next()
	if strings.ContainsRune(valid, r) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) backup() {
	l.pos = l.start
}

func (l *lexer) peek() rune {
	if l.pos < 0 || l.pos >= len(l.input) {
		return eof
	}
	return l.input[l.pos]
}

func (l *lexer) next() rune {
	if l.pos > len(l.input) {
		return eof
	}
	r := l.input[l.pos]
	l.pos++
	return r
}

func (l *lexer) run() []token {
	for l.state != nil {
		l.state = l.state(l)
	}
	return l.output
}

func newLexer(input string) *lexer {
	runes := []rune{}
	for pos := 0; pos < len(input); {
		r, width := utf8.DecodeRuneInString(input[pos:])
		runes = append(runes, r)
		pos += width
	}

	l := lexer{
		input:  runes,
		start:  0,
		pos:    0,
		state:  lexPath,
		output: []token{},
	}

	return &l
}

func parse(toks []token) (*Path, error) {
	return nil, nil
}

// ParsePath parses the given input string into a Path object.
func ParsePath(input string) (*Path, error) {
	l := newLexer(input)
	toks := l.run()

	path, err := parse(toks)
	if err != nil {
		return nil, err
	}

	return path, nil
}
