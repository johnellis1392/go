package main

import (
	"fmt"
	"io/ioutil"
)

func chunk(input []byte) []int16 {
	var res []int16
	for i := 0; i < len(input); i += 2 {
		res[i/2] = int16(input[i]<<8) | int16(input[i+1])
	}
	return res
}

// ReadProgram reads a nes file specified by filename and returns
// a Program representing the data found in the file.
func ReadProgram(filename string) (Program, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Program{}, err
	}

	// TODO: Handle any headers or blocks that may be found in these files.

	var insts []Inst
	vs := chunk(data)
	for _, b := range vs {
		// TODO: Handle remainder of instructions
		switch b {
		case 0x01:
			break
		default:
			// Unknown Instruction
			break
		}
	}

	// TODO: Read Data
	// TODO: Read remainder of segments

	return Program{
		Insts: insts,
		Raw:   vs,
	}, nil
}

func main() {
	const filename = "./test.txt"
	fmt.Printf("Reading file: %s\n", filename)
	p, err := ReadProgram(filename)
	if err != nil {
		panic(err)
	}

	s := NewSystem(p)
	for inst := p.Insts[p.Ptr]; inst != nil; {
		fmt.Printf("Executing Instruction: %v\n", inst.String())
		inst.Exec(s)
	}

	fmt.Println("Done")
}
