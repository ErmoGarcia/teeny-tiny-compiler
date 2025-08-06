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
		panic(fmt.Sprintf("Parsing error. Expected %d, got %d.", kind, parser.curToken.kind))
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

// Production rules.

// program ::= {statement}
func program(parser *Parser) {
	// fmt.Println("PROGRAM")

	// Since some newlines are required in our grammar, need to skip the excess.
	for checkToken(parser, NEWLINE) {
		nextToken(parser)
	}

	// Parse all the statements in the program.
	for !checkToken(parser, EOF) {
		statement(parser)
	}
}

// One of the following statements...
func statement(parser *Parser) {
	// Check the first token to see what kind of statement this is.

	if checkToken(parser, PRINT) { // "PRINT" (expression | string)
		fmt.Println("STATEMENT-PRINT")
		nextToken(parser)

		if checkToken(parser, STRING) {
			// Simple string.
			nextToken(parser)
		} else {
			// Expect an expression.
			expression(parser)
		}
	} else if checkToken(parser, IF) { // "IF" comparison "THEN" {statement} "ENDIF"
		fmt.Println("STATEMENT-IF")
		nextToken(parser)
		comparison(parser)

		match(parser, THEN)
		nl(parser)

		for !checkToken(parser, ENDIF) {
			statement(parser)
		}
		match(parser, ENDIF)
	} else if checkToken(parser, WHILE) { // "WHILE" comparison "REPEAT" {statement} "ENDWHILE"
		fmt.Println("STATEMENT-WHILE")
		nextToken(parser)
		comparison(parser)

		match(parser, REPEAT)
		nl(parser)

		for !checkToken(parser, ENDWHILE) {
			statement(parser)
		}
		match(parser, ENDWHILE)
	} else if checkToken(parser, LABEL) { // "LABEL" ident
		fmt.Println("STATEMENT-LABEL")
		nextToken(parser)
		match(parser, IDENT)
	} else if checkToken(parser, GOTO) { // "GOTO" ident
		fmt.Println("STATEMENT-GOTO")
		nextToken(parser)
		match(parser, IDENT)
	} else if checkToken(parser, LET) { // "LET" ident "=" expresion
		fmt.Println("STATEMENT-LET")
		nextToken(parser)
		match(parser, IDENT)
		match(parser, EQ)
		expression(parser)
	} else if checkToken(parser, INPUT) { // "INPIT" ident
		fmt.Println("STATEMENT-INPUT")
		nextToken(parser)
		match(parser, IDENT)
	} else { //This is not a valid statement. Error!
		panic(fmt.Sprintf("Parsing error. Invalid statement at %s (%d)", parser.curToken.text, parser.curToken.kind))
	}

	// Newline.
	nl(parser)
}

func expression(parser *Parser) {

}

func comparison(parser *Parser) {

}

// nl ::= '\n'+
func nl(parser *Parser) {
	fmt.Println("NEWLINE")

	// Require at least one newline.
	match(parser, NEWLINE)
	// But we will allow extra newlines too, of course.
	for checkToken(parser, NEWLINE) {
		nextToken(parser)
	}
}
