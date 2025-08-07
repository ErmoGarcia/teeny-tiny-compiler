package main

import (
	"fmt"
)

type Parser struct {
	lexer               *Lexer
	emitter             *Emitter
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
	headerLine(parser.emitter, "#include <stdio.h>")
	headerLine(parser.emitter, "int main(void){")

	// Since some newlines are required in our grammar, need to skip the excess.
	for checkToken(parser, NEWLINE) {
		nextToken(parser)
	}

	// Parse all the statements in the program.
	for !checkToken(parser, EOF) {
		statement(parser)
	}

	// Wrap things up.
	emitLine(parser.emitter, "return 0;")
	emitLine(parser.emitter, "}")

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
		nextToken(parser)

		if checkToken(parser, STRING) {
			// Simple string.
			emitLine(parser.emitter, "printf(\""+parser.curToken.text+"\\n\");")
			nextToken(parser)
		} else {
			// Expect an expression.
			emit(parser.emitter, "printf(\"%"+".2f\\n\", (float)(")
			expression(parser)
			emitLine(parser.emitter, "));")
		}

	} else if checkToken(parser, IF) { // "IF" comparison "THEN" {statement} "ENDIF"
		nextToken(parser)
		emit(parser.emitter, "if(")
		comparison(parser)

		match(parser, THEN)
		nl(parser)
		emitLine(parser.emitter, "){")

		for !checkToken(parser, ENDIF) {
			statement(parser)
		}

		match(parser, ENDIF)
		emitLine(parser.emitter, "}")

	} else if checkToken(parser, WHILE) { // "WHILE" comparison "REPEAT" {statement} "ENDWHILE"
		nextToken(parser)
		emit(parser.emitter, "while (")
		comparison(parser)

		match(parser, REPEAT)
		nl(parser)
		emitLine(parser.emitter, "){")

		for !checkToken(parser, ENDWHILE) {
			statement(parser)
		}

		match(parser, ENDWHILE)
		emitLine(parser.emitter, "}")

	} else if checkToken(parser, LABEL) { // "LABEL" ident
		nextToken(parser)
		match(parser, IDENT)

		// Make sure this label doesn't already exist.
		_, ok := parser.labelsDeclared[parser.curToken.text]
		if ok {
			panic(fmt.Sprintf("Label already exists: %s", parser.curToken.text))
		}
		parser.labelsDeclared[parser.curToken.text] = 1

		emitLine(parser.emitter, parser.curToken.text+":")
		match(parser, IDENT)

	} else if checkToken(parser, GOTO) { // "GOTO" ident
		nextToken(parser)
		match(parser, IDENT)

		parser.labelsGotoed[parser.curToken.text] = 1
		emitLine(parser.emitter, "goto "+parser.curToken.text)
		match(parser, IDENT)

	} else if checkToken(parser, LET) { // "LET" ident "=" expresion
		nextToken(parser)

		// Check if ident exists in symbol table. If not, declare it.
		_, ok := parser.symbols[parser.curToken.text]
		if !ok {
			parser.symbols[parser.curToken.text] = 1
			headerLine(parser.emitter, "float "+parser.curToken.text+";")
		}

		emit(parser.emitter, parser.curToken.text+" = ")
		match(parser, IDENT)
		match(parser, EQ)
		expression(parser)
		emitLine(parser.emitter, ";")

	} else if checkToken(parser, INPUT) { // "INPIT" ident
		nextToken(parser)

		// If variable doesn't already exist, declare it.
		_, ok := parser.symbols[parser.curToken.text]
		if !ok {
			parser.symbols[parser.curToken.text] = 1
			headerLine(parser.emitter, "float "+parser.curToken.text+";")
		}

		// Emit scanf but also validate the input. If invalid, set the variable to 0 and clear the input.
		emitLine(parser.emitter, "if(0 == scanf(\"%"+"f\", &"+parser.curToken.text+")) {")
		emitLine(parser.emitter, parser.curToken.text+" = 0;")
		emit(parser.emitter, "scanf(\"%")
		emitLine(parser.emitter, "*s\");")
		emitLine(parser.emitter, "}")
		match(parser, IDENT)

	} else { //This is not a valid statement. Error!
		panic(fmt.Sprintf("Invalid statement at %s (%d)", parser.curToken.text, parser.curToken.kind))
	}

	// Newline.
	nl(parser)
}

// comparison ::= expression (("==" | "!=" | ">" | ">=" | "<" | "<=") expression)+
func comparison(parser *Parser) {
	expression(parser)

	// Must be at least one comparison operator and another expression.
	if isComparisonOperator(parser.curToken.kind) {
		emit(parser.emitter, parser.curToken.text)
		nextToken(parser)
		expression(parser)
	} else {
		panic(fmt.Sprintf("Expected comparison operator at: %s", parser.curToken.text))
	}

	// Can have 0 or more comparison operator and expressions.
	for isComparisonOperator(parser.curToken.kind) {
		emit(parser.emitter, parser.curToken.text)
		nextToken(parser)
		expression(parser)
	}
}

// expression ::= term {( "-" | "+" ) term}
func expression(parser *Parser) {
	term(parser)
	// Can have 0 or more +/- and expressions.
	for checkToken(parser, PLUS) || checkToken(parser, MINUS) {
		emit(parser.emitter, parser.curToken.text)
		nextToken(parser)
		term(parser)
	}
}

// term ::= unary {( "/" | "*" ) unary}
func term(parser *Parser) {
	unary(parser)
	// Can have 0 or more *// and expressions.
	for checkToken(parser, ASTERISK) || checkToken(parser, SLASH) {
		emit(parser.emitter, parser.curToken.text)
		nextToken(parser)
		unary(parser)
	}
}

// unary ::= ["+" | "-"] primary
func unary(parser *Parser) {
	// Optional unary +/-
	if checkToken(parser, PLUS) || checkToken(parser, MINUS) {
		emit(parser.emitter, parser.curToken.text)
		nextToken(parser)
	}
	primary(parser)
}

// primary ::= number | ident
func primary(parser *Parser) {
	if checkToken(parser, NUMBER) {
		emit(parser.emitter, parser.curToken.text)
		nextToken(parser)
	} else if checkToken(parser, IDENT) {
		// Ensure the variable already exists.
		_, ok := parser.symbols[parser.curToken.text]
		if !ok {
			panic(fmt.Sprintf("Referencing variable before assignment: %s", parser.curToken.text))
		}
		emit(parser.emitter, parser.curToken.text)
		nextToken(parser)
	} else {
		// Error!
		panic(fmt.Sprintf("Unexpected token at %s", parser.curToken.text))
	}
}

// nl ::= '\n'+
func nl(parser *Parser) {
	// Require at least one newline.
	match(parser, NEWLINE)
	// But we will allow extra newlines too, of course.
	for checkToken(parser, NEWLINE) {
		nextToken(parser)
	}
}

func isComparisonOperator(kind TokenType) bool {
	switch kind {
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
