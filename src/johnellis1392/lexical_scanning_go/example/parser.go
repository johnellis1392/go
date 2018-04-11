package main

import "fmt"

type nodeType int

const (
	nodeError nodeType = iota

	nodeString
	nodeIdent
	nodeNumber

	nodeDecl
	nodeObject
	nodePath
)

type parseFn func(*parser) parseFn

type node struct {
	typ nodeType
	val string
}

type parser struct {
	state parseFn
	nodes chan node
	items chan item

	stack []*node
	item  item
}

// file := { decl }
// decl := ident '=' value
// value := ( ident | string | number | object )
// ident := [a-zA-Z][a-zA-Z0-9_]*
// string := '"' ( '\\' . | [^\n\0"] ) '"'
// number := [0-9]* (\. [0-9]*)?
// object := '{' { decl } '}'

func parseInsideObject(p *parser) parseFn {
	return nil
}

func parseObject(p *parser) parseFn {
	return nil
}

func parseValue(p *parser) parseFn {
	return nil
}

func parseDecl(p *parser) parseFn {
	switch i, ok := p.next(); {
	case !ok:
		if i.typ == itemEOF {
			return nil
		} else {
			return p.errorf("channel closed unexpectedly")
		}
	case i.typ == itemError:
		// Convert to nodeError
		return p.errorf(i.val)
	case i.typ == itemIdent:
		p.push(&node{nodeIdent, i.val})
		return nil
	default:
		return nil
	}
}

func parseFile(p *parser) parseFn {
	return parseDecl
}

// Parser Functions

func (p *parser) errorf(format string, args ...interface{}) parseFn {
	err := fmt.Sprintf(format, args...)
	p.nodes <- node{nodeError, err}
	close(p.nodes)
	return nil
}

func (p *parser) peek() *node {
	return p.stack[len(p.stack)-1]
}

func (p *parser) next() (*item, bool) {
	select {
	case i, ok := <-p.items:
		p.item = i
		return &i, ok
	}
}

func (p *parser) pop() *node {
	l := len(p.stack)
	if l == 0 {
		return nil
	}
	n := p.stack[l-1]
	p.stack = p.stack[0 : l-1]
	return n
}

func (p *parser) push(n *node) {
	p.stack = append(p.stack, n)
}

func (p *parser) run(l *lexer) {
	p.items = l.items
	for p.state = parseFile; p.state != nil; {
		p.state = p.state(p)
	}
	close(p.nodes)
}

func newParser() *parser {
	p := &parser{
		nodes: make(chan node, 1),
		stack: []*node{},
	}
	return p
}

func parse(l *lexer) *parser {
	p := newParser()
	go p.run(l)
	return p
}
