package main

import (
	"fmt"
)

type Inst interface {
	Exec(System)
	String() string
}

var _ Inst = (*AddInst)(nil)

type AddInst struct {
	Target string
	Val    int16
}

func (i AddInst) Exec(s System) {
	switch i.Target {
	case "eax":
		s.CPU.Eax += i.Val
	case "ebx":
		s.CPU.Ebx += i.Val
	case "ecx":
		s.CPU.Ecx += i.Val
	case "edx":
		s.CPU.Edx += i.Val
	default:
		// Unknown Instruction
		break
	}
}

func (i AddInst) String() string {
	return fmt.Sprintf("AddInst{%s, %s}", i.Target, i.Val)
}
