package main

import "fmt"

type nodeType int

const (
	nodeErr nodeType = iota
	nodeEOF

	nodeFile
	nodeObj
	nodeDecl

	nodeValue
	nodeNumber
	nodeString
	nodeIdent
	nodePath

	nodeEquals
	nodeSemicolon
	nodeLeftBrace
	nodeRightBrace
)

func (t nodeType) String() string {
	switch t {
	case nodeErr:
		return "nodeErr"
	case nodeObj:
		return "nodeObj"
	case nodeDecl:
		return "nodeDecl"
	case nodeNumber:
		return "nodeNumber"
	case nodeString:
		return "nodeString"
	case nodeIdent:
		return "nodeIdent"
	case nodePath:
		return "nodePath"
	case nodeEquals:
		return "nodeEquals"
	case nodeSemicolon:
		return "nodeSemicolon"
	case nodeLeftBrace:
		return "nodeLeftBrace"
	case nodeRightBrace:
		return "nodeRightBrace"
	default:
		return "(unknown)"
	}
}

func (t nodeType) isValue() bool {
	switch t {
	case nodeValue, nodeIdent, nodeString, nodePath, nodeNumber, nodeObj:
		return true
	default:
		return false
	}
}

type parseFn func(*parser) parseFn

type node interface {
	Type() nodeType
	String() string
	semant() error // Defined in `./marshaller.go`
}

type nerror struct {
	val string
}

func (n nerror) Error() string {
	return n.val
}

func (n nerror) Type() nodeType {
	return nodeErr
}

func (n nerror) String() string {
	return fmt.Sprintf("nerror{\"%s\"}", n.val)
}

type ndecl struct {
	ident node
	val   node
}

func (n ndecl) Type() nodeType {
	return nodeDecl
}

func (n ndecl) String() string {
	return fmt.Sprintf("ndecl{ident: %s, val: %s}", n.ident.String(), n.val.String())
}

type nobject struct {
	decls []node
}

func (n nobject) Type() nodeType {
	return nodeObj
}

func (n nobject) String() string {
	return fmt.Sprintf("nobject{decls: %q}", n.decls)
}

type nterm struct {
	typ nodeType
	val string
}

func (n nterm) Type() nodeType {
	return n.typ
}

func (n nterm) String() string {
	return fmt.Sprintf("nterm{%q, %s}", n.typ, n.val)
}

type nfile struct {
	decls []node
}

func (n nfile) Type() nodeType {
	return nodeFile
}

func (n nfile) String() string {
	var s string
	s += "nfile{"
	for _, d := range n.decls {
		s += "  " + d.String() + "\n"
	}
	s += "}"
	return s
}

// State Functions

// Reduce File
func parseAfterFile(p *parser) parseFn {
	if n, ok := p.pop(); !ok || n.Type() != nodeEOF {
		return p.errorf("illegal token at end of file: %q", n)
	}

	var decls []node
	for {
		switch n, ok := p.pop(); {
		case !ok: // No More; Done
			p.push(nfile{decls})
			return nil
		case n.Type() == nodeDecl:
			decls = append(decls, n)
			continue
		default:
			return p.errorf("illegal node found in file parse: %q", n)
		}
	}
}

// Reduce Object
func parseAfterObject(p *parser) parseFn {
	if n, ok := p.pop(); !ok || n.Type() != nodeRightBrace {
		return p.errorf("parse error, failed object: nodeRightBrace")
	}

	var decls []node
	for {
		switch n, ok := p.pop(); {
		case !ok:
			return p.errorf("failed parse of object: missing nodeLeftBrace")
		case n.Type() == nodeDecl:
			decls = append(decls, n)
			continue
		case n.Type() == nodeLeftBrace:
			p.push(nobject{decls})
			// return p.popState()
			return parseAfterValue
		default:
			return p.errorf("illegal node found in object: %q", n)
		}
	}
}

func parseInsideObject(p *parser) parseFn {
	switch n := p.shift(); n.Type() {
	case nodeRightBrace:
		return parseAfterObject
	case nodeIdent:
		return parseDecl
	default:
		return p.errorf("illegal token in object: %q", n)
	}
}

