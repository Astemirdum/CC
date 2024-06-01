package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Astemirdum/CC-lab2/internal/algo"
	"os"
)

func main() {
	var epsFlag bool
	var inp, out string

	flag.StringVar(&inp, "f", "", "input")
	flag.StringVar(&out, "o", "", "output")
	flag.BoolVar(&epsFlag, "e", false, "is eps-rules deleting algorithm")

	flag.Parse()

	if inp == "" {
		fmt.Println("need -f -e")
		os.Exit(1)
	}

	g := new(algo.Grammar)

	f, err := os.Open(inp)
	if err != nil {
		fmt.Println("Unable to read the input file.")
		os.Exit(1)
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(g); err != nil {
		fmt.Println("Invalid json.")
		os.Exit(1)
	}

	var result *algo.Grammar
	if epsFlag {
		result = algo.DeleteEpsRules(g)
	} else {
		result = algo.DeleteLeftRecursion(g)
	}
	output := os.Stdout
	if out != "" {
		output, err = os.Create(out)
		if err != nil {
			fmt.Println("out create", err.Error())
			os.Exit(-1)
		}
		defer output.Close()
	}
	if result != nil {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Fprint(output, string(data))
	} else {
		fmt.Println("Error occured while grammar transforming.")
		os.Exit(1)
	}
}
