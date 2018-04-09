package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const debug = true

const (
	digits    = "0123456789"
	hexDigits = "0123456789abcdefABCDEF"

	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetFull  = alphabetLower + alphabetUpper

	underscore    = "_"
	alphaNumerics = alphabetFull + digits + underscore
	whitespace    = " \n\t\r"
	terminals     = "{};="
)

// End-Of-File
const eof rune = -1

type itemType int

const (
	itemError itemType = iota
	itemEOF

	itemIdentifier
	itemNumber
	itemMajorVersion
	itemPath

	itemString
	itemLeftBrace
	itemRightBrace
	itemEquals
	itemSemiColon
)

func (t itemType) String() string {
	switch t {
	case itemError:
		return "itemError"
	case itemEOF:
		return "itemEOF"
	case itemIdentifier:
		return "itemIdentifier"
	case itemNumber:
		return "itemNumber"
	case itemMajorVersion:
		return "itemMajorVersion"
	case itemPath:
		return "itemPath"
	case itemString:
		return "itemString"
	case itemLeftBrace:
		return "itemLeftBrace"
	case itemRightBrace:
		return "itemRightBrace"
	case itemEquals:
		return "itemEquals"
	case itemSemiColon:
		return "itemSemiColon"
	default:
		return "(unknown)"
	}
}

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	return fmt.Sprintf("item{%v, \"%v\"}", i.typ, i.val)
}

type lexer struct {
	input string
	start int
	pos   int
	width int
	items chan item
}

type stateFn func(*lexer) stateFn

// file := { decl }
// decl := ident '=' value ';'
// ident := [a-zA-Z][a-zA-Z0-9_]*
// value := number | string | path | object
// number := [a bunch of stuff]
// string := '"' ( '\\' . | [^\n'"'] ) '"'
// path := [ '/' ] filename { '/' filename }
// filename := ( '.' | '..' | ident )
// object := '{' { field } '}'
// field :=

func lexPath(l *lexer) stateFn {
	// Check Starts with ('.' | '/')
	switch {
	case l.accept("."):
		break
	case l.accept("/"):
		if !l.accept(alphaNumerics) {
			// No Identifier Follows; Continue Parsing File
			l.emit(itemPath)
			return lexFile
		}

		// Identifier Follows; Path Continues
		for l.accept(alphaNumerics) {
		}

		break
	case l.accept(alphaNumerics):
		break
	default:
		// Unmatched Character; Quit
		return l.errorf("illegal character in path: %v", l.peek())
	}

	// Check Remainder of Path Parts
	for {
		if !l.accept("/") {
			// Check Leading "/"
			l.emit(itemPath)
			return lexFile
		}

		// Consume Remainder of Path Chars
		for l.accept(alphaNumerics + ".-") {
		}
	}
}

func lexString(l *lexer) stateFn {
	if !l.accept("\"") {
		return l.errorf("illegal start of string: %q", l.peek())
	}

	for {
		switch r := l.next(); {
		case r == '"':
			l.emit(itemString)
			return lexFile
		case r == '\\':
			// Escape Sequence: Skip Next Char
			l.next()
			continue
		case r == eof || r == '\n':
			return l.errorf("illegal character in string: %v", r)
		default:
			continue
		}
	}
}

func lexNumber(l *lexer) stateFn {
	// Must start with [1-9]
	if !l.accept("123456789") {
		return l.errorf("illegal start of number: %v", l.peek())
	}

	// Remainder of number body
	for l.accept(digits) {
	}

	// Parse Decimal Point
	if !l.accept(".") {
		// No Decimal Point; Found Integer
		l.emit(itemNumber)
		return lexFile
	}

	// Decimal Point Consumed; Consume Remainder of Float
	if !l.accept(digits) {
		// Error: Not Valid Float Character
		return l.errorf("illegal character in float: %v", l.peek())
	}

	// Parse Remainder of Float
	for l.accept(digits) {
	}

	// Emit Final Number
	l.emit(itemNumber)
	return lexFile
}

func lexIdentMajorVersion(l *lexer) stateFn {
	if !l.accept(digits) {
		return l.errorf("illegal start of major version: %v", l.peek())
	}

	for l.accept(digits) {
	}

	// Consume Remainder of Sem-ver Tokens
	for {
		if !l.accept(".") {
			l.emit(itemIdentifier)
			return lexFile
		}

		for l.accept(alphaNumerics) {
		}
	}
}

func lexIdentifier(l *lexer) stateFn {
	if r := l.next(); !isAlphaNumeric(r) {
		return l.errorf("illegal start of identifier: %v", r)
	}

	for l.accept(alphaNumerics) {
	}

	// Identifier Could End in Major Version
	if l.accept("-") {
		return lexIdentMajorVersion
	}

	// No Major Version
	l.emit(itemIdentifier)
	return lexFile
}

func lexDeclaration(l *lexer) stateFn {
	if !isAlphaNumeric(l.peek()) {
		return l.errorf("illegal start of declaration: %v", l.peek())
	}
	return lexIdentifier
}

func lexFile(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			l.backup()
			return lexDeclaration
		case r == '"':
			l.backup()
			return lexString
		case r == '.' || r == '/':
			l.backup()
			return lexPath
		case isWhitespace(r):
			l.ignore()
			continue
		case r == '{':
			l.emit(itemLeftBrace)
			continue
		case r == '}':
			l.emit(itemRightBrace)
			continue
		case r == '=':
			l.emit(itemEquals)
			continue
		case r == ';':
			l.emit(itemSemiColon)
			continue
		case r == eof:
			return nil
		default:
			return l.errorf("unexpected character: %v, %q", l.peek(), l.peek())
		}
	}
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	err := fmt.Sprintf(format, args...)
	// fmt.Println(err)
	l.items <- item{itemError, err}
	return nil
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
	l.width = 0
}

func (l *lexer) match(prefix string) bool {
	return strings.HasPrefix(l.input[l.pos:], prefix)
}

func (l *lexer) peek() rune {
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
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

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) run() {
	for state := lexFile; state != nil; {
		state = state(l)
	}
	fmt.Println("Closing Channel...")
	close(l.items)
}

func (l *lexer) log(format string, args ...interface{}) {
	if debug {
		fmt.Printf(format, args...)
		fmt.Println()
	}
}

func (l *lexer) dump() {
	if debug {
		fmt.Println(" * Lexer:", l.String())
		fmt.Println(" * Cursor:", l.input[l.start:l.pos])
	}
}

func (l *lexer) String() string {
	return fmt.Sprintf(
		"lexer{start: %d, pos: %d, width: %d, items: _, input: _}",
		l.start,
		l.pos,
		l.width,
	)
}

func newLexer(input string) *lexer {
	const chanSize = 10
	l := lexer{
		input: input,
		start: 0,
		pos:   0,
		width: 0,
		items: make(chan item, chanSize),
	}
	return &l
}

func lex(input string) chan item {
	l := newLexer(input)
	go l.run()
	return l.items
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\n', '\t':
		return true
	default:
		return false
	}
}

func isCharacter(r rune) bool {
	return strings.ContainsRune(alphabetFull, r)
}

func isAlphaNumeric(r rune) bool {
	return strings.ContainsRune(alphaNumerics, r)
}
