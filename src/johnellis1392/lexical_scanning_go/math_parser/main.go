package main

import (
	"fmt"
)

func withLexer() {
	const input = `1 + 0 * 1`
	tokchan := lex(input)

	var toks []token
	for t := range tokchan {
		toks = append(toks, t)
	}

	fmt.Println("Received Tokens:")
	for _, t := range toks {
		fmt.Println(t)
	}
}

func main() {
	withLexer()
}
