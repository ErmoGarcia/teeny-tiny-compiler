package main

import "fmt"

type Parser struct {
	lexer               *Lexer
	curToken, peekToken Token
}

// Return true if the current token matches.
func checkToken(parser *Parser, kind TokenType) bool {
	return kind == parser.curToken.kind
}

// Return true if the next token matches.
func checkPeek(parser *Parser, kind TokenType) bool {
	return kind == parser.peekToken.kind
}

// Try to match current token. If not, error. Advances the current token.
func match(parser *Parser, kind TokenType) {
	if !checkToken(parser, kind) {
		abort(fmt.Sprintf("Expected %d, got %d", kind, parser.curToken.kind))
	}
	nextToken(parser)
}

// Advances the current token.
func nextToken(parser *Parser) error {
	parser.curToken = parser.peekToken
	token, err := getToken(parser.lexer)
	parser.peekToken = token
	// No need to worry about passing the EOF, lexer handles that.
	return err
}
