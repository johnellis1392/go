package main

// Simple Parse Tree Structure for Math Language

// * Grammar:
// root := expr
// expr := value | add | mul
// add := mul '+' add | mul
// mul := value '*' mul | '(' expr ')' | value
// value := ident | number

// * parseRoot:
// root := . expr
// - parseRoot(digit|alphabet) => parseExpr
// - parseRoot('(') => parseExpr

// * parseExpr:
// expr := . value | . add | . mul
// add := . mul '+' add | . mul
// mul := . value '*' mul | . '(' expr ')' | . value
// value := . number | . ident
// - parseExpr(digit|alphabet) => parseValue
// - parseExpr('(') => parseMul

// * parseMul:
// mul := . value '*' mul | . '(' expr ')' | . value
// - parseMul(digit|alphabet) => parseValue
// - parseMul('(') => accept('('), parseExpr

// * parseValue:
// value := number . | ident .
// number := [1-9][0-9]* .
// ident := [a-zA-Z_][a-zA-Z0-9_]* .
// mul := value . '*' mul | value .
// add := mul . '+' add | mul .
// - parseValue('*') => parseInsideMul
// - parseValue('+') => parseInsideAdd

// * parseInsideMul:
// mul := value '*' . mul
// mul := . value '*' mul
// mul := . '(' expr ')'
// mul := . value
// - parseInsideMul(alpha) => parseValue
// - parseInsideMul('(') => parseExpr

// * parseInsideAdd:
// add := mul '+' . add
// add := . mul '+' add
// add := . mul
// mul := . value '*' mul
// mul := . value
// mul := '(' expr ')'
// - parseInsideMul(alpha) => parseValue
// - parseInsideMul('(') => parseExpr
