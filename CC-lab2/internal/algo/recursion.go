package algo

import "fmt"

const (
	N   = "nonterminal"
	T   = "terminal"
	EPS = "eps"
	NEW = "'"
)

type Part struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewN(v string) Part {
	return Part{
		Type:  N,
		Value: v,
	}
}

type Rule struct {
	Left  Part   `json:"left"`
	Right []Part `json:"right"`
}

type Grammar struct {
	NonTerminals []string `json:"nonterminals"`
	Terminals    []string `json:"terminals"`
	Start        string   `json:"start"`
	Rules        []*Rule  `json:"rules"`
}

func DeleteLeftRecursion(g *Grammar) *Grammar {
	for i := range g.NonTerminals {
		for j := 0; j < i; j++ {
			g.DeleteIndirect(g.NonTerminals[i], g.NonTerminals[j])
		}
		g.DeleteImmediate(g.NonTerminals[i])
	}
	g.LeftFactor()
	return g
}

// DeleteIndirect method, which replaces left-recursive rules with equivalent non-left-recursive rules.
func (g *Grammar) DeleteIndirect(ai, aj string) {
	nI, nJ := NewN(ai), NewN(aj)
	var with, without []*Rule
	for _, rule := range g.Rules {
		if rule.Left == nI {
			if rule.Right[0] == nJ {
				with = append(with, rule)
			}
		} else if rule.Left == nJ {
			without = append(without, rule)
		}
	}

	for _, rule := range with {
		replacePart := rule.Right[1:]
		rule.Right = append(without[0].Right, replacePart...)
		for i := 1; i < len(without); i++ {
			g.Rules = append(g.Rules, &Rule{Left: rule.Left, Right: append(without[i].Right, replacePart...)})
		}
	}
}

// DeleteImmediate method, it eliminates immediate left recursion,
// which introduces new nonterminal symbols and modifies rules to remove direct left recursion.
func (g *Grammar) DeleteImmediate(ai string) {
	part := NewN(ai)
	var with, without []*Rule
	for _, rule := range g.Rules {
		if rule.Left == part {
			if rule.Right[0] == part {
				with = append(with, rule)
			} else {
				without = append(without, rule)
			}
		}
	}
	if len(with) == 0 {
		return
	}

	pt := ai + NEW
	g.NonTerminals = append(g.NonTerminals, pt)
	part = Part{Type: N, Value: pt}
	epsPart := Part{Type: T, Value: EPS}
	for _, rule := range without {
		if rule.Right[0] != epsPart {
			rule.Right = append(rule.Right, part)
		} else {
			rule.Right[0] = part
		}
	}
	for _, rule := range with {
		rule.Left = part
		rule.Right = append(rule.Right[1:], part)
	}
	g.Rules = append(g.Rules, &Rule{Left: part, Right: []Part{epsPart}})
}

// LeftFactor method is used to factor out common prefixes in the production rules of the grammar. It groups rules by their left-hand side nonterminal and then calls the factorize function
func (g *Grammar) LeftFactor() {
	groups := make(map[Part][]*Rule)
	for _, rule := range g.Rules {
		groups[rule.Left] = append(groups[rule.Left], rule)
	}
	for nt, group := range groups {
		nonterms, rules := factorize(nt.Value, group)
		g.NonTerminals = append(g.NonTerminals, nonterms...)
		g.Rules = append(g.Rules, rules...)
	}
}

// factorize function, which creates new nonterminal symbols and rules to eliminate the common prefixes.
func factorize(nt string, group []*Rule) ([]string, []*Rule) {
	//The factorize function iterates over the rules in a group and identifies common prefixes between pairs of rules.
	//When a common prefix is found, it creates a new nonterminal symbol and modifies the original rules to use this new symbol,
	//with the remaining parts of the original rules as the right-hand side. This process continues until no more common prefixes can be found.
	counter := 0
	epsPart := Part{Type: T, Value: EPS}
	nonterms := make([]string, 0)
	rules := make([]*Rule, 0)
	for {
		n, flag := len(group), false
		var i, j, k int
		for i = 0; i < n-1 && !flag; i++ {
			for j = i + 1; j < n && !flag; j++ {
				if group[i].Right[0] != group[j].Right[0] || group[i].Right[0] == epsPart {
					continue
				}
				flag = true
				limit := min(len(group[i].Right), len(group[j].Right))
				for k = 1; k < limit; k++ {
					if group[i].Right[k] != group[j].Right[k] {
						break
					}
				}
			}
		}
		i--
		j--
		if !flag {
			return nonterms, rules
		}

		ntNew := fmt.Sprintf("%s%d", nt, counter)
		counter++
		nonterms = append(nonterms, ntNew)
		part := Part{Type: N, Value: ntNew}
		rule := &Rule{
			Left:  Part{Type: N, Value: nt},
			Right: append(group[i].Right[:k], part),
		}
		group = append(group, rule)
		rules = append(rules, rule)
		if len(group[i].Right[k:]) > 0 {
			group[i].Left, group[i].Right = part, group[i].Right[k:]
		} else {
			group[i].Left, group[i].Right = part, []Part{epsPart}
		}
		if len(group[j].Right[k:]) > 0 {
			group[j].Left, group[j].Right = part, group[j].Right[k:]
		} else {
			group[j].Left, group[j].Right = part, []Part{epsPart}
		}
	}
}
