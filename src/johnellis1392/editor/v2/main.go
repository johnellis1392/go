package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
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

const (
	fps = 2
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

// Staging Functions for Testing
func example1() {
	e := newEditor()
	go e.eventLoop()

	// Main Render Loop
	for {
		switch ev := <-e.events; ev.typ {
		case errorEvent:
			fmt.Println(ev.val)
		case keyPress:
			e.node.handle(ev)
			e.render(os.Stdout)
		default:
			continue
		}
	}
}

func example2() {
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

// Example Based on Code Here:
// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-2-drawing-the-game-board
func example3() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	// vao := makeVao(triangle)
	cells := makeCells()

	for !window.ShouldClose() {
		t := time.Now()

		for x := range cells {
			for _, c := range cells[x] {
				c.checkState(cells)
			}
		}

		draw(cells, window, program)

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

// Main
func main() {
	fmt.Println("Starting Editor...")

	// example1()
	// example2()
	example3()
}
