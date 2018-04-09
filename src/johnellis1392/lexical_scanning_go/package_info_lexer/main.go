package main

import (
	"fmt"
)

func runLexer(input string) {
	var items []item
	itemchan := lex(input)

loop:
	for {
		select {
		case i, open := <-itemchan:
			if !open {
				break loop
			}

			// Channel Open
			switch i.typ {
			case itemError:
				panic(fmt.Sprintf("received error: %s", i.val))
			case itemEOF:
				fmt.Println("Received EOF")
			default:
				items = append(items, i)
			}
		}
	}

	fmt.Println("Found Tokens:")
	for _, i := range items {
		fmt.Println(" *", i)
	}
}

func lexAndParse(input string) {
	itemchan := lex(input)
	parsechan := parse(itemchan)

	var nodes []node
	for n := range parsechan {
		nodes = append(nodes, n)
	}

	fmt.Printf("Received (%v) Nodes:\n", len(nodes))
	for _, n := range nodes {
		fmt.Println(" *", n)
	}
}

func main() {
	fmt.Println("Starting Lexer...")
	const input = `
		base = {
		  workspace = jhelli_FintechICEIngestionService;
		  versionSet = "FintechICEIngestionService/jhelli";
		};
		packages = {
		  FintechICEIngestionServiceModel-1.0 = .;
		  FintechICEIngestionServiceClientConfig-1.1 = .;
		  FintechICEIngestionServiceJavaClient-1.1 = .;
		  FintechICEIngestionServiceTests-1.0 = .;
		  FintechICEIngestionService-1.0 = .;
		  FintechICEIngestionServiceCustomTraits-1.0 = .;
		};
	`

	runLexer(input)
}
