package main

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

// 0b 0110 0001 => 61 - ADC - (Indirect,X)
// 0b 0110 0101 => 65 - ADC - Zero Page
// 0b 0110 1001 => 69 - ADC - Immediate
// 0b 0110 1101 => 6D - ADC - Absolute
// 0b 0111 0001 => 71 - ADC - (Indirect),Y
// 0b 0111 0101 => 75 - ADC - Zero Page,X
// 0b 0111 1001 => 79 - ADC - Absolute,Y
// 0b 0111 1101 => 7D - ADC - Absolute,X

type Inst func(s *system)

var instructionSet = map[byte]Inst{
	0x00: BRK,
	0x01: ORA,
	0x02: NOP, // Future Expansion
	0x03: NOP, // Future Expansion
	0x04: NOP, // Future Expansion
	0x05: ORA,
	0x06: ASL,
	0x07: NOP, // Future Expansion
	0x08: PHP,
	0x09: ORA,
	0x0A: ASL,
	0x0B: NOP, // Future Expansion
	0x0C: NOP, // Future Expansion
	0x0D: ORA,
	0x0E: ASL,
	0x0F: NOP, // Future Expansion
	0x10: BPL,
	0x11: ORA,
	0x12: NOP, // Future Expansion
	0x13: NOP, // Future Expansion
	0x14: NOP, // Future Expansion
	0x15: ORA,
	0x16: ASL,
	0x17: NOP, // Future Expansion
	0x18: CLC,
	0x19: ORA,
	0x1A: NOP, // Future Expansion
	0x1B: NOP, // Future Expansion
	0x1C: NOP, // Future Expansion
	0x1D: ORA,
	0x1E: ASL,
	0x1F: NOP, // Future Expansion
	0x20: JSR,
	0x21: AND,
	0x22: NOP, // Future Expansion
	0x23: NOP, // Future Expansion
	0x24: BIT,
	0x25: AND,
	0x26: ROL,
	0x27: NOP, // Future Expansion
	0x28: PLP,
	0x29: AND,
	0x2A: ROL,
	0x2B: NOP, // Future Expansion
	0x2C: BIT,
	0x2D: AND,
	0x2E: ROL,
	0x2F: NOP, // Future Expansion
	0x30: BMI,
	0x31: AND,
	0x32: NOP, // Future Expansion
	0x33: NOP, // Future Expansion
	0x34: NOP, // Future Expansion
	0x35: AND,
	0x36: ROL,
	0x37: NOP, // Future Expansion
	0x38: SEC,
	0x39: AND,
	0x3A: NOP, // Future Expansion
	0x3B: NOP, // Future Expansion
	0x3C: NOP, // Future Expansion
	0x3D: AND,
	0x3E: ROL,
	0x3F: NOP, // Future Expansion
	0x40: RTI,
	0x41: EOR,
	0x42: NOP, // Future Expansion
	0x43: NOP, // Future Expansion
	0x44: NOP, // Future Expansion
	0x45: EOR,
	0x46: LSR,
	0x47: NOP, // Future Expansion
	0x48: PHA,
	0x49: EOR,
	0x4A: LSR,
	0x4B: NOP, // Future Expansion
	0x4C: JMP,
	0x4D: EOR,
	0x4E: LSR,
	0x4F: NOP, // Future Expansion
	0x50: BVC,
	0x51: EOR,
	0x52: NOP, // Future Expansion
	0x53: NOP, // Future Expansion
	0x54: NOP, // Future Expansion
	0x55: EOR,
	0x56: LSR,
	0x57: NOP, // Future Expansion
	0x58: CLI,
	0x59: EOR,
	0x5A: NOP, // Future Expansion
	0x5B: NOP, // Future Expansion
	0x5C: NOP, // Future Expansion
	0x5D: EOR,
	0x5E: LSR,
	0x5F: NOP, // Future Expansion
	0x60: RTS,
	0x61: ADC,
	0x62: NOP, // Future Expansion
	0x63: NOP, // Future Expansion
	0x64: NOP, // Future Expansion
	0x65: ADC,
	0x66: ROR,
	0x67: NOP, // Future Expansion
	0x68: PLA,
	0x69: ADC,
	0x6A: ROR,
	0x6B: NOP, // Future Expansion
	0x6C: JMP,
	0x6D: ADC,
	0x6E: ROR,
	0x6F: NOP, // Future Expansion
	0x70: BVS,
	0x71: ADC,
	0x72: NOP, // Future Expansion
	0x73: NOP, // Future Expansion
	0x74: NOP, // Future Expansion
	0x75: ADC,
	0x76: ROR,
	0x77: NOP, // Future Expansion
	0x78: SEI,
	0x79: ADC,
	0x7A: NOP, // Future Expansion
	0x7B: NOP, // Future Expansion
	0x7C: NOP, // Future Expansion
	0x7D: ADC,
	0x7E: ROR,
	0x7F: NOP, // Future Expansion
	0x80: NOP, // Future Expansion
	0x81: STA,
	0x82: NOP, // Future Expansion
	0x83: NOP, // Future Expansion
	0x84: STY,
	0x85: STA,
	0x86: STX,
	0x87: NOP, // Future Expansion
	0x88: DEY,
	0x89: NOP, // Future Expansion
	0x8A: TXA,
	0x8B: NOP, // Future Expansion
	0x8C: STY,
	0x8D: STA,
	0x8E: STX,
	0x8F: NOP, // Future Expansion
	0x90: BCC,
	0x91: STA,
	0x92: NOP, // Future Expansion
	0x93: NOP, // Future Expansion
	0x94: STY,
	0x95: STA,
	0x96: STX,
	0x97: NOP, // Future Expansion
	0x98: TYA,
	0x99: STA,
	0x9A: TXS,
	0x9B: NOP, // Future Expansion
	0x9C: NOP, // Future Expansion
	0x9D: STA,
	0x9E: NOP, // Future Expansion
	0x9F: NOP, // Future Expansion
	0xA0: LDY,
	0xA1: LDA,
	0xA2: LDX,
	0xA3: NOP, // Future Expansion
	0xA4: LDY,
	0xA5: LDA,
	0xA6: LDX,
	0xA7: NOP, // Future Expansion
	0xA8: TAY,
	0xA9: LDA,
	0xAA: TAX,
	0xAB: NOP, // Future Expansion
	0xAC: LDY,
	0xAD: LDA,
	0xAE: LDX,
	0xAF: NOP, // Future Expansion
	0xB0: BCS,
	0xB1: LDA,
	0xB2: NOP, // Future Expansion
	0xB3: NOP, // Future Expansion
	0xB4: LDY,
	0xB5: LDA,
	0xB6: LDX,
	0xB7: NOP, // Future Expansion
	0xB8: CLV,
	0xB9: LDA,
	0xBA: TSX,
	0xBB: NOP, // Future Expansion
	0xBC: LDY,
	0xBD: LDA,
	0xBE: LDX,
	0xBF: NOP, // Future Expansion
	0xC0: Cpy,
	0xC1: CMP,
	0xC2: NOP, // Future Expansion
	0xC3: NOP, // Future Expansion
	0xC4: CPY,
	0xC5: CMP,
	0xC6: DEC,
	0xC7: NOP, // Future Expansion
	0xC8: INY,
	0xC9: CMP,
	0xCA: DEX,
	0xCB: NOP, // Future Expansion
	0xCC: CPY,
	0xCD: CMP,
	0xCE: DEC,
	0xCF: NOP, // Future Expansion
	0xD0: BNE,
	0xD1: CMP,
	0xD2: NOP, // Future Expansion
	0xD3: NOP, // Future Expansion
	0xD4: NOP, // Future Expansion
	0xD5: CMP,
	0xD6: DEC,
	0xD7: NOP, // Future Expansion
	0xD8: CLD,
	0xD9: CMP,
	0xDA: NOP, // Future Expansion
	0xDB: NOP, // Future Expansion
	0xDC: NOP, // Future Expansion
	0xDD: CMP,
	0xDE: DEC,
	0xDF: NOP, // Future Expansion
	0xE0: CPX,
	0xE1: SBC,
	0xE2: NOP, // Future Expansion
	0xE3: NOP, // Future Expansion
	0xE4: CPX,
	0xE5: SBC,
	0xE6: INC,
	0xE7: NOP, // Future Expansion
	0xE8: INX,
	0xE9: SBC,
	0xEA: NOP,
	0xEB: NOP, // Future Expansion
	0xEC: CPX,
	0xED: SBC,
	0xEE: INC,
	0xEF: NOP, // Future Expansion
	0xF0: BEQ,
	0xF1: SBC,
	0xF2: NOP, // Future Expansion
	0xF3: NOP, // Future Expansion
	0xF4: NOP, // Future Expansion
	0xF5: SBC,
	0xF6: INC,
	0xF7: NOP, // Future Expansion
	0xF8: SED,
	0xF9: SBC,
	0xFA: NOP, // Future Expansion
	0xFB: NOP, // Future Expansion
	0xFC: NOP, // Future Expansion
	0xFD: SBC,
	0xFE: INC,
	0xFF: NOP, // Future Expansion
}

