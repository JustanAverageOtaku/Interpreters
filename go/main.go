package main

import (
	"GoInterpreter/src/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome to the most Go-ated programming language, %s\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
