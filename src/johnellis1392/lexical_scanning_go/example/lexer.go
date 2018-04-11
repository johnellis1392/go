package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const debug = false
const eof rune = -1

const (
	digits    = "0123456789"
	hexDigits = "0123456789abcdefABCDEF"

	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetFull  = alphabetLower + alphabetUpper

	underscore    = "_"
	alphaNumerics = underscore + digits + alphabetFull
	whitespace    = " \n\t"
)

type itemType int

const (
	itemEOF itemType = -(iota + 1)
	itemError

	itemIdent
	itemNumber
	itemString
	itemMajorVersion
	itemPath

	itemEquals
	itemLeftBrace
	itemRightBrace
	itemSemicolon
	itemSlash
	itemDot
)

func (t itemType) String() string {
	switch t {
	case itemEOF:
		return "itemEOF"
	case itemError:
		return "itemError"
	case itemIdent:
		return "itemIdent"
	case itemNumber:
		return "itemNumber"
	case itemString:
		return "itemString"
	case itemMajorVersion:
		return "itemMajorVersion"
	case itemEquals:
		return "itemEquals"
	case itemLeftBrace:
		return "itemLeftBrace"
	case itemRightBrace:
		return "itemRightBrace"
	case itemSemicolon:
		return "itemSemicolon"
	default:
		return "(unknown)"
	}
}

type stateFn func(*lexer) stateFn

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	return fmt.Sprintf("item{%q, %q}", i.typ, i.val)
}

type lexer struct {
	input string
	start int
	pos   int
	width int

	items chan item
	state stateFn
}

// file := { decl }
// decl := ident '=' value
// value := ( ident | string | number | object )
// ident := [a-zA-Z][a-zA-Z0-9_]*
// string := '"' ( '\\' . | [^\n\0"] ) '"'
// number := [0-9]* (\. [0-9]*)?
// object := '{' { decl } '}'

func lexNumber(l *lexer) stateFn {
	l.log("lexNumber()")
	var r rune

loop: // Whole Numbers
	for {
		switch r = l.next(); {
		case isNumber(r):
			continue
		default:
			l.backup()
			break loop
		}
	}

	// Decimal Point
	if r = l.next(); r != '.' {
		l.backup()
		l.emit(itemNumber)
		return lexText
	}

	if r = l.next(); !isNumber(r) {
		return l.errorf("unexpected character in float: %v", r)
	}

	// Predicate
	for {
		switch r = l.next(); {
		case isNumber(r):
			continue
		default:
			l.backup()
			l.emit(itemNumber)
			return lexText
		}
	}
}

func lexString(l *lexer) stateFn {
	l.log("lexString()")
	var r rune
	if r = l.next(); r != '"' {
		return l.errorf("illegal start of string: %q", r)
	} else {
		l.ignore()
	}

	for {
		switch r = l.next(); {
		case r == '"':
			// End of String
			l.backup()
			l.emit(itemString)
			l.next()
			l.ignore()
			return lexText
		case r == '\\':
			l.next() // Skip escaped character
			continue
		case r == '\n' || r == eof:
			return l.errorf("unclosed string")
		default:
			continue
		}
	}
}

