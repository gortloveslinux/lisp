package main

import "fmt"

const prompt = "> "

func main() {
	for {
		repl()
	}
}

func repl() {
	fmt.Print(prompt)
}
