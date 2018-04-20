package main

import "fmt"

type InterruptCode uint8

const (
	PrintCode InterruptCode = iota
	DrawPixelCode
)

type Interrupt interface {
	Inst
	Code() InterruptCode
}

var _ Inst = (*PrintInt)(nil)
var _ Interrupt = (*PrintInt)(nil)

type PrintInt struct{}

func (i PrintInt) Exec(s System) {
	// TODO
}

func (i PrintInt) String() string {
	return fmt.Sprintf("PrintInt{%s}", i.Code())
}

func (i PrintInt) Code() InterruptCode {
	return PrintCode
}
