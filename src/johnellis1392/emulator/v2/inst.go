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

func (s *system) exec() {
	switch i := s.next(); i {
	case 0x00: // 00 - BRK
		// NOTE: Cannot be Shut Off by Setting status.I
		// Addressing Mode: Implied
		// Operation: Forced Interrupt PC + 2 toS P toS
		break

	case 0x01: // 01 - ORA - (Indirect,X)
	case 0x02: // 02 - Future Expansion
	case 0x03: // 03 - Future Expansion
	case 0x04: // 04 - Future Expansion
	case 0x05: // 05 - ORA - Zero Page
	case 0x06: // 06 - ASL - Zero Page
	case 0x07: // 07 - Future Expansion
	case 0x08: // 08 - PHP
	case 0x09: // 09 - ORA - Immediate
	case 0x0A: // 0A - ASL - Accumulator
	case 0x0B: // 0B - Future Expansion
	case 0x0C: // 0C - Future Expansion
	case 0x0D: // 0D - ORA - Absolute
	case 0x0E: // 0E - ASL - Absolute
	case 0x0F: // 0F - Future Expansion

	case 0x10: // 10 - BPL
	case 0x11: // 11 - ORA - (Indirect),Y
	case 0x12: // 12 - Future Expansion
	case 0x13: // 13 - Future Expansion
	case 0x14: // 14 - Future Expansion
	case 0x15: // 15 - ORA - Zero Page,X
	case 0x16: // 16 - ASL - Zero Page,X
	case 0x17: // 17 - Future Expansion
	case 0x18: // 18 - CLC
	case 0x19: // 19 - ORA - Absolute,Y
	case 0x1A: // 1A - Future Expansion
	case 0x1B: // 1B - Future Expansion
	case 0x1C: // 1C - Future Expansion
	case 0x1D: // 1D - ORA - Absolute,X
	case 0x1E: // 1E - ASL - Absolute,X
	case 0x1F: // 1F - Future Expansion

	case 0x20: // 20 - JSR
	case 0x21: // 21 - AND - (Indirect,X)
	case 0x22: // 22 - Future Expansion
	case 0x23: // 23 - Future Expansion
	case 0x24: // 24 - BIT - Zero Page
	case 0x25: // 25 - AND - Zero Page
	case 0x26: // 26 - ROL - Zero Page
	case 0x27: // 27 - Future Expansion
	case 0x28: // 28 - PLP
	case 0x29: // 29 - AND - Immediate
	case 0x2A: // 2A - ROL - Accumulator
	case 0x2B: // 2B - Future Expansion
	case 0x2C: // 2C - BIT - Absolute
	case 0x2D: // 2D - AND - Absolute
	case 0x2E: // 2E - ROL - Absolute
	case 0x2F: // 2F - Future Expansion

	case 0x30: // 30 - BMI
	case 0x31: // 31 - AND - (Indirect),Y
	case 0x32: // 32 - Future Expansion
	case 0x33: // 33 - Future Expansion
	case 0x34: // 34 - Future Expansion
	case 0x35: // 35 - AND - Zero Page,X
	case 0x36: // 36 - ROL - Zero Page,X
	case 0x37: // 37 - Future Expansion
	case 0x38: // 38 - SEC
	case 0x39: // 39 - AND - Absolute,Y
	case 0x3A: // 3A - Future Expansion
	case 0x3B: // 3B - Future Expansion
	case 0x3C: // 3C - Future Expansion
	case 0x3D: // 3D - AND - Absolute,X
	case 0x3E: // 3E - ROL - Absolute,X
	case 0x3F: // 3F - Future Expansion

	case 0x40: // 40 - RTI
	case 0x41: // 41 - EOR - (Indirect,X)
	case 0x42: // 42 - Future Expansion
	case 0x43: // 43 - Future Expansion
	case 0x44: // 44 - Future Expansion
	case 0x45: // 45 - EOR - Zero Page
	case 0x46: // 46 - LSR - Zero Page
	case 0x47: // 47 - Future Expansion
	case 0x48: // 48 - PHA
	case 0x49: // 49 - EOR - Immediate
	case 0x4A: // 4A - LSR - Accumulator
	case 0x4B: // 4B - Future Expansion
	case 0x4C: // 4C - JMP - Absolute
	case 0x4D: // 4D - EOR - Absolute
	case 0x4E: // 4E - LSR - Absolute
	case 0x4F: // 4F - Future Expansion

	case 0x50: // 50 - BVC
	case 0x51: // 51 - EOR - (Indirect),Y
	case 0x52: // 52 - Future Expansion
	case 0x53: // 53 - Future Expansion
	case 0x54: // 54 - Future Expansion
	case 0x55: // 55 - EOR - Zero Page,X
	case 0x56: // 56 - LSR - Zero Page,X
	case 0x57: // 57 - Future Expansion
	case 0x58: // 58 - CLI
	case 0x59: // 59 - EOR - Absolute,Y
	case 0x5A: // 5A - Future Expansion
	case 0x5B: // 5B - Future Expansion
	case 0x5C: // 5C - Future Expansion
	// case 0x50: // 50 - EOR - Absolute,X
	case 0x5D: // 50 - EOR - Absolute,X
	case 0x5E: // 5E - LSR - Absolute,X
	case 0x5F: // 5F - Future Expansion

	case 0x60: // 60 - RTS
	case 0x61: // 61 - ADC - (Indirect,X)
		var carry uint8

		// Get Carry Flag
		if s.cpu.sr.c {
			carry = 1
		} else {
			carry = 0
		}

		temp := int8(i + s.cpu.ac + carry)

		// Set Zero Flag
		s.cpu.sr.z = temp == 0

		// Decimal Flag Set
		if s.cpu.sr.d && (s.cpu.ac&0xF+i&0xF+carry) > 9 {
			temp += 6
		}

		// Set Sign Flag
		s.cpu.sr.s = temp < 0

		// Set Overflow Flag
		// TODO

	case 0x62: // 62 - Future Expansion
	case 0x63: // 63 - Future Expansion
	case 0x64: // 64 - Future Expansion
	case 0x65: // 65 - ADC - Zero Page
	case 0x66: // 66 - ROR - Zero Page
	case 0x67: // 67 - Future Expansion
	case 0x68: // 68 - PLA
	case 0x69: // 69 - ADC - Immediate
	case 0x6A: // 6A - ROR - Accumulator
	case 0x6B: // 6B - Future Expansion
	case 0x6C: // 6C - JMP - Indirect
	case 0x6D: // 6D - ADC - Absolute
	case 0x6E: // 6E - ROR - Absolute
	case 0x6F: // 6F - Future Expansion

	case 0x70: // 70 - BVS
	case 0x71: // 71 - ADC - (Indirect),Y
	case 0x72: // 72 - Future Expansion
	case 0x73: // 73 - Future Expansion
	case 0x74: // 74 - Future Expansion
	case 0x75: // 75 - ADC - Zero Page,X
	case 0x76: // 76 - ROR - Zero Page,X
	case 0x77: // 77 - Future Expansion
	case 0x78: // 78 - SEI
	case 0x79: // 79 - ADC - Absolute,Y
	case 0x7A: // 7A - Future Expansion
	case 0x7B: // 7B - Future Expansion
	case 0x7C: // 7C - Future Expansion
	// case 0x70: // 70 - ADC - Absolute,X
	case 0x7D: // 70 - ADC - Absolute,X
	case 0x7E: // 7E - ROR - Absolute,X
	case 0x7F: // 7F - Future Expansion

	case 0x80: // 80 - Future Expansion
	case 0x81: // 81 - STA - (Indirect,X)
	case 0x82: // 82 - Future Expansion
	case 0x83: // 83 - Future Expansion
	case 0x84: // 84 - STY - Zero Page
	case 0x85: // 85 - STA - Zero Page
	case 0x86: // 86 - STX - Zero Page
	case 0x87: // 87 - Future Expansion
	case 0x88: // 88 - DEY
	case 0x89: // 89 - Future Expansion
	case 0x8A: // 8A - TXA
	case 0x8B: // 8B - Future Expansion
	case 0x8C: // 8C - STY - Absolute
	// case 0x80: // 80 - STA - Absolute
	case 0x8D: // 80 - STA - Absolute
	case 0x8E: // 8E - STX - Absolute
	case 0x8F: // 8F - Future Expansion

	case 0x90: // 90 - BCC
	case 0x91: // 91 - STA - (Indirect),Y
	case 0x92: // 92 - Future Expansion
	case 0x93: // 93 - Future Expansion
	case 0x94: // 94 - STY - Zero Page,X
	case 0x95: // 95 - STA - Zero Page,X
	case 0x96: // 96 - STX - Zero Page,Y
	case 0x97: // 97 - Future Expansion
	case 0x98: // 98 - TYA
	case 0x99: // 99 - STA - Absolute,Y
	case 0x9A: // 9A - TXS
	case 0x9B: // 9B - Future Expansion
	case 0x9C: // 9C - Future Expansion
	// case 0x90: // 90 - STA - Absolute,X
	case 0x9D: // 90 - STA - Absolute,X
	case 0x9E: // 9E - Future Expansion
	case 0x9F: // 9F - Future Expansion

	case 0xA0: // A0 - LDY - Immediate
	case 0xA1: // A1 - LDA - (Indirect,X)
	case 0xA2: // A2 - LDX - Immediate
	case 0xA3: // A3 - Future Expansion
	case 0xA4: // A4 - LDY - Zero Page
	case 0xA5: // A5 - LDA - Zero Page
	case 0xA6: // A6 - LDX - Zero Page
	case 0xA7: // A7 - Future Expansion
	case 0xA8: // A8 - TAY
	case 0xA9: // A9 - LDA - Immediate
	case 0xAA: // AA - TAX
	case 0xAB: // AB - Future Expansion
	case 0xAC: // AC - LDY - Absolute
	case 0xAD: // AD - LDA - Absolute
	case 0xAE: // AE - LDX - Absolute
	case 0xAF: // AF - Future Expansion

	case 0xB0: // B0 - BCS
	case 0xB1: // B1 - LDA - (Indirect),Y
	case 0xB2: // B2 - Future Expansion
	case 0xB3: // B3 - Future Expansion
	case 0xB4: // B4 - LDY - Zero Page,X
	// case 0xBS: // BS - LDA - Zero Page,X
	case 0xB5: // BS - LDA - Zero Page,X
	case 0xB6: // B6 - LDX - Zero Page,Y
	case 0xB7: // B7 - Future Expansion
	case 0xB8: // B8 - CLV
	case 0xB9: // B9 - LDA - Absolute,Y
	case 0xBA: // BA - TSX
	case 0xBB: // BB - Future Expansion
	case 0xBC: // BC - LDY - Absolute,X
	case 0xBD: // BD - LDA - Absolute,X
	case 0xBE: // BE - LDX - Absolute,Y
	case 0xBF: // BF - Future Expansion

	case 0xC0: // C0 - Cpy - Immediate
	case 0xC1: // C1 - CMP - (Indirect,X)
	case 0xC2: // C2 - Future Expansion
	case 0xC3: // C3 - Future Expansion
	case 0xC4: // C4 - CPY - Zero Page
	case 0xC5: // C5 - CMP - Zero Page
	case 0xC6: // C6 - DEC - Zero Page
	case 0xC7: // C7 - Future Expansion
	case 0xC8: // C8 - INY
	case 0xC9: // C9 - CMP - Immediate
	case 0xCA: // CA - DEX
	case 0xCB: // CB - Future Expansion
	case 0xCC: // CC - CPY - Absolute
	case 0xCD: // CD - CMP - Absolute
	case 0xCE: // CE - DEC - Absolute
	case 0xCF: // CF - Future Expansion

	case 0xD0: // D0 - BNE
	case 0xD1: // D1 - CMP   (Indirect@,Y
	case 0xD2: // D2 - Future Expansion
	case 0xD3: // D3 - Future Expansion
	case 0xD4: // D4 - Future Expansion
	case 0xD5: // D5 - CMP - Zero Page,X
	case 0xD6: // D6 - DEC - Zero Page,X
	case 0xD7: // D7 - Future Expansion
	case 0xD8: // D8 - CLD
	case 0xD9: // D9 - CMP - Absolute,Y
	case 0xDA: // DA - Future Expansion
	case 0xDB: // DB - Future Expansion
	case 0xDC: // DC - Future Expansion
	case 0xDD: // DD - CMP - Absolute,X
	case 0xDE: // DE - DEC - Absolute,X
	case 0xDF: // DF - Future Expansion

	case 0xE0: // E0 - CPX - Immediate
	case 0xE1: // E1 - SBC - (Indirect,X)
	case 0xE2: // E2 - Future Expansion
	case 0xE3: // E3 - Future Expansion
	case 0xE4: // E4 - CPX - Zero Page
	case 0xE5: // E5 - SBC - Zero Page
	case 0xE6: // E6 - INC - Zero Page
	case 0xE7: // E7 - Future Expansion
	case 0xE8: // E8 - INX
	case 0xE9: // E9 - SBC - Immediate
	case 0xEA: // EA - NOP
	case 0xEB: // EB - Future Expansion
	case 0xEC: // EC - CPX - Absolute
	case 0xED: // ED - SBC - Absolute
	case 0xEE: // EE - INC - Absolute
	case 0xEF: // EF - Future Expansion

	case 0xF0: // F0 - BEQ
	case 0xF1: // F1 - SBC - (Indirect),Y
	case 0xF2: // F2 - Future Expansion
	case 0xF3: // F3 - Future Expansion
	case 0xF4: // F4 - Future Expansion
	case 0xF5: // F5 - SBC - Zero Page,X
	case 0xF6: // F6 - INC - Zero Page,X
	case 0xF7: // F7 - Future Expansion
	case 0xF8: // F8 - SED
	case 0xF9: // F9 - SBC - Absolute,Y
	case 0xFA: // FA - Future Expansion
	case 0xFB: // FB - Future Expansion
	case 0xFC: // FC - Future Expansion
	case 0xFD: // FD - SBC - Absolute,X
	case 0xFE: // FE - INC - Absolute,X
	case 0xFF: // FF - Future Expansion

	default: // Invalid Operation
		break
	}
}
