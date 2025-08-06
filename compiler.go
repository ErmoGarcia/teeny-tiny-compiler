package main

import (
	"bufio"
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
	file, err := os.Open(file_path)
	check(err)
	defer func() {
		err := file.Close()
		check(err)
	}()

	scanner := bufio.NewScanner(file)
	lexer := Lexer{curChar: ""}
	parser := Parser{lexer: &lexer}
	fmt.Println("PROGRAM")
	for scanner.Scan() {
		lexer.source = scanner.Text() + "\n"
		lexer.curPos = -1
		nextChar(&lexer)
		nextToken(&parser)
		nextToken(&parser) // Call this twice to initialize current and peek.
		// for peek(&lexer) != "\x00" {
		// for lexer.curPos < len(lexer.source) {
		// 	// curChar := nextChar(&lexer)
		// 	// fmt.Printf("%s\n", curChar)
		// 	token, err := getToken(&lexer)
		// 	check(err)
		// 	fmt.Printf("%s: %d\n", token.text, token.kind)
		// }

		program(&parser) // Start the parser.
	}
	fmt.Println("Parsing completed.")
	check(scanner.Err())
}
