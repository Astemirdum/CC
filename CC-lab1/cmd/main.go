package main

import (
	"flag"
	"fmt"

	"github.com/Astemirdum/CC-lab1/internal/fa"
)

func main() {
	regex := flag.String("r", "", "The regular expression to be used")
	str := flag.String("s", "", "The string to be tested against the regular expression")

	flag.Parse()

	if *regex == "" {
		fmt.Println("Please provide a regular expression: -r \"\" ")
		return
	}

	if *regex == "" || *str == "" {
		fmt.Println("Please provide a string to test against it: -s \"\" ")
		return
	}

	matched := runDFA(*regex, *str)
	fmt.Println("Matched:", matched)
}

func runDFA(regex, str string) bool {
	return fa.Build(regex).Match(str)
}
