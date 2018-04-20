package main

const (
	bufferSize = 10
)

type parser struct {
	input  string
	output chan Inst
	pos    int
}

func (p *parser) run() {

}

func newParser(input string) *parser {
	return &parser{
		input:  input,
		output: make(chan Inst, bufferSize),
		pos:    0,
	}
}

func Parse(input string) chan Inst {
	p := newParser(input)
	go p.run()
	return p.output
}
