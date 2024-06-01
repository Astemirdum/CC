package fa

import "github.com/Astemirdum/CC-lab1/internal/containers"

type DFA struct {
	q  containers.Set[int]
	f  containers.Set[int]
	q0 int
	t  containers.Set[byte]
	d  []DFATransition
}

type DFATransition struct {
	Symbol byte
	From   int
	To     int
}

func (nfa *NFA) toDFA() *DFA {
	dfa := &DFA{
		q: containers.NewSet[int](),
		f: containers.NewSet[int](),
		t: nfa.t,
	}

	stateMap := make(map[*state]int)
	// map states
	curState := 0
	for _, q := range nfa.q {
		if _, ok := stateMap[q]; !ok {
			stateMap[q] = curState
			dfa.q.Add(curState)
			curState++
		}
	}
	// map q0
	for k, v := range stateMap {
		if k.value.Equals(nfa.q0) {
			dfa.q0 = v
		}
	}
	// map final states
	for _, f := range nfa.f {
		for k, v := range stateMap {
			if k.value.Equals(f) {
				dfa.f.Add(v)
			}
		}
	}

	// map transitions
	for _, t := range nfa.s {
		founds := 0
		tran := DFATransition{Symbol: t.symbol}
		for k, v := range stateMap {
			if k.value.Equals(t.state) {
				tran.From = v
				founds++
			}
			if k.value.Equals(t.destState) {
				tran.To = v
				founds++
			}
			if founds == 2 {
				break
			}
		}
		dfa.d = append(dfa.d, tran)
	}

	return dfa
}
