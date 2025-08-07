package main

import (
	"fmt"
)

type Parser struct {
	lexer               *Lexer
	curToken, peekToken Token

	symbols        map[string]int // Variables declared so far.
	labelsDeclared map[string]int // Labels declared so far.
	labelsGotoed   map[string]int // Labels goto'ed so far.
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
		panic(fmt.Sprintf("Expected %d, got %d.", kind, parser.curToken.kind))
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
	fmt.Println("PROGRAM")

	// Since some newlines are required in our grammar, need to skip the excess.
	for checkToken(parser, NEWLINE) {
		nextToken(parser)
	}

	// Parse all the statements in the program.
	for !checkToken(parser, EOF) {
		statement(parser)
	}

	// Check that each label referenced in a GOTO is declared.
	for label := range parser.labelsGotoed {
		_, ok := parser.labelsDeclared[label]
		if !ok {
			panic(fmt.Sprintf("Attempting to GOTO to undeclared label: %s", label))
		}
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

		// Make sure this label doesn't already exist.
		_, ok := parser.labelsDeclared[parser.curToken.text]
		if ok {
			panic(fmt.Sprintf("Label already exists: %s", parser.curToken.text))
		}
		parser.labelsDeclared[parser.curToken.text] = 1

		match(parser, IDENT)
	} else if checkToken(parser, GOTO) { // "GOTO" ident
		fmt.Println("STATEMENT-GOTO")
		nextToken(parser)
		match(parser, IDENT)

		parser.labelsGotoed[parser.curToken.text] = 1
		match(parser, IDENT)
	} else if checkToken(parser, LET) { // "LET" ident "=" expresion
		fmt.Println("STATEMENT-LET")
		nextToken(parser)

		// Check if ident exists in symbol table. If not, declare it.
		_, ok := parser.symbols[parser.curToken.text]
		if !ok {
			parser.symbols[parser.curToken.text] = 1
		}

		match(parser, IDENT)
		match(parser, EQ)
		expression(parser)
	} else if checkToken(parser, INPUT) { // "INPIT" ident
		fmt.Println("STATEMENT-INPUT")
		nextToken(parser)

		// If variable doesn't already exist, declare it.
		_, ok := parser.symbols[parser.curToken.text]
		if !ok {
			parser.symbols[parser.curToken.text] = 1
		}

		match(parser, IDENT)
	} else { //This is not a valid statement. Error!
		panic(fmt.Sprintf("Invalid statement at %s (%d)", parser.curToken.text, parser.curToken.kind))
	}

	// Newline.
	nl(parser)
}

// comparison ::= expression (("==" | "!=" | ">" | ">=" | "<" | "<=") expression)+
func comparison(parser *Parser) {
	fmt.Println("COMPARISON")
	expression(parser)

	// Must be at least one comparison operator and another expression.
	if isComparisonOperator(parser) {
		nextToken(parser)
		expression(parser)
	} else {
		panic(fmt.Sprintf("Expected comparison operator at: %s", parser.curToken.text))
	}

	// Can have 0 or more comparison operator and expressions.
	for isComparisonOperator(parser) {
		nextToken(parser)
		expression(parser)
	}
}

// expression ::= term {( "-" | "+" ) term}
func expression(parser *Parser) {
	fmt.Println("EXPRESSION")

	term(parser)
	// Can have 0 or more +/- and expressions.
	for checkToken(parser, PLUS) || checkToken(parser, MINUS) {
		nextToken(parser)
		term(parser)
	}
}

// term ::= unary {( "/" | "*" ) unary}
func term(parser *Parser) {
	fmt.Println("TERM")

	unary(parser)
	// Can have 0 or more *// and expressions.
	for checkToken(parser, ASTERISK) || checkToken(parser, SLASH) {
		nextToken(parser)
		unary(parser)
	}
}

// unary ::= ["+" | "-"] primary
func unary(parser *Parser) {
	fmt.Println("UNARY")

	// Optional unary +/-
	if checkToken(parser, PLUS) || checkToken(parser, MINUS) {
		nextToken(parser)
	}
	primary(parser)
}

// primary ::= number | ident
func primary(parser *Parser) {
	fmt.Println("PRIMARY (" + parser.curToken.text + ")")

	if checkToken(parser, NUMBER) {
		nextToken(parser)
	} else if checkToken(parser, IDENT) {
		// Ensure the variable already exists.
		_, ok := parser.symbols[parser.curToken.text]
		if !ok {
			panic(fmt.Sprintf("Referencing variable before assignment: %s", parser.curToken.text))
		}
		nextToken(parser)
	} else {
		// Error!
		panic(fmt.Sprintf("Unexpected token at %s", parser.curToken.text))
	}
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

func isComparisonOperator(parser *Parser) bool {
	switch parser.curToken.kind {
	case EQEQ:
		return true
	case NOTEQ:
		return true
	case GT:
		return true
	case GTEQ:
		return true
	case LT:
		return true
	case LTEQ:
		return true
	default:
		return false
	}
}