func lexIdentMajorVersion(l *lexer) stateFn {
	l.log("lexIdentMajorVersion()")
	var r rune

	// Must start with number
	if r = l.next(); !isNumber(r) {
		l.log("First Character: %q", r)
		return l.errorf("invalid character in major version number: %v", r)
	}

loop1: // First Number
	for {
		switch r = l.next(); {
		case isNumber(r):
			l.log("Found Number: %q", r)
			continue
		default:
			l.log("Found Other Thing: %q", r)
			l.backup()
			break loop1
		}
	}

loop2: // Collect Remainder of versions
	for {
		// Decimal Point
		l.log("First Character in loop2: %q", r)
		if r = l.next(); r != '.' {
			return l.errorf("invalid character in major version: %v", r)
		}

		// Must Consume at least 1 Number
		if r = l.next(); !isNumber(r) {
			return l.errorf("invalid character in major version: %v", r)
		}

		l.log("Second Character in loop2: %q", r)
	loop3: // Collect Number
		for {
			switch r = l.next(); {
			case isNumber(r):
				l.log("Found number in loop3: %q", r)
				continue loop3
			case r == '.':
				// Another Segment Available for Version: Continue
				l.log("Found '.' in loop3: %q", r)
				l.backup()
				continue loop2
			default:
				l.log("Found something else in loop3: %q", r)
				// Version Number Done
				l.backup()
				l.emit(itemIdent)
				return lexText
			}
		}
	}
}

func lexIdent(l *lexer) stateFn {
	l.log("lexIdent()")
	// Get first character (must be letter)
	if r := l.next(); !isCharacter(r) {
		return l.errorf("illegal start of identifier: %v", r)
	}

	// Get rest of letters
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			l.log("lexIdent() => %q", r)
			continue
		case r == '-':
			l.log("lexIdent() => %q", r)
			// Suffixed Major Version
			return lexIdentMajorVersion
		default:
			l.log("lexIdent() => %q", r)
			l.backup()
			l.emit(itemIdent)
			return lexText
		}
	}
}

func lexValue(l *lexer) stateFn {
	l.log("lexValue()")
	switch r := l.next(); {
	case r == '"':
		l.backup()
		return lexString
	case isNumber(r):
		l.backup()
		return lexNumber
	case isCharacter(r):
		l.backup()
		return lexIdent
	default:
		return l.errorf("unknown value type - unmatched character: %q", r)
	}
}

func lexText(l *lexer) stateFn {
	l.log("lexText()")
	for {
		switch r := l.next(); {
		case r == '"' || isAlphaNumeric(r):
			l.backup()
			return lexValue
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
			l.log("lexText() => ';'")
			l.emit(itemSemicolon)
			continue
		case r == '/':
			l.log("lexText() => '/'")
			l.emit(itemSlash)
			continue
		case r == '.':
			l.log("lexText() => '.'")
			l.emit(itemDot)
			continue
		case isWhitespace(r):
			l.log("lexText() => '%q'", r)
			l.ignore()
			continue
		case r == eof:
			l.log("lexText() => EOF")
			l.emit(itemEOF)
			return nil
		default:
			// Invalid Character
			return l.errorf("invalid character: %q, %v", r, r)
		}
	}
}

func newLexer(input string) *lexer {
	const capacity = 10
	l := &lexer{
		input: input,
		start: 0,
		pos:   0,
		width: 0,
		items: make(chan item, capacity),
		state: lexText,
	}
	return l
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	const context = 2
	if debug {
		fmt.Printf(format, args...)
		fmt.Println()
	}

	contextStr := l.input[l.pos-context : l.pos+context]
	inputStr := fmt.Sprintf(format, args...)
	err := fmt.Sprintf("ERR: near '%s' -- %s", contextStr, inputStr)

	l.items <- item{itemError, err}
	return nil
}

func (l *lexer) log(format string, args ...interface{}) {
	if debug {
		fmt.Printf(format, args...)
		fmt.Println()
	}
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
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

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) run() {
	for l.state != nil {
		l.state = l.state(l)
	}
	fmt.Println("Closing Channel...")
	close(l.items)
}

func lex(input string) *lexer {
	l := newLexer(input)
	go l.run()
	return l
}

func isWhitespace(r rune) bool {
	return strings.ContainsRune(whitespace, r)
}

func isNumber(r rune) bool {
	return strings.ContainsRune(digits, r)
}

func isCharacter(r rune) bool {
	return strings.ContainsRune(alphabetFull, r)
}

func isAlphaNumeric(r rune) bool {
	return strings.ContainsRune(alphaNumerics, r)
}
