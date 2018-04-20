package main

type System struct {
	CPU
	Mem
	Program
}

func NewSystem(p Program) System {
	return System{
		Program: p,
		CPU:     CPU{},
		Mem:     Mem{},
	}
}

type CPU struct {
	Eax, Ebx, Ecx, Edx int16
	Ptr                int
}

type Mem []int16

type Program struct {
	Insts []Inst
	Data  []int16
	Bss   []int16
	Raw   []int16
}
