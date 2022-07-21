package main

import (
	"log"
	"os"
)

func main() {
	arguments := os.Args[1:]
	envDirectory := arguments[0]
	commandArgs := arguments[1:]

	if len(arguments) < 2 {
		log.Fatalf("use more commandArgs")
	}
	env, err := ReadDir(envDirectory)
	if err != nil {
		log.Fatalf("Can't to find directory \"%v\". Error : %v", envDirectory, err)
	}
	os.Exit(RunCmd(commandArgs, env))
}
