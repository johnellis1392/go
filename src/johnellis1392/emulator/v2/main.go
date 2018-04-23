package main

import (
	"fmt"
	"io/ioutil"
)

// Addressing Modes:
// 1.   Immediate              =>   Operand IS Numeric Value
// 2.   Absolute               =>   Operand is Address of Value in Memory
// 3.   Zero-Page Absolute     =>   Operand is Address of Value in Zero-Page Memory
// 4.   Implied                =>   Operand Addresses Implied by Instruction
// 5.   Accumulator            =>   Instruction Operates on Data in Accumulator
// 6.   Indexed                =>   Operand is Added to Value of X / Y Index to Get Address of Value
// 7.   Zero-Page Indexed      =>   X / Y Index is Used to Get Value from Zero-Page Memory
// 8.   Indirect               =>   (Applies only to JMP) Operand is Address of Bytes whose Value is the Location to Jump to
// 9.   Pre-Indexed Indirect   =>   Operand is Zero-Page Address that's Added to X-Index to get Location of Address of Value
// 10.  Post-Indexed Indirect  =>   Next Two Bytes Added to Y-Index to Get Location of Address of Value
// 11.  Relative               =>   Signed Operand is Added to Program Counter to Get New Execution Location

// Instruction Set:
// ADC => Add Mem to Acc with Carry
// AND => And Mem w/ Acc
// ASL => Shift Left One Bit
// BCC => Branch on Carry Clear
// BCS => Branch on Carry Set
// BEQ => Branch on Result Zero
// BIT => Test Bits in Mem with Acc
// BMI => Branch on Result Minus
// BNE => Branch Not Zero
// BPL => Branch Result Plus
// BRK => Force Break
// BVC => Branch on Overflow Clear
// BVS => Branch on Overflow Set
// CLC => Clera Carry
// CLD => Clear Decimal
// CLI => Clear Interrupt
// CLV => Clear Overflow
// CMP => Compare Mem & Acc
// CPX => Cmp Mem & XR
// CPY => Cmp Mem & YR
// DEC => Decrement Mem by One
// DEX => Decrement X-Index by One
// DEY => Decrement Y-Index by One
// EOR => Exclusive Or
// INC => Increment Mem by One
// INX => Increment X-Index by One
// INY => Increment Y-Index by One
// JMP => Jump
// JSR => Jump & Save Return Address
// LDA => Load Acc with Mem
// LDX => Load XR with Mem
// LDY => Loda YR with Mem
// LSR => Shift Right One Bit
// NOP => No-op
// ORA => Or Mem & Acc
// PHA => Push Acc on Stack
// PHP => Push Processor Status on Stack
// PLA => Pull Acc from Stack
// PLP => Pull Processor Status from Stack
// ROL => Rotate One Bit Left (Mem | Acc)
// ROR => Rotate One Bit Right (Mem | Acc)
// RTI => Return from Interrupt
// RTS => Return from Subroutine
// SBC => Subtract Mem from Acc with Borrow
// SEC => Set Carry Flag
// SED => Set Decimal Flag
// SEI => Set Interrupt Disable Status
// STA => Store Acc in Mem
// STX => Store XR in Mem
// STY => Store YR in Mem
// TAX => Transfer Acc to XR
// TAY => Transfer Acc to YR
// TSX => Transfer Stack Pointer to XR
// TXA => Transfer XR to Acc
// TXS => Transfer XR to SP
// TYA => Transfer YR to Acc

const screenRes = 640 * 480
const memSize = 1024

type status struct {
	s bool // Sign Flag: Set if Result of Arithmetic Operation is Negative
	v bool // Overflow Flag: Set if Result of Arithmetic Operation Exceeds Register Size
	// _ bool // Bit 5: Unused
	b bool // Break Flag: Interrupt Signal Executed
	d bool // Decimal Mode Status Flag
	i bool // Interrupt Enable / Disable
	z bool // Zero Flag
	c bool // Carry Flag
}

func (s status) carry() byte {
	if s.c {
		return 1
	}
	return 0
}

func (s *status) setCarry(i byte) {
	if i != 0 {
		s.c = true
		return
	}
	s.c = false
}

type cpu struct {
	ac  byte   // Accumulator
	xr  byte   // X Index Register
	yr  byte   // Y Index Register
	sr  status // Status Register
	pc  byte   // Program Counter
	sp  byte   // Stack Pointer
	clk uint64 // Clock Timer

	inst byte // Current Instruction
}

func (c *cpu) tick() {
	c.clk++
}

type page [memSize]byte

type mem [memSize]page

type system struct {
	cpu
	mem
}

func (s *system) stopped() bool {
	return false
}

func (s *system) next() byte {
	pc := int(s.cpu.pc)
	i := s.mem[pc]
	s.cpu.pc++
	return i
}

func (s *system) peek() byte {
	return s.cpu.inst
}

func newSystem(data []byte) *system {
	return &system{
		cpu: cpu{},
		mem: mem{},
		rom: rom(data),
	}
}

func main() {
	const filename = "./test.rom"
	fmt.Printf("Reading File: %s\n", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println("Running File.")
	for s := newSystem(data); !s.stopped(); {
		s.exec()
	}
}
