package main

import (
	"log"
	"os"
)

func main() {
	dir := os.Args[1]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatalln("cannot read directory, error: %w", err)
	}

	code := RunCmd(os.Args[2:], env)
	os.Exit(code)
}
