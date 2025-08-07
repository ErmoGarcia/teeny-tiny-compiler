package main

import (
	"fmt"
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: compiler <file_path>")
	}

	file_path := os.Args[1]
	file, err := os.ReadFile(file_path)
	check(err)

	lexer := Lexer{
		source:  string(file) + "\n",
		curPos:  -1,
		curChar: "",
	}
	parser := Parser{
		lexer:          &lexer,
		symbols:        make(map[string]int),
		labelsDeclared: make(map[string]int),
		labelsGotoed:   make(map[string]int),
	}
	nextChar(&lexer)
	nextToken(&parser)
	nextToken(&parser) // Call this twice to initialize current and peek.
	program(&parser)   // Start the parser.
	fmt.Println("Parsing completed.")
}
