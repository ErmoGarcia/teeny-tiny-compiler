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
	for scanner.Scan() {
		lexer := Lexer{scanner.Text() + "\n", -1, ""}
		nextChar(&lexer)
		// for peek(&lexer) != "\x00" {
		for lexer.curPos < len(lexer.source) {
			// curChar := nextChar(&lexer)
			// fmt.Printf("%s\n", curChar)
			token, err := getToken(&lexer)
			check(err)
			fmt.Printf("%s: %d\n", token.text, token.kind)
		}
	}
	check(scanner.Err())
}
