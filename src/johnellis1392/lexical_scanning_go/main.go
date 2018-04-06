package main

import (
	"fmt"
)

func main() {
	fmt.Println("Lexing...")
	_, ichan := lex("ExampleLexer", `Something {{1}} Else`)

	var items []item
	// loop:
	// 	for {
	// 		select {
	// 		case i := <-ichan:
	// 			if i.typ == itemEOF {
	// 				// goto END
	// 				break loop
	// 			}
	// 			items = append(items, i)
	// 			break
	// 		default:
	// 			// Nothing received from channel; continue
	// 			fmt.Println("...(nothing received from ichan)")
	// 			continue
	// 		}
	// 	}

loop:
	for {
		select {
		case i, ok := <-ichan:

			if !ok {
				// Channel Closed; Haven't found EOF
				panic(fmt.Errorf("an error occurred while lexing: channel closed early"))
			} else if i.typ == itemEOF {
				// Finished
				break loop
			} else {
				// Other
				items = append(items, i)
			}

			// if ok {
			// 	continue
			// } else if i.typ == itemEOF {
			// 	break loop
			// } else {
			// 	// An error occurred; channel closed without encountering EOF
			// 	fmt.Println("\n\n")
			// 	fmt.Printf("Received Item: %v, %q\n", i, i)
			// 	fmt.Printf("Received Item (string): %q\n", i.String())
			// 	fmt.Printf(" * typ: %q, %v\n", i.typ, i.typ)
			// 	fmt.Printf(" * val: %q, %v\n", i.val, i.val)
			// 	fmt.Printf(" * ok:  %q, %v\n", ok, ok)
			// 	fmt.Println("\n\n")
			// 	panic(fmt.Errorf("an error occurred while lexing"))
			// }

		}
	}

	// END:
	fmt.Println("Received Tokens:")
	fmt.Println(items)
}
