package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const debug = true

const (
	leftMeta  = "{{"
	rightMeta = "}}"
	digits    = "0123456789"
	hexDigits = "0123456789abcdefABCDEF"

	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabet      = alphabetLower + alphabetUpper

	underscore    = "_"
	alphaNumerics = alphabet + digits + underscore
)

const eof rune = -1

type itemType int

const (
	itemError itemType = iota
	itemDot
	itemEOF
	itemElse
	itemEnd
	itemField
	itemIdentifier
	itemIf
	itemLeftMeta
	itemNumber
	itemPipe
	itemRange
	itemRawString
	itemRightMeta
	itemString
	itemText
)

func (i itemType) String() string {
	switch i {
	case itemError:
		return "itemError"
	case itemDot:
		return "itemDot"
	case itemEOF:
		return "itemEOF"
	case itemElse:
		return "itemElse"
	case itemEnd:
		return "itemEnd"
	case itemField:
		return "itemField"
	case itemIdentifier:
		return "itemIdentifier"
	case itemIf:
		return "itemIf"
	case itemLeftMeta:
		return "itemLeftMeta"
	case itemNumber:
		return "itemNumber"
	case itemPipe:
		return "itemPipe"
	case itemRange:
		return "itemRange"
	case itemRawString:
		return "itemRawString"
	case itemRightMeta:
		return "itemRightMeta"
	case itemString:
		return "itemString"
	case itemText:
		return "itemText"
	default:
		return "(unknown)"
	}
}

type item struct {
	typ itemType
	val string
}

type lexer struct {
	name  string
	input string
	start int
	pos   int
	width int
	items chan item
}

type stateFn func(*lexer) stateFn

// String Overloads Printf target method for debug-printing items.
func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}

	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}

	return fmt.Sprintf("%q", i.val)
}

func (i item) equals(o item) bool {
	return i.typ == o.typ && i.val == o.val
}

func lexIdentifier(l *lexer) stateFn {
	l.log("lexIdentifier(*lexer)\n")

	if !l.accept(alphabet) {
		return l.errorf("invalid identifier start character: %v", l.peek())
	}

	l.acceptRun(alphaNumerics)
	l.emit(itemIdentifier)
	return lexInsideAction
}

func lexNumber(l *lexer) stateFn {
	l.log("lexNumber(*lexer)\n")

	l.accept("+-")
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	l.accept("i")
	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(itemNumber)
	return lexInsideAction
}

func lexRawQuote(l *lexer) stateFn {
	l.log("lexRawQuote(*lexer)\n")

	// Scan until we find end-raw-quote.
	// If we encounter eof or unescaped \n, error
	// Return item{itemRawString, l.input[l.start:l.pos]}

	for {
		if l.peek() == '`' {
			l.emit(itemRawString)
			l.accept("`")
			l.ignore()
			return lexInsideAction
		}

		switch r := l.next(); {
		case r == eof:
			return l.errorf("unclosed raw quote")
		case r == '\\':
			l.next()
		default:
			continue
		}
	}
}

func lexQuote(l *lexer) stateFn {
	l.log("lexQuote(*lexer)\n")

	// Scan until we find end-quote.
	// If we encounter eof or unescaped \n, error
	// Return item{itemString, l.input[l.start:l.pos]}

	for {
		if l.peek() == '"' {
			l.emit(itemString)
			l.accept("\"")
			l.ignore()
			return lexInsideAction
		}

		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed quote")
		case r == '\\':
			l.next()
		default:
			continue
		}
	}
}

func lexRightMeta(l *lexer) stateFn {
	l.log("lexRightMeta(*lexer)\n")
	l.pos += len(rightMeta)
	l.emit(itemRightMeta)
	return lexText
}

func lexInsideAction(l *lexer) stateFn {
	l.log("lexInsideAction(*lexer)\n")
	for {
		if strings.HasPrefix(l.input[l.pos:], rightMeta) {
			return lexRightMeta
		}

		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed action")
		case isSpace(r):
			l.ignore()
		case r == '|':
			l.emit(itemPipe)
		case r == '"':
			return lexQuote
		case r == '`':
			return lexRawQuote
		case r == '+' || r == '-' || '0' <= r && r <= '9':
			l.backup()
			return lexNumber
		case isAlphaNumeric(r):
			l.backup()
			return lexIdentifier
		}
	}
}

func lexLeftMeta(l *lexer) stateFn {
	l.log("lexLeftMeta(*lexer)\n")
	l.pos += len(leftMeta)
	l.emit(itemLeftMeta)
	return lexInsideAction
}

func lexText(l *lexer) stateFn {
	l.log("lexText(*lexer)\n")
	for {
		if strings.HasPrefix(l.input[l.pos:], leftMeta) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeftMeta
		}

		if l.next() == eof {
			break
		}
	}

	if l.pos > l.start {
		l.emit(itemText)
	}

	l.emit(itemEOF)
	return nil
}

// New lexing function for evaluating states.
// func (l *lexer) nextItem() item {
// 	for {
// 		select {
// 		case item := <-l.items:
// 			return item
// 		default:
// 			l.state = l.state(l)
// 		}
// 	}
// 	panic("not reached")
// }

func (l *lexer) log(format string, args ...interface{}) {
	if debug {
		fmt.Printf(format, args...)
	}
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.log("lexer.errorf(%q, ...)\n", format)
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func (l *lexer) emit(t itemType) {
	l.log("lexer.emit(%v)\n", t)
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	l.log("lexer.accept(%q)\n", valid)
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	l.log("lexer.acceptRun(%q)\n", valid)
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) peek() rune {
	l.log("lexer.peek()\n")
	rune := l.next()
	l.backup()
	return rune
}

func (l *lexer) ignore() {
	l.log("lexer.ignore()\n")
	l.start = l.pos
}

func (l *lexer) backup() {
	l.log("lexer.backup()\n")
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

	l.log("lexer.next() => (%v, %q)\n", r, r)
	return r
}

func (l *lexer) run() {
	l.log("lexer.run()\n")
	for state := lexText; state != nil; {
		state = state(l)
	}
	fmt.Printf(" * Lexer Finished: Closing Channel\n")
	close(l.items)
}

// Revised lex function which hides the token channel,
// and buffers the token channel for better parallelization.
// func lex(name, input string) *lexer {
//   l := &lexer{
//     name: name,
//     input: input,
//     state: lexText,
//     items: make(chan item, 2), // Make a chan that can hold 2 items
//   }
//   return l
// }

func lex(name, input string) (*lexer, chan item) {
	fmt.Println("Starting Lexer...", name, input)
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item, 10),
	}
	go l.run()
	return l, l.items
}

func isAlphaNumeric(r rune) bool {
	return strings.ContainsRune(alphaNumerics, r)
}

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\n':
		return true
	default:
		return false
	}
}
