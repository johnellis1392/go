package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

// type commandType uint32
//
// const (
// 	cmdInsert commandType = iota
// 	cmdDelete
// 	cmdUp
// 	cmdReplace
// 	cmdAppend
// 	cmdEnter
// 	cmdNext
// 	cmdPrev
// 	cmdAppendParent
// )

type eventType uint32

const (
	errorEvent eventType = iota
	keyPress
	keyRepeat
	keyRelease
)

func (et eventType) String() string {
	switch et {
	case errorEvent:
		return "errorEvent"
	case keyPress:
		return "keyPress"
	case keyRepeat:
		return "keyRepeat"
	case keyRelease:
		return "keyRelease"
	default:
		return "[unkown]"
	}
}

type event struct {
	typ eventType
	val rune
}

func (e event) String() string {
	return fmt.Sprintf("event{%v, %q}", e.typ, e.val)
}

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

// editNode is a parse tree node that allows for user editing
// type editNode struct {
// 	parent node
// 	val    []rune
// 	pos    int
// }
//
// var _ node = (*editNode)(nil)
//
// func (n *editNode) insert(r rune) {
// 	if n.pos <= 0 {
// 		n.pos = 0
// 	}
//
// 	if n.pos >= len(n.val) {
// 		n.pos = len(n.val)
// 	}
//
// 	result := make([]rune, len(n.val)+1)
// 	// fmt.Printf(" * (%v).insert(%q) => len: %v\n", n, r, len(n.val))
// 	for i, c := range n.val[:n.pos] {
// 		result[i] = c
// 	}
//
// 	result[n.pos] = r
//
// 	for i, c := range n.val[n.pos:] {
// 		result[int(n.pos)+i+1] = c
// 	}
//
// 	n.val = result
// }
//
// func (n *editNode) delete() {
// 	if len(n.val) == 0 {
// 		return
// 	}
//
// 	if n.pos <= 0 {
// 		n.pos = 0
// 	}
//
// 	if n.pos >= len(n.val) {
// 		n.pos = len(n.val)
// 	}
//
// 	n.val = append(n.val[n.pos:], n.val[n.pos+1:]...)
// }
//
// func (n *editNode) handle(e event) (node, error) {
// 	// fmt.Printf(" * (%v).handle(%v)\n", n, e)
//
// 	if n.val == nil {
// 		n.val = []rune{}
// 	}
//
// 	switch {
// 	case isAlphaNumeric(e.val), isWhitespace(e.val):
// 		// Insert Character into Buffer
// 		n.insert(e.val)
// 		return n, nil
// 	case e.val == '\n', e.val == '\r':
// 		n.pos = 0
// 		return n.parent, nil
// 	default:
// 		return n, nil
// 	}
// }
//
// func (n *editNode) render(w io.Writer) error {
// 	return nil
// }
//
// func (n *editNode) String() string {
// 	return fmt.Sprintf("editNode{pos: %v, val: \"%v\"}", n.pos, string(n.val))
// }

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

// editor represents the state of the editor object.
type editor struct {
	root   *rootNode
	node   node
	events chan event
}

func newEditor() *editor {
	const chanSize = 10
	root := &rootNode{nil}
	ed := editor{
		root:   root,
		node:   root,
		events: make(chan event, chanSize),
	}
	return &ed
}

func (ed *editor) handle(ev event) error {
	if n, err := ed.node.handle(ev); err != nil {
		return err
	} else {
		ed.node = n
		return nil
	}
}

func (ed *editor) eventLoop() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			// e.events <- event{errorEvent, fmt.Sprintf("an error occurred while reading from stdin: %v", err.Error())}
			ed.events <- event{errorEvent, 0}
			continue
		}

		ed.events <- event{keyPress, rune(input[0])}
	}
}

func (ed *editor) render(w io.Writer) {
	if err := ed.root.render(w); err != nil {
		w.Write([]byte(err.Error()))
	}
}

// textBuffer represents an editable buffer of text data.
type textBuffer struct {
	pos int
	val []rune
}

func (n *textBuffer) insert(r rune) {
	if n.val == nil {
		n.val = []rune{}
	}

	if n.pos <= 0 {
		n.pos = 0
	}

	if n.pos >= len(n.val) {
		n.pos = len(n.val)
	}

	result := make([]rune, len(n.val)+1)
	// fmt.Printf(" * (%v).insert(%q) => len: %v\n", n, r, len(n.val))
	for i, c := range n.val[:n.pos] {
		result[i] = c
	}

	result[n.pos] = r

	for i, c := range n.val[n.pos:] {
		result[int(n.pos)+i+1] = c
	}

	n.val = result
}

func (n *textBuffer) delete() {
	if n.val == nil {
		n.val = []rune{}
	}

	if len(n.val) == 0 {
		return
	}

	if n.pos <= 0 {
		n.pos = 0
	}

	if n.pos >= len(n.val) {
		n.pos = len(n.val)
	}

	n.val = append(n.val[n.pos:], n.val[n.pos+1:]...)
}

func (n *textBuffer) String() string {
	return fmt.Sprintf("textBuffer{pos: %v, val: \"%v\"}", n.pos, string(n.val))
}

// type editor2 struct {
// 	root  *rootNode
// 	node  node
// 	input chan rune
// }
//
// func (ed *editor2) run() {
// 	for key := range ed.input {
// 		ev := event{keyPress, key}
// 		ed.node.handle(ev)
// 	}
// }
//
// func newEditor2() *editor2 {
// 	ed := editor2{
// 		root:  nil,
// 		node:  nil,
// 		input: make(chan rune),
// 	}
// 	go ed.run()
// 	return &ed
// }

// Main
func main() {
	fmt.Println("Starting Editor...")

	// e := newEditor()
	// go e.eventLoop()

	// Main Render Loop
	// for {
	// 	switch ev := <-e.events; ev.typ {
	// 	case errorEvent:
	// 		fmt.Println(ev.val)
	// 	case keyPress:
	// 		e.node.handle(ev)
	// 		e.render(os.Stdout)
	// 	default:
	// 		continue
	// 	}
	// }

	e := newEditor()
	events := []event{
		{keyPress, 'i'},
		{keyPress, 'o'},
		{keyPress, 'a'},
		{keyPress, '"'},
		{keyPress, 'i'},
		{keyPress, '\n'},
		{keyPress, '['},
		{keyPress, '0'},
		{keyPress, '\n'},
	}

	fmt.Println()
	for _, ev := range events {
		if err := e.handle(ev); err != nil {
			panic(err)
		}
		e.render(os.Stdout)
		fmt.Println()
	}
	fmt.Println()
}