func (s *system) exec() {
	for inst := s.next(); inst != nil; inst = s.next() {
		if fn, ok := instructionSet[inst]; ok {
			// Found Instruction
			fn(s)
		}
		// Else: NOP
	}
}

const byteSize = 256

// ADC: Add with Carry
//
// Instruction:
// A + M + C -> A, C
//
// Addressing Modes:
// Immediate     | code: 69 | reads: 2 | clk: 2
// Zero Page     | code: 65 | reads: 2 | clk: 3
// Zero Page,X   | code: 75 | reads: 2 | clk: 4
// Absolute      | code: 60 | reads: 3 | clk: 4
// Absolute,X    | code: 70 | reads: 3 | clk: 4*
// Absolute,Y    | code: 79 | reads: 3 | clk: 4*
// (Indirect,X)  | code: 61 | reads: 2 | clk: 6
// (Indirect),Y  | code: 71 | reads: 2 | clk: 5*
func ADC(s *system) {
	carry := s.cpu.status.carry()
	acc := s.cpu.ac
	var (
		mem      byte
		res      byte
		rescarry byte
	)

	switch code := s.peek(); code {
	case 0x69: // Immediate
		mem = s.next()

	case 0x65: // Zero Page
		pageIndex := 0
		memIndex := s.next()
		mem = s.mem[pageIndex][memIndex]

	case 0x75: // Zero Page,X
		pageIndex := 0
		memIndex := s.next()
		xr := s.cpu.xr
		mem = s.mem[pageIndex][memIndex+xr]

	case 0x60: // Absolute
		// Little Endian
		memIndex := s.next()
		pageIndex := s.next()
		mem = s.mem[pageIndex][memIndex]

	case 0x70: // Absolute,X
		memIndex := s.next()
		pageIndex := s.next()
		xr := s.cpu.xr
		mem = s.mem[pageIndex][memIndex+xr]

	case 0x79: // Absolute,Y
		memIndex, pageIndex := s.next(), s.next()
		yr := s.cpu.yr
		mem = s.mem[pageIndex][memIndex+yr]

	case 0x61: // (Indirect,X)
		pageIndex := 0
		memIndex := s.next()
		xr := s.cpu.xr
		// Little Endian
		b1, b0 := s.mem[pageIndex][memIndex+xr], s.mem[pageIndex][memIndex+xr+1]
		mem = s.mem[b0][b1]

	case 0x71: // (Indirect),Y
		pageIndex := 0
		memIndex := s.next()
		// Little Endian
		b1, b0 := s.mem[pageIndex][memIndex], s.mem[pageIndex][memIndex+1]
		yr := s.cpu.yr
		b0 = (b0 + yr) % byteSize
		b1 = b1 + (b0+yr)/byteSize
		mem = s.mem[b0][b1]

	default: // Unknown: NOP
		return
	}

	res = (acc + mem + carry) % byteSize
	rescarry = (acc + mem + carry) / byteSize
	s.cpu.status.setCarry(rescarry)
	s.cpu.ac = res
}

