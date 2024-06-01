package fa

import (
	"github.com/Astemirdum/CC-lab1/internal/containers"
	"github.com/Astemirdum/CC-lab1/internal/tree"
)

type NFA struct {
	q  []*state              // конечное состояний
	t  containers.Set[byte]  // входной алфавит
	s  []NFATransition       // функция переходов
	q0 containers.Set[int]   // начальное состояние
	f  []containers.Set[int] // конечных состояний
}

type NFATransition struct {
	state     containers.Set[int]
	symbol    byte
	destState containers.Set[int]
}

type state struct {
	value  containers.Set[int]
	marked bool
}

func (nfa *NFA) hasUnmarked() bool {
	for _, s := range nfa.q {
		if !s.marked {
			return true
		}
	}

	return false
}

func (nfa *NFA) getUnmarkedPos() int {
	for i, s := range nfa.q {
		if !s.marked {
			return i
		}
	}

	return -1
}

func (nfa *NFA) getAlphabet(re string) {
	for i := range re {
		if tree.IsSymbol(re[i]) {
			nfa.t.Add(re[i])
		}
	}
}

func (nfa *NFA) hasState(s containers.Set[int]) bool {
	for _, st := range nfa.q {
		if st.value.Equals(s) {
			return true
		}
	}

	return false
}
