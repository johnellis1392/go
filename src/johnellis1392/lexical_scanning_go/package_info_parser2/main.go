package main

import (
	"fmt"
)

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

func withLexer() {
	tokchan := lex(input)

	var tokens []token
	for t := range tokchan {
		fmt.Println(" *** Received Token: ", t)
		tokens = append(tokens, t)
	}

	fmt.Println("Found Tokens:")
	for _, t := range tokens {
		fmt.Println(" *", t)
	}
}

func withParser() {
	tokchan := lex(input)
	nodechan := parse(tokchan)

	fmt.Println("Running...")
	for n := range nodechan {
		fmt.Println("Received Node: ", n)
	}
}

func withMarshaller() {
	tokchan := lex(input)
	nodechan := parse(tokchan)
	marchan := marshal(nodechan)

	fmt.Println("Running...")
	for n := range marchan {
		fmt.Println("Received Value: ", n)
	}
}

func main() {
	// withLexer()
	// withParser()
	withMarshaller()
}