func AND(s *system) {}

func ASL(s *system) {}

func BCC(s *system) {}

func BCS(s *system) {}

func BEQ(s *system) {}

func BIT(s *system) {}

func BMI(s *system) {}

func BNE(s *system) {}

func BPL(s *system) {}

func BRK(s *system) {}

func BVC(s *system) {}

func BVS(s *system) {}

func CLC(s *system) {}

func CLD(s *system) {}

func CLI(s *system) {}

func CLV(s *system) {}

func CMP(s *system) {}

func CPX(s *system) {}

func CPY(s *system) {}

func DEC(s *system) {}

func DEX(s *system) {}

func DEY(s *system) {}

func EOR(s *system) {}

func INC(s *system) {}

func INX(s *system) {}

func INY(s *system) {}

func JMP(s *system) {}

func JSR(s *system) {}

func LDA(s *system) {}

func LDX(s *system) {}

func LDY(s *system) {}

func LSR(s *system) {}

func NOP(s *system) {}

func ORA(s *system) {}

func PHA(s *system) {}

func PHP(s *system) {}

func PLA(s *system) {}

func PLP(s *system) {}

func ROL(s *system) {}

func ROR(s *system) {}

func RTI(s *system) {}

func RTS(s *system) {}

func SBC(s *system) {}

