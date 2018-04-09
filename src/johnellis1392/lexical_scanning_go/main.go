package main

import (
	"fmt"
)

func main() {
	fmt.Println("Lexing...")
	_, ichan := lex("ExampleLexer", `Something {{1}} Else`)
	var items []item

loop:
	for {
		select {
		case i, open := <-ichan:
			if !open {
				// Channel Closed
				if i.typ == itemEOF {
					panic(fmt.Errorf("channel closed unexpectedly, with terminal item: %v", i))
				}
				// Success
				break loop
			} else if i.typ == itemEOF {
				// Finished
				break loop
			} else {
				// Other
				items = append(items, i)
			}
		}
	}

	fmt.Println("Received Tokens:")
	fmt.Println(items)
}
