# Package Parser EBNF

// Non-Terminals
file := decl* eof
decl := ident '=' value ';'
value := (ident | string | number | path | object)
object := '{' (decl)* '}'

// Terminals
ident := letter alpha*
string := '"' ('\\' . | [^eof"\n]) '"'
number := leadingDigit digit* (\. digit+)?
path := filename ('/' filename)*
filename := (alpha | '.')+

letter := [a-zA-Z]
digit := [0-9]
leadingDigit := [1-9]
underscore := '\_'
alpha := (letter | digit | underscore)


## lexFile
file := . {decl} eof
decl := . ident '=' value ';'

- follow: {eof, ident}

- lexFile(eof) => nil // Done
- lexFile(ident) => pushState(lexFile), lexDecl
- _ => errorf

## lexDecl
decl := ident . '=' value ';'

- follow: {'='}

- lexDecl('=') => lexValue
- _ => errorf

## lexValue
decl := ident '=' . value ';'
value := . object
object := . '{' {decl} '}'
value := . ident
value := . string
value := . number
value := . path

- follow: {ident, string, number, path, '{'}

- lexValue(ident,string,number,path) => lexAfterValue
- lexValue('{') => lexObject
- _ => errorf

## lexAfterValue
value := ident .
value := string .
value := number .
value := path .
value := object .
decl := ident '=' value . ';'

- follow: {';'}

- lexAfterValue(';') => state := popState(), state
- _ => errorf

## lexObject
object := '{' . {decl} '}'
decl := . ident '=' value ';'
object := '{' {decl} . '}' // For Empty Object

- follow: {ident, '}'}

- lexObject(ident) => pushState(lexObject), lexDecl
- lexObject('}') => lexAfterValue

## lexAfterDecl
decl := ident '=' value ';' .
file := {decl} . eof
object := '{' {decl} . '}'

file := decl . {decl} eof
object := '{' decl . {decl} '}'
decl := . ident '=' value ';'

- follow: {eof, '}', ident}

- lexAfterDecl(eof) => state := popState(), state
- lexAfterDecl('}') => state := popState(), state
- lexAfterDecl(ident) => lexDecl
