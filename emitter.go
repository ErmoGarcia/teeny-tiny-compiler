package main

import "os"

type Emitter struct {
	fullPath, header, code string
}

func emit(emitter *Emitter, code string) {
	emitter.code += code
}

func emitLine(emitter *Emitter, code string) {
	emitter.code += code + "\n"
}

func headerLine(emitter *Emitter, code string) {
	emitter.header += code + "\n"
}

func writeFile(emitter *Emitter) {
	// f, err := os.Create(emitter.fullPath)
	// check(err)
	// defer f.Close()
	// _, err = f.WriteString(emitter.header + emitter.code)
	// check(err)
	// err = f.Sync()
	// check(err)
	d1 := []byte(emitter.header + emitter.code)
	err := os.WriteFile(emitter.fullPath, d1, 0644)
	check(err)
}
