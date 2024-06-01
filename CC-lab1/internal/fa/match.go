package fa

func (dfa *DFA) Match(s string) bool {
	curState := dfa.q0

	transitionFromSymbolMap := make(map[int]map[byte]int)
	for _, t := range dfa.d {
		if _, ok := transitionFromSymbolMap[t.From]; !ok {
			transitionFromSymbolMap[t.From] = make(map[byte]int)
		}
		transitionFromSymbolMap[t.From][t.Symbol] = t.To
	}

	for i := range s {
		nextState, ok := transitionFromSymbolMap[curState][s[i]]
		if !ok {
			return false
		}
		curState = nextState
	}
	return dfa.f[curState]
}
