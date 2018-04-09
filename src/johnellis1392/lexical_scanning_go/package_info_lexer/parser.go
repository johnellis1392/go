package main

import "fmt"

type nodeType int

const (
	nodeError nodeType = iota

	nodeIdent
	nodeNumber
	nodeString
	nodePath

	nodeObject
	nodeDecl

	nodeEquals
	nodeSemiColon
	nodeLeftBrace
	nodeRightBrace
)

type node interface {
	Type() nodeType
}

type ndecl struct {
	ident *nident
	val   node
}

func (n *ndecl) Type() nodeType {
	return nodeDecl
}

type nnumber struct {
	val string
}

func (n *nnumber) Type() nodeType {
	return nodeNumber
}

type nident struct {
	val string
}

func (n *nident) Type() nodeType {
	return nodeIdent
}

type nerror struct {
	val string
}

func (n *nerror) Type() nodeType {
	return nodeError
}

type nstring struct {
	val string
}

func (n *nstring) Type() nodeType {
	return nodeString
}

type npath struct {
	val string
}

func (n *npath) Type() nodeType {
	return nodePath
}

type nobject struct {
	fields []*node
}

func (n *nobject) Type() nodeType {
	return nodeObject
}

type parser struct {
	input  chan item
	output chan node
	stack  []node
	state  parseFn
	node   node
}

type parseFn func(*parser) parseFn

func parseObject(p *parser) parseFn {
	return nil
}

func parseDecl(p *parser) parseFn {
	switch n := p.shift(); n.Type() {
	case nodeRightBrace:
		return parseObject
	case nodeNumber, nodeIdent, nodePath, nodeString:
		// Convert to node and return
		_ = p.pop()
		ident := p.pop().(*nident)
		decl := &ndecl{ident, n}
		p.push(decl)
		return parseFile
	default:
		return p.errorf("invalid parse: %v", n)
	}
}

func parseIdent(p *parser) parseFn {
	switch n := p.shift(); {
	case n.Type() == nodeEquals:
		return parseDecl
	default:
		// Invalid Parse
		return p.errorf("invalid token: %v", n)
	}
}

func parseFile(p *parser) parseFn {
	switch n := p.shift(); {
	case n.Type() == nodeError:
		return p.errorf("received lexer error: %v", n.Type())
	case n.Type() == nodeIdent:
		// return parseDecl
		return parseIdent
	default:
		// An error occurred
		return p.errorf("received illegal token: %v", n)
	}
}

// Parser Functions
func (p *parser) errorf(format string, args ...interface{}) parseFn {
	err := fmt.Sprintf(format, args...)
	p.output <- &nerror{err}
	return nil
}

// func (p *parser) reduce(args ...nodeType)

func (p *parser) shift() node {
	i, open := <-p.input
	if !open {
		p.node = nil
		return nil
	}
	n := newNode(&i)
	p.node = n
	return n
}

func newNode(i *item) node {
	var n node
	switch i.typ {
	case itemIdentifier:
		n = &nident{i.val}
	case itemNumber:
		n = &nnumber{i.val}
	case itemPath:
		n = &npath{i.val}
	case itemString:
		n = &nstring{i.val}
	default:
		n = nil
	}
	return n
}

func (p *parser) push(n node) {
	p.stack = append(p.stack, n)
}

func (p *parser) empty() bool {
	return len(p.stack) == 0
}

func (p *parser) pop() node {
	if p.empty() {
		return nil
	}
	n := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return n
}

func (p *parser) value() (*node, bool) {
	if len(p.stack) != 1 {
		// Invalid number of nodes
		return nil, false
	}
	return &p.stack[0], true
}

func (p *parser) run() {
	for p.state != nil {
		p.state = p.state(p)
	}

	if v, ok := p.value(); ok {
		p.output <- *v
	} else {
		p.errorf("invalid parse tree state: closing")
	}

	close(p.output)
}

func newParser(input chan item) *parser {
	p := parser{
		input: input,
		state: parseFile,
		node:  nil,
	}
	return &p
}

func parse(input chan item) chan node {
	p := newParser(input)
	go p.run()
	return p.output
}
