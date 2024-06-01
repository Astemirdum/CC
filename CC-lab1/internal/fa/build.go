package fa

import (
	"fmt"
	"strings"

	"github.com/Astemirdum/CC-lab1/internal/containers"
	"github.com/Astemirdum/CC-lab1/internal/tree"
)

func Build(re string) *DFA {
	re = rebuild(re)
	//fmt.Println(re)
	root := tree.RegexToTree(re)
	tree.Print(root)
	dfa := newDFA(re, root)

	fmt.Println("DFA:")
	printConfiguration(dfa)

	dfa.minimize()

	fmt.Println("Min DFA:")
	printConfiguration(dfa)

	return dfa
}

func newDFA(re string, root *tree.NodeTree) *DFA {
	var nfa NFA
	nfa.t = containers.NewSet[byte]()

	nfa.getAlphabet(re)

	followPos := root.PrepareTree()
	fmt.Println(followPos)
	treeMap := root.ToMap()
	fmt.Println(treeMap)
	nfa.q0 = root.FirstPos
	nfa.q = append(nfa.q, &state{value: nfa.q0})
	for nfa.hasUnmarked() {
		a := nfa.q[nfa.getUnmarkedPos()]
		a.marked = true

		for symbol := range nfa.t {
			u := containers.NewSet[int]()
			for p := range a.value {
				if treeMap[p].Value == symbol {
					u = u.Unite(followPos[p])
				}
			}

			if !nfa.hasState(u) {
				nfa.q = append(nfa.q, &state{value: u})
			}
			nfa.s = append(nfa.s, NFATransition{state: a.value, symbol: symbol, destState: u})
		}
	}

	endPos := root.Children[1].Pos
	for _, s := range nfa.q {
		if s.value.Contains(endPos) {
			nfa.f = append(nfa.f, s.value)
		}
	}

	return nfa.toDFA()
}

func rebuild(re string) string {
	builder := strings.Builder{}
	builder.WriteByte('(')

	for i := 0; i < len(re)-1; i++ {
		builder.WriteByte(re[i])
		if !tree.IsBinaryOperation(re[i]) && tree.Operation(re[i]) != tree.OpenBracket &&
			!tree.IsBinaryOperation(re[i+1]) && tree.Operation(re[i+1]) != tree.CloseBracket && !tree.IsUnaryOperation(re[i+1]) {
			builder.WriteByte('.')
		}
	}

	builder.WriteByte(re[len(re)-1])
	builder.WriteString(").#")

	return builder.String()
}
