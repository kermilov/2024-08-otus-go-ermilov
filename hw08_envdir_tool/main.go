package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go-envdir <dir> <command> [args...]")
	}

	dir := os.Args[1]
	cmd := os.Args[2:]
	env, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := RunCmd(cmd, env)

	os.Exit(exitCode)
}
