package algo

import (
	"reflect"
	"sort"
)

// different states that a nonterminal can be in during this process.
const (
	UNDEFINED = 0
	PENDING   = 1
	DEFINED   = 2
)

type State struct {
	Status int
	Flag   bool
}

// The removeDuplicates function is a helper function to remove duplicate Rule pointers from a slice of Rule
// pointers based on their underlying values to ensure no redundant rules are added to the grammar.
func removeDuplicates(slice []*Rule) []*Rule {
	allKeys := make(map[*Rule]Rule, 0)
	list := make([]*Rule, 0)
	for _, item := range slice {
		if value, ok := allKeys[item]; !ok && !reflect.DeepEqual(*item, value) {
			allKeys[item] = *item
			list = append(list, item)
		}
	}
	return list
}

func DeleteEpsRules(g *Grammar) *Grammar {
	// It starts by identifying nonterminals that can produce ε and then eliminates ε-rules by creating new production rules wherever necessary.

	//In the DeleteEpsRules function, once empty (ε-producing) nonterminals are identified, the function creates all possible
	//combinations of the symbols on the right-hand side of a rule, excluding combinations that produce ε when not allowed.
	//New production rules are created considering these combinations.
	//The resulting grammar will have all original ε-rules removed, preserving language equivalence.
	epsPart := Part{Type: T, Value: EPS}
	emptyNonTerms := getEmptyNonterminals(g.NonTerminals, g.Rules)
	emptyRules := make([]int, 0)
	newRules := make([]*Rule, 0)
	for i, rule := range g.Rules {
		if rule.Right[0] == epsPart {
			emptyRules = append(emptyRules, i)
		} else {
			term := false
			index := make([]int, 0)
			for j, part := range rule.Right {
				search, limit := sort.SearchStrings(emptyNonTerms, part.Value), len(emptyNonTerms)
				if !term && search == limit {
					term = true
				} else if search < limit {
					index = append(index, j)
				}
			}
			combs := all(index)
			if !term {
				combs = combs[:len(combs)-1]
			}
			for _, comb := range combs {
				r := &Rule{Left: rule.Left, Right: []Part{}}
				start := 0
				for _, j := range comb {
					r.Right = append(r.Right, rule.Right[start:j]...)
					start = j + 1
				}
				r.Right = append(r.Right, rule.Right[start:]...)
				newRules = append(newRules, r)
			}
		}
	}

	for i := range emptyRules {
		g.Rules = append(g.Rules[:emptyRules[i]-i], g.Rules[emptyRules[i]-i+1:]...)
	}
	g.Rules = append(g.Rules, newRules...)
	g.Rules = removeDuplicates(g.Rules)
	if sort.SearchStrings(emptyNonTerms, g.Start) < len(emptyNonTerms) {
		start := g.Start + NEW
		g.NonTerminals = append(g.NonTerminals, start)
		g.Rules = append(g.Rules, &Rule{
			Left: Part{Type: N, Value: start},
			Right: []Part{{
				Type:  N,
				Value: g.Start,
			}},
		}, &Rule{
			Left:  Part{Type: N, Value: start},
			Right: []Part{epsPart},
		})
		g.Start = start
	}
	return g
}

// getEmptyNonterminals function takes a slice of nonterminals and a set of production rules and returns a sorted list of
// nonterminals that can produce ε directly or indirectly. It initializes a map of State pointers for each nonterminal and
// uses the searchEmpty function to recursively mark nonterminals that can produce ε.
func getEmptyNonterminals(nonterms []string, rules []*Rule) []string {
	states := make(map[string]*State, len(nonterms))
	for _, nt := range nonterms {
		states[nt] = new(State)
	}
	for _, nt := range nonterms {
		states[nt].Status = PENDING
		searchEmpty(states, nt, rules)
		states[nt].Status = DEFINED
	}

	result := make([]string, 0)
	for nt, state := range states {
		if state.Flag {
			result = append(result, nt)
		}
	}
	sort.Strings(result)
	return result
}

// The searchEmpty function takes a map of nonterminal State pointers, a specific nonterminal nt, and a set of production rules.
// It modifies the State pointers to indicate whether a given nonterminal can produce ε. It does so by checking the production
// rules and recursively applying the same logic
// to any nonterminals on the right-hand side of production rules that only consist of nonterminals.
func searchEmpty(states map[string]*State, nt string, rules []*Rule) {
	if states[nt].Status == DEFINED {
		return
	}

	index := make([]int, 0)
	part, epsPart := Part{Type: N, Value: nt}, Part{Type: T, Value: EPS}
	for i, rule := range rules {
		if rule.Left == part {
			if rule.Right[0] == epsPart {
				states[nt].Flag = true
				return
			}
			index = append(index, i)
		}
	}
	nonterms := make([]string, 0)
	for _, i := range index {
		for _, right := range rules[i].Right {
			if right.Type == N && states[right.Value].Status == UNDEFINED {
				nonterms = append(nonterms, right.Value)
			}
		}
	}

	for _, v := range nonterms {
		states[v].Status = PENDING
		searchEmpty(states, v, rules)
		states[v].Status = DEFINED
		if states[v].Flag {
			states[nt].Flag = true
		}
	}
}

// All returns all possible combinations of the input slice of ints.
func All(input []int) [][]int {
	// Helper function to create combinations.
	var helper func([]int, int) [][]int
	helper = func(comb []int, n int) [][]int {
		if n == len(input) {
			// Make a copy of the current combination to avoid issues with references.
			combCopy := make([]int, len(comb))
			copy(combCopy, comb)
			return [][]int{combCopy}
		}

		// Do not include the current element at 'n'.
		combsWithOutN := helper(comb, n+1)

		// Include the current element at 'n'.
		combWithN := make([]int, len(comb)+1)
		copy(combWithN, comb)
		combWithN[len(comb)] = input[n]
		combsWithN := helper(combWithN, n+1)

		// Return combination of both scenarios.
		return append(combsWithOutN, combsWithN...)
	}

	// Start the helper function with an empty combination and starting index 0.
	return helper([]int{}, 0)
}

func all(set []int) (subsets [][]int) {
	length := uint(len(set))

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		var subset []int

		for object := uint(0); object < length; object++ {
			// checks if object is contained in subset
			// by checking if bit 'object' is set in subsetBits
			if (subsetBits>>object)&1 == 1 {
				// add object to subset
				subset = append(subset, set[object])
			}
		}
		// add subset to subsets
		subsets = append(subsets, subset)
	}
	return subsets
}