// Reduce Declaration
func parseAfterDecl(p *parser) parseFn {
	if !p.match(nodeIdent, nodeEquals, nodeValue, nodeSemicolon) {
		return p.errorf("invalid parse: declaration failed")
	}

	var ok bool
	var id, v node

	// Semicolon
	if _, ok = p.pop(); !ok {
		return p.errorf("failed parse declaration: unexpected end of tokent stack")
	}

	// Value
	if v, ok = p.pop(); !ok {
		return p.errorf("failed parse declaration: unexpected end of tokent stack")
	}

	// Equals
	if _, ok = p.pop(); !ok {
		return p.errorf("failed parse declaration: unexpected end of tokent stack")
	}

	// Ident
	if id, ok = p.pop(); !ok {
		return p.errorf("failed parse declaration: unexpected end of tokent stack")
	}

	p.push(ndecl{id, v})
	// return p.popState()
	// return parseInsideObject

	// Continue Parse:
	// * If next token is '}', we're in an object; close
	// * If next token is ident, we're uh...?; parse declaration
	// * If next token is eof, assume end; return
	// * Otherwise, error
	switch n := p.shift(); n.Type() {
	case nodeRightBrace:
		return parseAfterObject
	case nodeIdent:
		return parseDecl
	case nodeEOF:
		return parseAfterFile
	default:
		return p.errorf("illegal token following declaration: %q", n)
	}
}

func parseAfterValue(p *parser) parseFn {
	switch n := p.shift(); n.Type() {
	case nodeSemicolon:
		return parseAfterDecl
	default:
		return p.errorf("illegal token after value: %q", n)
	}
}

func parseValue(p *parser) parseFn {
	switch n := p.shift(); n.Type() {
	case nodeLeftBrace:
		// p.pushState(parseAfterValue)
		return parseInsideObject
	case nodeIdent, nodeString, nodeNumber, nodePath:
		return parseAfterValue
	default:
		return p.errorf("expected value, found: %q", n)
	}
}

func parseDecl(p *parser) parseFn {
	switch n := p.shift(); n.Type() {
	case nodeEquals:
		return parseValue
	default:
		return p.errorf("illegal token in declaration: %q", n)
	}
}

func parseFile(p *parser) parseFn {
	switch n := p.shift(); n.Type() {
	case nodeIdent:
		return parseDecl
	case nodeEOF:
		return nil
	default:
		return p.errorf("illegal start of declaration: %q", n)
	}
}

// Parser Functions
type parser struct {
	input  chan token
	output chan node
	state  parseFn
	stack  []node
	// states []parseFn
}

func (p *parser) errorf(format string, args ...interface{}) parseFn {
	err := fmt.Sprintf(format, args...)
	p.output <- nerror{err}
	return nil
}

func (p *parser) match(args ...nodeType) bool {
	n := len(args)
	m := len(p.stack)
	if n > m {
		return false
	}

	for i := range args {
		switch t1, t2 := args[i], p.stack[m-n+i].Type(); {
		case t1 == nodeValue && t2.isValue():
			continue
		case t1 == t2:
			continue
		default:
			return false
		}
	}

	return true
}

func (p *parser) shift() node {
	var n node

	switch t := <-p.input; t.typ {
	case tokenErr:
		n = nerror{t.val}
	case tokenEOF:
		n = nterm{nodeEOF, t.val}
	case tokenString:
		n = nterm{nodeString, t.val}
	case tokenNumber:
		n = nterm{nodeNumber, t.val}
	case tokenIdent:
		n = nterm{nodeIdent, t.val}
	case tokenPath:
		n = nterm{nodePath, t.val}
	case tokenEquals:
		n = nterm{nodeEquals, t.val}
	case tokenSemicolon:
		n = nterm{nodeSemicolon, t.val}
	case tokenLeftBrace:
		n = nterm{nodeLeftBrace, t.val}
	case tokenRightBrace:
		n = nterm{nodeRightBrace, t.val}
	default:
		n = nerror{fmt.Sprintf("illegal token: %q", t)}
	}

	// fmt.Printf(" * Received Token: %q\n", n)
	p.push(n)
	return n
}

func (p *parser) push(n node) {
	p.stack = append(p.stack, n)
}

func (p *parser) pop() (node, bool) {
	if len(p.stack) == 0 {
		return nil, false
	}
	n := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return n, true
}

// func (p *parser) pushState(f parseFn) {
// 	p.states = append(p.states, f)
// }
//
// func (p *parser) popState() parseFn {
// 	if len(p.states) == 0 {
// 		return nil
// 	}
// 	f := p.states[len(p.states)-1]
// 	p.states = p.states[:len(p.states)-1]
// 	return f
// }

func (p *parser) run() {
	for p.state != nil {
		p.state = p.state(p)
	}

	if len(p.stack) != 1 {
		p.output <- nerror{"invalid parse"}
	} else {
		p.output <- p.stack[0]
	}

	close(p.output)
}

func newParser(input chan token) *parser {
	p := parser{
		input:  input,
		output: make(chan node),
		state:  parseFile,
		// states: []parseFn{},
	}
	return &p
}

func parse(input chan token) chan node {
	p := newParser(input)
	go p.run()
	return p.output
}
