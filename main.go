package main

import (
	"fmt"
	"my-interpreter/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s, This is my own Interpreter \n", user.Username)
	fmt.Printf("Feel free to type in commands \n")
	repl.Start(os.Stdin, os.Stdout)
}
