package fa

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func printConfiguration(dfa *DFA) {
	fmt.Println("q0:", dfa.q0)

	finalStates := make([]string, 0, len(dfa.f))
	for f := range dfa.f {
		finalStates = append(finalStates, strconv.Itoa(f))
	}

	fmt.Println("f: {" + strings.Join(finalStates, ",") + "}")

	fmt.Println("transition table:")
	fmt.Print("\t")

	states := make([]int, 0, len(dfa.q))
	for st := range dfa.q {
		states = append(states, st)
	}
	sort.Ints(states)
	line := make([]string, len(states))
	for i, element := range states {
		line[i] = fmt.Sprint(element)
	}
	fmt.Println(strings.Join(line, "\t"))

	transitionMap := make(map[int]map[int]byte)
	for _, t := range dfa.d {
		if _, ok := transitionMap[t.From]; !ok {
			transitionMap[t.From] = make(map[int]byte)
		}
		transitionMap[t.From][t.To] = t.Symbol
	}

	for _, st := range states {
		fmt.Print(st, "\t")
		for toSt := range states {
			if symbol, ok := transitionMap[st][toSt]; ok {
				fmt.Print(string(symbol), "\t")
			} else {
				fmt.Print("\t")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
