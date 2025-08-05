package main

import (
	"fmt"
)

type Lexer struct {
	source  string
	curPos  int
	curChar string
}

type TokenType int

const (
	EOF     TokenType = iota // 0
	NEWLINE                  // 1
	NUMBER                   // 2
	IDENT                    // 3
	STRING                   // 4
	// Keywords.
	LABEL    // 5
	GOTO     // 6
	PRINT    // 7
	INPUT    // 8
	LET      // 9
	IF       // 10
	THEN     // 11
	ENDIF    // 12
	WHILE    // 13
	REPEAT   // 14
	ENDWHILE // 15
	// Operators.
	EQ       // 16
	PLUS     // 17
	MINUS    // 18
	ASTERISK // 19
	SLASH    // 20
	EQEQ     // 21
	NOTEQ    // 22
	LT       // 23
	LTEQ     // 24
	GT       // 25
	GTEQ     // 26
)

type Token struct {
	text string
	kind TokenType
}

// Process the next character.
func nextChar(lexer *Lexer) string {
	lexer.curPos++
	if lexer.curPos >= len(lexer.source) {
		lexer.curChar = "\x00" // EOF
		return lexer.curChar
	}
	lexer.curChar = string(lexer.source[lexer.curPos])
	return lexer.curChar
}

// Return the lookahead character.
func peek(lexer *Lexer) string {
	if lexer.curPos+1 >= len(lexer.source) {
		return "\x00" // EOF
	}
	return string(lexer.source[lexer.curPos+1])
}

// Invalid token found, print error message and exit.
func abort(msg string) error {
	return fmt.Errorf("Lexing error. %s", msg)
}

// Skip whitespace except newlines, which we will use to indicate the end of a statement.
func skipWhitespace(lexer *Lexer) {
	for lexer.curChar == " " || lexer.curChar == "\t" || lexer.curChar == "\r" {
		nextChar(lexer)
	}
}

// Skip comments in the code.
func skipComment(lexer *Lexer) {
	if lexer.curChar == "#" {
		for lexer.curChar != "\n" {
			nextChar(lexer)
		}
	}
}

// Return the next token.
func getToken(lexer *Lexer) (Token, error) {
	skipWhitespace(lexer)
	skipComment(lexer)

	token := Token{text: lexer.curChar}
	switch lexer.curChar {
	case "+":
		token.kind = PLUS
	case "-":
		token.kind = MINUS
	case "*":
		token.kind = ASTERISK
	case "/":
		token.kind = SLASH
	case "=":
		if peek(lexer) == "=" { // Use peek to check for double characters
			nextChar(lexer)
			token.kind = EQEQ
			token.text += "="
		} else {
			token.kind = EQ
		}
	case ">":
		if peek(lexer) == "=" {
			nextChar(lexer)
			token.kind = GTEQ
			token.text += "="
		} else {
			token.kind = GT
		}
	case "<":
		if peek(lexer) == "=" {
			nextChar(lexer)
			token.kind = LTEQ
			token.text += "="
		} else {
			token.kind = LT
		}
	case "!":
		if peek(lexer) == "=" {
			nextChar(lexer)
			token.kind = NOTEQ
			token.text += "="
		} else {
			return token, abort("Expected !=, got !" + peek(lexer))
		}
	case "\"": // Get characters between quotations
		nextChar(lexer)
		startPos := lexer.curPos
		for lexer.curChar != "\"" {
			// Don't allow special characters in the string. No escape characters, newlines, tabs, or %.
			// We will be using C's printf on this string.
			if lexer.curChar == "\r" || lexer.curChar == "\t" || lexer.curChar == "\\" || lexer.curChar == "%" {
				return token, abort("Illegal character in string.")
			}
			nextChar(lexer)
		}
		token.kind = STRING
		token.text = lexer.source[startPos:lexer.curPos] // Get the substring.
	case "\n":
		token.kind = NEWLINE
	case "\x00":
		token.kind = EOF
	default:
		return token, abort("Unknown token: " + lexer.curChar)
	}

	nextChar(lexer)
	return token, nil
}
