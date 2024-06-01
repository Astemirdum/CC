package tree

import (
	"fmt"
	"github.com/Astemirdum/CC-lab1/internal/containers"
)

type Operation byte

const (
	Concat       Operation = '.'
	Union        Operation = '|'
	KleeneStar   Operation = '*'
	OpenBracket  Operation = '('
	CloseBracket Operation = ')'
)

type NodeTree struct {
	Value    byte
	Pos      int
	Nullable bool
	FirstPos containers.Set[int]
	LastPos  containers.Set[int]
	Children []*NodeTree
}

func (t *NodeTree) String() string {
	return fmt.Sprintf("[v=%s p=%d null=%t, (f=%v, l=%v)]", string(t.Value), t.Pos, t.Nullable, t.FirstPos, t.LastPos)
}

// Shunting Yard algorithm into a syntax tree
func RegexToTree(expr string) *NodeTree {
	curPos := 0

	byteSt := containers.NewStack()
	nodeSt := containers.NewStack()

	p := new(NodeTree)
	for i := range expr {
		switch el := expr[i]; {
		case Operation(el) == CloseBracket:
			for Operation(byteSt.Top().(byte)) != OpenBracket {
				processOperation(byteSt, nodeSt, curPos)
				curPos++
			}
			byteSt.Pop()
		case Operation(el) == OpenBracket:
			byteSt.Push(el)
		case IsBinaryOperation(el) || IsUnaryOperation(el):
			for byteSt.Len() > 0 && (opPrior(byteSt.Top().(byte)) >= opPrior(el)) {
				processOperation(byteSt, nodeSt, curPos)
				curPos++
			}
			byteSt.Push(el)
		default:
			p = &NodeTree{Value: el, Pos: curPos}
			curPos++
			nodeSt.Push(p)
		}
	}

	for byteSt.Len() > 0 {
		processOperation(byteSt, nodeSt, curPos)
	}
	root := nodeSt.Top().(*NodeTree)
	return root
}

func Print(root *NodeTree) {
	var rec func(root *NodeTree) string
	rec = func(root *NodeTree) string {
		if root == nil {
			return ""
		}
		s := fmt.Sprintf("\t%s(p%d)  \t", string(root.Value), root.Pos)
		for _, c := range root.Children {
			s += rec(c)
		}
		return s
	}

	fmt.Println("root", rec(root))
}

func IsBinaryOperation(c byte) bool {
	return Operation(c) == Concat || Operation(c) == Union
}

func IsUnaryOperation(c byte) bool {
	return Operation(c) == KleeneStar
}

func opPrior(op byte) int {
	switch Operation(op) {
	case Union:
		return 1
	case Concat:
		return 2
	case KleeneStar:
		return 3
	default:
		return 0
	}
}

func processOperation(byteSt containers.Stack, nodeSt containers.Stack, curPos int) {
	value := byteSt.Pop().(byte)
	p := &NodeTree{Value: value, Pos: curPos}

	right := nodeSt.Pop().(*NodeTree)
	if IsBinaryOperation(value) {
		left := nodeSt.Pop().(*NodeTree)
		p.Children = append(p.Children, left, right)
	} else {
		p.Children = append(p.Children, right)
	}

	nodeSt.Push(p)
}

func (t *NodeTree) PrepareTree() map[int]containers.Set[int] {
	m := make(map[int]containers.Set[int])
	prepareTreeRecursive(t, m)
	return m
}

func prepareTreeRecursive(root *NodeTree, m map[int]containers.Set[int]) {
	for _, child := range root.Children {
		prepareTreeRecursive(child, m)
	}

	root.Nullable = nullable(root)
	root.FirstPos = firstPos(root)
	root.LastPos = lastPos(root)
	followPos(root, m)
}

func nullable(t *NodeTree) bool {
	switch Operation(t.Value) {
	case Union:
		for _, child := range t.Children {
			if child.Nullable {
				return true
			}
		}
		return false
	case Concat:
		for _, child := range t.Children {
			if !child.Nullable {
				return false
			}
		}
		return true
	case KleeneStar:
		return true
	default:
		return false
	}
}

func firstPos(t *NodeTree) containers.Set[int] {
	s := containers.NewSet[int]()

	switch Operation(t.Value) {
	case Union:
		u := t.Children[0]
		v := t.Children[1]
		s = u.FirstPos.Unite(v.FirstPos)
	case Concat:
		u := t.Children[0]
		v := t.Children[1]
		if u.Nullable {
			s = u.FirstPos.Unite(v.FirstPos)
		} else {
			s = u.FirstPos
		}
	case KleeneStar:
		u := t.Children[0]
		s = u.FirstPos
	default:
		s.Add(t.Pos)
	}

	return s
}

func lastPos(t *NodeTree) containers.Set[int] {
	s := containers.NewSet[int]()

	switch Operation(t.Value) {
	case Union:
		u := t.Children[0]
		v := t.Children[1]
		s = u.LastPos.Unite(v.LastPos)
	case Concat:
		u := t.Children[0]
		v := t.Children[1]
		if v.Nullable {
			s = u.FirstPos.Unite(v.FirstPos)
		} else {
			s = v.FirstPos
		}
	case KleeneStar:
		u := t.Children[0]
		s = u.FirstPos
	default:
		s.Add(t.Pos)
	}

	return s
}

func followPos(t *NodeTree, m map[int]containers.Set[int]) {
	switch Operation(t.Value) {
	case Concat:
		for i := range t.Children[0].LastPos.ToSet() {
			curSet := m[i]
			m[i] = curSet.Unite(t.Children[1].FirstPos)
		}
	case KleeneStar:
		for i := range t.LastPos.ToSet() {
			curSet := m[i]
			m[i] = curSet.Unite(t.FirstPos)
		}
	}
}

func IsSymbol(s byte) bool {
	return !IsUnaryOperation(s) && !IsBinaryOperation(s) && Operation(s) != OpenBracket && Operation(s) != CloseBracket && s != '#'
}

func (t *NodeTree) ToMap() map[int]*NodeTree {
	m := make(map[int]*NodeTree)
	m[t.Pos] = t

	for _, c := range t.Children {
		c.MapRec(m)
	}

	return m
}

func (t *NodeTree) MapRec(m map[int]*NodeTree) {
	m[t.Pos] = t

	for _, c := range t.Children {
		c.MapRec(m)
	}
}
