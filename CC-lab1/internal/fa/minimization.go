package fa

import (
	"github.com/Astemirdum/CC-lab1/internal/containers"
)

func (dfa *DFA) minimize() {
	p := []containers.Set[int]{dfa.f, dfa.q.Subtract(dfa.f)}
	stateClass := getStateClass(p)
	invMap := mapTransitions(dfa.d) // inverse transitions map

	charQueue := containers.NewQueue()
	setQueue := containers.NewQueue()

	for c := range dfa.t {
		charQueue.Push(c)
		charQueue.Push(c)
		setQueue.Push(dfa.f)
		setQueue.Push(dfa.q.Subtract(dfa.f))
	}

	for charQueue.Len() > 0 {
		c := charQueue.Pop().(byte)
		s := setQueue.Pop().(containers.Set[int])

		involved := make(map[int]containers.Set[int])
		for q := range s {
			if fromSet, ok := invMap[c][q]; ok {
				for r := range fromSet {
					class := stateClass[r]
					if _, ok := involved[class]; !ok {
						involved[class] = make(containers.Set[int])
					}
					involved[class].Add(r)
				}
			}
		}

		for class := range involved {
			if involved[class].Size() < len(p[class]) {
				p = append(p, containers.NewSet[int]())
				j := len(p) - 1

				for r := range involved[class] {
					p = swapState(p, r, class, j)
				}

				if len(p[j]) > len(p[class]) {
					p[j], p[class] = p[class], p[j]
				}

				for r := range p[j] {
					stateClass[r] = j
				}

				for c := range dfa.t {
					charQueue.Push(c)
					setQueue.Push(p[j])
				}
			}
		}
	}

	dfa.mapOptimalStates(p)
}

func (dfa *DFA) mapOptimalStates(P []containers.Set[int]) {
	stateMap := make(map[int]int)

	// map states
	newQ := make(containers.Set[int], len(P))
	for i, q := range P {
		for state := range q {
			stateMap[state] = i
		}
		newQ[i] = true
	}

	// map fininal states
	resultF := containers.NewSet[int]()
	for f := range dfa.f {
		for k, v := range stateMap {
			if k == f {
				resultF.Add(v)
			}
		}
	}

	// map transitions
	resultTrans := make([]DFATransition, 0)
	for _, t := range dfa.d {
		founds := 0
		tran := DFATransition{Symbol: t.Symbol}
		for k, v := range stateMap {
			if k == t.From {
				tran.From = v
				founds++
			}
			if k == t.To {
				tran.To = v
				founds++
			}
			if founds == 2 {
				break
			}
		}
		resultTrans = append(resultTrans, tran)
	}

	dfa.q0 = stateMap[dfa.q0]
	dfa.q = newQ
	dfa.f = resultF
	dfa.d = resultTrans
}

func swapState(p []containers.Set[int], state int, from int, to int) []containers.Set[int] {
	for v := range p[from] {
		if v == state {
			p[from].Remove(v)
			break
		}
	}

	p[to].Add(state)
	return p
}

func getStateClass(p []containers.Set[int]) map[int]int {
	stateClass := make(map[int]int)

	for classIndex := range p {
		for v := range p[classIndex] {
			stateClass[v] = classIndex
		}
	}
	return stateClass
}

func mapTransitions(trans []DFATransition) map[byte]map[int]containers.Set[int] {
	inv := make(map[byte]map[int]containers.Set[int])

	for _, t := range trans {
		if _, ok := inv[t.Symbol]; !ok {
			inv[t.Symbol] = make(map[int]containers.Set[int])
		}
		if _, ok := inv[t.Symbol][t.To]; !ok {
			inv[t.Symbol][t.To] = make(containers.Set[int])
		}
		inv[t.Symbol][t.To].Add(t.From)
	}

	return inv
}
