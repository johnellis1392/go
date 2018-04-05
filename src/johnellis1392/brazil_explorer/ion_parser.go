package main

import (
  "fmt"
  "io/ioutil"
  "strings"
)


// Links:
// https://stackoverflow.com/questions/8422146/go-how-to-create-a-parser?utm_medium=organic&utm_source=google_rich_qa&utm_campaign=google_rich_qa
// https://blog.gopheracademy.com/advent-2014/parsers-lexers/
// https://github.com/tcolgate/mp3
// http://www.lihaoyi.com/post/EasyParsingwithParserCombinators.html
const (
  AlphabetLower = "abcdefghijklmnopqrstuvwxyz"
  AlphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
  Alphabet = AlphabetLower + AlphabetUpper

  Numbers = "0123456789"

  Alphanumerics = Alphabet + Numbers
)


type Ion struct {
  Directives []Directive
}

type Directive struct {
  Key string
  Value IonObject
}

type IonObject struct {
  Directives []Directive
}

type IonValue struct {
  Value string
}

func tryParseDirective(data string) (*Directive, string) {
  for _, c := range data {

    switch c {
    case '':
      break
    default:
      break
    }
  }
}



type TOKEN uint8
const (
  LPAREN TOKEN = '('
  RPAREN TOKEN = ')'

  LSQBRACE TOKEN = '['
  RSQBRACE TOKEN = ']'

  LCRLBRACE TOKEN = '{'
  RCRLBRACE TOKEN = '}'

  EQUALS TOKEN = '='
  SEMICOLON TOKEN = ';'
  COLON TOKEN = ':'
)




type Position struct {
  Row uint32
  Column uint32
}

type Token struct {
  ID TOKEN
  Position Position
}

func isAlphanumeric(c uint8) bool {
  return strings.ContainsAny(Alphanumerics, string(c))
}

func tokenize(data string) ([]Token, error) {
  return nil, nil
}

// Parse parses text into Ion data.
func Parse(data string) (*Ion, error) {
  tokens, err := tokenize(data)
  if err != nil {
    return nil, err
  }

  return nil, nil
}

// ParseFile parses the given file into an Ion object.
func ParseFile(filename string) (*Ion, error) {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }

  result, err := Parse(data)
  if err != nil {
    return nil, err
  }

  return result, nil
}
