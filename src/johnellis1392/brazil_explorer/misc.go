// Miscellaneous Example Code
package main

import (
	"strings"
)

// Code taken from Rune discussion here:
// https://stackoverflow.com/questions/19310700/what-is-a-rune
func misc_SwapRune(r rune) rune {
	switch {
	case 'a' <= r && r <= 'z':
		return r - 'a' + 'A'
	case 'A' <= r && r <= 'Z':
		return r - 'A' + 'a'
	default:
		return r
	}
}

// Change case of string, using strings.Map(...) function
// to map over characters in a string.
//
// From: `go doc strings.Map`
// func Map(mapping func(rune) rune, s string) string
func misc_SwapCase(str string) string {
	return strings.Map(misc_SwapRune, str)
}
