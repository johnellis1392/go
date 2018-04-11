package main

import (
	"fmt"
)

func testLexer() {
	l := lex(`
		base = {
			workspace = jhelli_FintechICEInvoiceService;
			versionSet = "FintechICEInvoice/jhelli";
		};
		packages = {
			FintechICEInvoiceService-1.0 = .;
			FintechICEInvoiceServiceModel-1.0 = .;
			FintechICEInvoiceServiceClientConfig-1.0 = .;
		};
	`)

	var items []item
loop:
	for {
		select {
		case i, ok := <-l.items:
			if !ok {
				// Channel Closed Unexpectedly
				if i.typ == itemEOF {
					break loop
				}
				panic(fmt.Errorf("Channel closed before EOF was reached"))
			}

			switch i.typ {
			case itemEOF:
				// Done
				// fmt.Println("Found EOF")
				items = append(items, i)
				break loop
			case itemError:
				fmt.Println(i.val)
				panic(fmt.Errorf("A lexer error occurred: %v", i))
			default:
				// fmt.Printf("Found Token: %q\n", i)
				items = append(items, i)
				continue
			}
		}
	}

	fmt.Printf("Found %d Tokens:\n", len(items))
	for _, i := range items {
		fmt.Printf(" * %q\n", i)
	}
}

func main() {
	testLexer()
}