func SEC(s *system) {}

func SED(s *system) {}

func SEI(s *system) {}

func STA(s *system) {}

func STX(s *system) {}

func STY(s *system) {}

func TAX(s *system) {}

func TAY(s *system) {}

func TSX(s *system) {}

func TXA(s *system) {}

func TXS(s *system) {}

func TYA(s *system) {}

// func (s *system) exec() {
// 	switch i := s.next(); i {
// 	case 0x00: // 00 - BRK
// 	case 0x01: // 01 - ORA - (Indirect,X)
// 	case 0x02: // 02 - Future Expansion
// 	case 0x03: // 03 - Future Expansion
// 	case 0x04: // 04 - Future Expansion
// 	case 0x05: // 05 - ORA - Zero Page
// 	case 0x06: // 06 - ASL - Zero Page
// 	case 0x07: // 07 - Future Expansion
// 	case 0x08: // 08 - PHP
// 	case 0x09: // 09 - ORA - Immediate
// 	case 0x0A: // 0A - ASL - Accumulator
// 	case 0x0B: // 0B - Future Expansion
// 	case 0x0C: // 0C - Future Expansion
// 	case 0x0D: // 0D - ORA - Absolute
// 	case 0x0E: // 0E - ASL - Absolute
// 	case 0x0F: // 0F - Future Expansion
//
// 	case 0x10: // 10 - BPL
// 	case 0x11: // 11 - ORA - (Indirect),Y
// 	case 0x12: // 12 - Future Expansion
// 	case 0x13: // 13 - Future Expansion
// 	case 0x14: // 14 - Future Expansion
// 	case 0x15: // 15 - ORA - Zero Page,X
// 	case 0x16: // 16 - ASL - Zero Page,X
// 	case 0x17: // 17 - Future Expansion
// 	case 0x18: // 18 - CLC
// 	case 0x19: // 19 - ORA - Absolute,Y
// 	case 0x1A: // 1A - Future Expansion
// 	case 0x1B: // 1B - Future Expansion
// 	case 0x1C: // 1C - Future Expansion
// 	case 0x1D: // 1D - ORA - Absolute,X
// 	case 0x1E: // 1E - ASL - Absolute,X
// 	case 0x1F: // 1F - Future Expansion
//
// 	case 0x20: // 20 - JSR
// 	case 0x21: // 21 - AND - (Indirect,X)
// 	case 0x22: // 22 - Future Expansion
// 	case 0x23: // 23 - Future Expansion
// 	case 0x24: // 24 - BIT - Zero Page
// 	case 0x25: // 25 - AND - Zero Page
// 	case 0x26: // 26 - ROL - Zero Page
// 	case 0x27: // 27 - Future Expansion
// 	case 0x28: // 28 - PLP
// 	case 0x29: // 29 - AND - Immediate
// 	case 0x2A: // 2A - ROL - Accumulator
// 	case 0x2B: // 2B - Future Expansion
// 	case 0x2C: // 2C - BIT - Absolute
// 	case 0x2D: // 2D - AND - Absolute
// 	case 0x2E: // 2E - ROL - Absolute
// 	case 0x2F: // 2F - Future Expansion
//
// 	case 0x30: // 30 - BMI
// 	case 0x31: // 31 - AND - (Indirect),Y
// 	case 0x32: // 32 - Future Expansion
// 	case 0x33: // 33 - Future Expansion
// 	case 0x34: // 34 - Future Expansion
// 	case 0x35: // 35 - AND - Zero Page,X
// 	case 0x36: // 36 - ROL - Zero Page,X
// 	case 0x37: // 37 - Future Expansion
// 	case 0x38: // 38 - SEC
// 	case 0x39: // 39 - AND - Absolute,Y
// 	case 0x3A: // 3A - Future Expansion
// 	case 0x3B: // 3B - Future Expansion
// 	case 0x3C: // 3C - Future Expansion
// 	case 0x3D: // 3D - AND - Absolute,X
// 	case 0x3E: // 3E - ROL - Absolute,X
// 	case 0x3F: // 3F - Future Expansion
//
// 	case 0x40: // 40 - RTI
// 	case 0x41: // 41 - EOR - (Indirect,X)
// 	case 0x42: // 42 - Future Expansion
// 	case 0x43: // 43 - Future Expansion
// 	case 0x44: // 44 - Future Expansion
// 	case 0x45: // 45 - EOR - Zero Page
// 	case 0x46: // 46 - LSR - Zero Page
// 	case 0x47: // 47 - Future Expansion
// 	case 0x48: // 48 - PHA
// 	case 0x49: // 49 - EOR - Immediate
// 	case 0x4A: // 4A - LSR - Accumulator
// 	case 0x4B: // 4B - Future Expansion
// 	case 0x4C: // 4C - JMP - Absolute
// 	case 0x4D: // 4D - EOR - Absolute
// 	case 0x4E: // 4E - LSR - Absolute
// 	case 0x4F: // 4F - Future Expansion
//
// 	case 0x50: // 50 - BVC
// 	case 0x51: // 51 - EOR - (Indirect),Y
// 	case 0x52: // 52 - Future Expansion
// 	case 0x53: // 53 - Future Expansion
// 	case 0x54: // 54 - Future Expansion
// 	case 0x55: // 55 - EOR - Zero Page,X
// 	case 0x56: // 56 - LSR - Zero Page,X
// 	case 0x57: // 57 - Future Expansion
// 	case 0x58: // 58 - CLI
// 	case 0x59: // 59 - EOR - Absolute,Y
// 	case 0x5A: // 5A - Future Expansion
// 	case 0x5B: // 5B - Future Expansion
// 	case 0x5C: // 5C - Future Expansion
// 	// case 0x50: // 50 - EOR - Absolute,X
// 	case 0x5D: // 50 - EOR - Absolute,X
// 	case 0x5E: // 5E - LSR - Absolute,X
// 	case 0x5F: // 5F - Future Expansion
//
// 	case 0x60: // 60 - RTS
// 	case 0x61: // 61 - ADC - (Indirect,X)
// 	case 0x62: // 62 - Future Expansion
// 	case 0x63: // 63 - Future Expansion
// 	case 0x64: // 64 - Future Expansion
// 	case 0x65: // 65 - ADC - Zero Page
// 	case 0x66: // 66 - ROR - Zero Page
// 	case 0x67: // 67 - Future Expansion
// 	case 0x68: // 68 - PLA
// 	case 0x69: // 69 - ADC - Immediate
// 	case 0x6A: // 6A - ROR - Accumulator
// 	case 0x6B: // 6B - Future Expansion
// 	case 0x6C: // 6C - JMP - Indirect
// 	case 0x6D: // 6D - ADC - Absolute
// 	case 0x6E: // 6E - ROR - Absolute
// 	case 0x6F: // 6F - Future Expansion
//
// 	case 0x70: // 70 - BVS
// 	case 0x71: // 71 - ADC - (Indirect),Y
// 	case 0x72: // 72 - Future Expansion
// 	case 0x73: // 73 - Future Expansion
// 	case 0x74: // 74 - Future Expansion
// 	case 0x75: // 75 - ADC - Zero Page,X
// 	case 0x76: // 76 - ROR - Zero Page,X
// 	case 0x77: // 77 - Future Expansion
// 	case 0x78: // 78 - SEI
// 	case 0x79: // 79 - ADC - Absolute,Y
// 	case 0x7A: // 7A - Future Expansion
// 	case 0x7B: // 7B - Future Expansion
// 	case 0x7C: // 7C - Future Expansion
// 	case 0x7D: // 70 - ADC - Absolute,X
// 	case 0x7E: // 7E - ROR - Absolute,X
// 	case 0x7F: // 7F - Future Expansion
//
// 	case 0x80: // 80 - Future Expansion
// 	case 0x81: // 81 - STA - (Indirect,X)
// 	case 0x82: // 82 - Future Expansion
// 	case 0x83: // 83 - Future Expansion
// 	case 0x84: // 84 - STY - Zero Page
// 	case 0x85: // 85 - STA - Zero Page
// 	case 0x86: // 86 - STX - Zero Page
// 	case 0x87: // 87 - Future Expansion
// 	case 0x88: // 88 - DEY
// 	case 0x89: // 89 - Future Expansion
// 	case 0x8A: // 8A - TXA
// 	case 0x8B: // 8B - Future Expansion
// 	case 0x8C: // 8C - STY - Absolute
// 	case 0x8D: // 80 - STA - Absolute
// 	case 0x8E: // 8E - STX - Absolute
// 	case 0x8F: // 8F - Future Expansion
//
// 	case 0x90: // 90 - BCC
// 	case 0x91: // 91 - STA - (Indirect),Y
// 	case 0x92: // 92 - Future Expansion
// 	case 0x93: // 93 - Future Expansion
// 	case 0x94: // 94 - STY - Zero Page,X
// 	case 0x95: // 95 - STA - Zero Page,X
// 	case 0x96: // 96 - STX - Zero Page,Y
// 	case 0x97: // 97 - Future Expansion
// 	case 0x98: // 98 - TYA
// 	case 0x99: // 99 - STA - Absolute,Y
// 	case 0x9A: // 9A - TXS
// 	case 0x9B: // 9B - Future Expansion
// 	case 0x9C: // 9C - Future Expansion
// 	case 0x9D: // 90 - STA - Absolute,X
// 	case 0x9E: // 9E - Future Expansion
// 	case 0x9F: // 9F - Future Expansion
//
// 	case 0xA0: // A0 - LDY - Immediate
// 	case 0xA1: // A1 - LDA - (Indirect,X)
// 	case 0xA2: // A2 - LDX - Immediate
// 	case 0xA3: // A3 - Future Expansion
// 	case 0xA4: // A4 - LDY - Zero Page
// 	case 0xA5: // A5 - LDA - Zero Page
// 	case 0xA6: // A6 - LDX - Zero Page
// 	case 0xA7: // A7 - Future Expansion
// 	case 0xA8: // A8 - TAY
// 	case 0xA9: // A9 - LDA - Immediate
// 	case 0xAA: // AA - TAX
// 	case 0xAB: // AB - Future Expansion
// 	case 0xAC: // AC - LDY - Absolute
// 	case 0xAD: // AD - LDA - Absolute
// 	case 0xAE: // AE - LDX - Absolute
// 	case 0xAF: // AF - Future Expansion
//
// 	case 0xB0: // B0 - BCS
// 	case 0xB1: // B1 - LDA - (Indirect),Y
// 	case 0xB2: // B2 - Future Expansion
// 	case 0xB3: // B3 - Future Expansion
// 	case 0xB4: // B4 - LDY - Zero Page,X
// 	case 0xB5: // BS - LDA - Zero Page,X
// 	case 0xB6: // B6 - LDX - Zero Page,Y
// 	case 0xB7: // B7 - Future Expansion
// 	case 0xB8: // B8 - CLV
// 	case 0xB9: // B9 - LDA - Absolute,Y
// 	case 0xBA: // BA - TSX
// 	case 0xBB: // BB - Future Expansion
// 	case 0xBC: // BC - LDY - Absolute,X
// 	case 0xBD: // BD - LDA - Absolute,X
// 	case 0xBE: // BE - LDX - Absolute,Y
// 	case 0xBF: // BF - Future Expansion
//
// 	case 0xC0: // C0 - Cpy - Immediate
// 	case 0xC1: // C1 - CMP - (Indirect,X)
// 	case 0xC2: // C2 - Future Expansion
// 	case 0xC3: // C3 - Future Expansion
// 	case 0xC4: // C4 - CPY - Zero Page
// 	case 0xC5: // C5 - CMP - Zero Page
// 	case 0xC6: // C6 - DEC - Zero Page
// 	case 0xC7: // C7 - Future Expansion
// 	case 0xC8: // C8 - INY
// 	case 0xC9: // C9 - CMP - Immediate
// 	case 0xCA: // CA - DEX
// 	case 0xCB: // CB - Future Expansion
// 	case 0xCC: // CC - CPY - Absolute
// 	case 0xCD: // CD - CMP - Absolute
// 	case 0xCE: // CE - DEC - Absolute
// 	case 0xCF: // CF - Future Expansion
//
// 	case 0xD0: // D0 - BNE
// 	case 0xD1: // D1 - CMP   (Indirect@,Y
// 	case 0xD2: // D2 - Future Expansion
// 	case 0xD3: // D3 - Future Expansion
// 	case 0xD4: // D4 - Future Expansion
// 	case 0xD5: // D5 - CMP - Zero Page,X
// 	case 0xD6: // D6 - DEC - Zero Page,X
// 	case 0xD7: // D7 - Future Expansion
// 	case 0xD8: // D8 - CLD
// 	case 0xD9: // D9 - CMP - Absolute,Y
// 	case 0xDA: // DA - Future Expansion
// 	case 0xDB: // DB - Future Expansion
// 	case 0xDC: // DC - Future Expansion
// 	case 0xDD: // DD - CMP - Absolute,X
// 	case 0xDE: // DE - DEC - Absolute,X
// 	case 0xDF: // DF - Future Expansion
//
// 	case 0xE0: // E0 - CPX - Immediate
// 	case 0xE1: // E1 - SBC - (Indirect,X)
// 	case 0xE2: // E2 - Future Expansion
// 	case 0xE3: // E3 - Future Expansion
// 	case 0xE4: // E4 - CPX - Zero Page
// 	case 0xE5: // E5 - SBC - Zero Page
// 	case 0xE6: // E6 - INC - Zero Page
// 	case 0xE7: // E7 - Future Expansion
// 	case 0xE8: // E8 - INX
// 	case 0xE9: // E9 - SBC - Immediate
// 	case 0xEA: // EA - NOP
// 	case 0xEB: // EB - Future Expansion
// 	case 0xEC: // EC - CPX - Absolute
// 	case 0xED: // ED - SBC - Absolute
// 	case 0xEE: // EE - INC - Absolute
// 	case 0xEF: // EF - Future Expansion
//
// 	case 0xF0: // F0 - BEQ
// 	case 0xF1: // F1 - SBC - (Indirect),Y
// 	case 0xF2: // F2 - Future Expansion
// 	case 0xF3: // F3 - Future Expansion
// 	case 0xF4: // F4 - Future Expansion
// 	case 0xF5: // F5 - SBC - Zero Page,X
// 	case 0xF6: // F6 - INC - Zero Page,X
// 	case 0xF7: // F7 - Future Expansion
// 	case 0xF8: // F8 - SED
// 	case 0xF9: // F9 - SBC - Absolute,Y
// 	case 0xFA: // FA - Future Expansion
// 	case 0xFB: // FB - Future Expansion
// 	case 0xFC: // FC - Future Expansion
// 	case 0xFD: // FD - SBC - Absolute,X
// 	case 0xFE: // FE - INC - Absolute,X
// 	case 0xFF: // FF - Future Expansion
//
// 	default: // Invalid Operation
// 		break
// 	}
// }
