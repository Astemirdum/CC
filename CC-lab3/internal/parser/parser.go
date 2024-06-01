package parser

import (
	"errors"
	"fmt"
	"github.com/Astemirdum/CC-lab3/internal/lexer"
	"log/slog"
)

type Parser interface {
	Parse() (string, error)
}

func New(lex lexer.Lexer) Parser {
	return &parser{
		lex: lex,
	}
}

type parser struct {
	lex lexer.Lexer
}

func (p *parser) Parse() (string, error) {
	root, err := p.program()
	if err != nil {
		return "", p.makeErr(err.Error())
	}
	g := newGraph()
	root.ToAst(g, "")
	return g.String(), nil
}

var (
	ErrNextFail   = errors.New("next fail")
	ErrWrongToken = errors.New("wrong token")
)

func (p *parser) makeErr(msg string) error {
	pos := p.lex.CurPos()
	return fmt.Errorf("msg: %s {token: %s, pos: %d}", msg, p.lex.Token(), pos)
}

func (p *parser) log(msg string) {
	slog.Debug("msg",
		slog.String("token", fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())),
		slog.Int("pos", p.lex.CurPos()),
		slog.String("msg", msg),
	)
}

func (p *parser) program() (*Node, error) {
	p.log("program")
	root := NewNode("program")

	node, err := p.block()
	if err != nil {
		return nil, err
	}
	root.AddChild(node)
	return node, nil
}

func (p *parser) block() (*Node, error) {
	p.log("block")
	root := NewNode("block")

	if !p.lex.Next() || p.lex.Token() != "{" {
		return nil, errors.New("expected token \"{\"")
	}
	root.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))

	node, err := p.operators()
	if err != nil {
		return nil, fmt.Errorf("operators() %v", err)
	}
	root.AddChild(node)

	if !p.lex.Next() || p.lex.Token() != "}" {
		return nil, errors.New("expected token \"}\"")
	}
	root.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))

	return root, nil
}

func (p *parser) operators() (*Node, error) {
	p.log("operators()")

	root := NewNode("operators")

	operatorNode, err := p.operator()
	if err != nil {
		return nil, err
	}
	root.AddChild(operatorNode)

	if p.lex.Next() && p.lex.Token() == ";" {
		_ = p.lex.Rollback()
		//tailNode := NewNode("tail")
		//tailNode.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))
		tailNode, err := p.tail()
		if err != nil {
			return nil, err
		}
		root.AddChild(tailNode)
		return root, nil
	}
	_ = p.lex.Rollback()

	return root, nil
}

func (p *parser) tail() (*Node, error) {
	p.log("tail()")

	root := NewNode("tail")

	if !p.lex.Next() || p.lex.Token() != ";" {
		_ = p.lex.Rollback()
		return nil, fmt.Errorf("expected \";\" %w", ErrWrongToken)
	}
	root.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))

	node, err := p.operator()
	if err != nil {
		_ = p.lex.Rollback()
		return nil, err
	}
	root.AddChild(node)

	ok := p.lex.Next()
	token := p.lex.Token()
	_ = p.lex.Rollback()
	if ok && token == ";" {
		node, err := p.tail()
		if err != nil {
			return nil, err
		}
		root.AddChild(node)
		return root, nil
	}

	return root, nil
}

func (p *parser) operator() (*Node, error) {
	p.log("operator()")
	root := NewNode("operator")

	if !p.lex.Next() || p.lex.Token() != "id" {
		_ = p.lex.Rollback()
		return nil, errors.New(`id expected`)
	}
	root.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))

	if !p.lex.Next() || p.lex.Token() != "=" {
		_ = p.lex.Rollback()
		_ = p.lex.Rollback()
		return nil, errors.New(`"=" expected`)
	}
	root.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))

	node, err := p.expr()
	if err != nil {
		_ = p.lex.Rollback()
		return nil, err
	}
	root.AddChild(node)

	return root, nil
}

func (p *parser) expr() (*Node, error) {
	p.log("expr()")

	root := NewNode("expr")

	arithmExprNode1, err := p.arithmExpr()
	if err != nil {
		return nil, err
	}
	root.AddChild(arithmExprNode1)

	signCmpNode, err := p.signCmp()
	if err != nil {
		return nil, err
	}
	root.AddChild(signCmpNode)

	arithmExprNode2, err := p.arithmExpr()
	if err != nil {
		return nil, err
	}
	root.AddChild(arithmExprNode2)

	return root, nil
}

func (p *parser) arithmExpr() (*Node, error) {
	p.log("arithmExpr()")
	root := NewNode("arithmExpr")

	if signAdd, err := p.signAdd(); err == nil {
		root.AddChild(signAdd)
	}

	termNode, err := p.term()
	if err != nil {
		return nil, err
	}
	root.AddChild(termNode)

	if arithmExprNode_, err := p.arithmExpr_(); err == nil {
		root.AddChild(arithmExprNode_)
	}

	return root, nil
}

func (p *parser) arithmExpr_() (*Node, error) {
	p.log("arithmExpr_()")
	root := NewNode("arithmExpr_")

	signAddNode, err := p.signAdd()
	if err != nil {
		return nil, err
	}
	root.AddChild(signAddNode)

	termNode, err := p.term()
	if err != nil {
		return nil, err
	}
	root.AddChild(termNode)

	if arithmExprNode_, err := p.arithmExpr_(); err == nil {
		root.AddChild(arithmExprNode_)
	}

	return root, nil
}

func (p *parser) term() (*Node, error) {
	p.log("term()")
	root := NewNode("term")
	multiplier, err := p.multiplier()
	if err != nil {
		return nil, err
	}
	root.AddChild(multiplier)

	if termNode_, err := p.term_(); err == nil {
		root.AddChild(termNode_)
	}
	return root, nil
}

func (p *parser) term_() (*Node, error) {
	p.log("term_()")
	root := NewNode("term_")

	signMultiNode, err := p.signMulti()
	if err != nil {
		return nil, err
	}
	root.AddChild(signMultiNode)
	multiplierNode, err := p.multiplier()
	if err != nil {
		return nil, err
	}
	root.AddChild(multiplierNode)
	if termNode_, err := p.term_(); err == nil {
		root.AddChild(termNode_)
	}
	return root, nil
}

func (p *parser) multiplier() (*Node, error) {
	p.log("multiplier()")
	root := NewNode("multiplier")
	primeExprNode, err := p.primeExpr()
	if err != nil {
		return nil, err
	}
	root.AddChild(primeExprNode)
	if multiplierNode_, err := p.multiplier_(); err == nil {
		root.AddChild(multiplierNode_)
	}

	return root, nil
}

func (p *parser) multiplier_() (*Node, error) {
	p.log("multiplier_()")
	root := NewNode("multiplier_")

	if !p.lex.Next() || p.lex.Token() != "^" {
		_ = p.lex.Rollback()
		return nil, fmt.Errorf("need \"^\" err:%w", ErrWrongToken)
	}
	root.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))

	primeExprNode, err := p.primeExpr()
	if err != nil {
		_ = p.lex.Rollback()
		return nil, err
	}
	root.AddChild(primeExprNode)

	if multiplierNode_, err := p.multiplier_(); err == nil {
		root.AddChild(multiplierNode_)
	}

	return root, nil
}

func (p *parser) primeExpr() (*Node, error) {
	p.log("primeExpr()")
	root := NewNode("primeExpr")

	if !p.lex.Next() {
		return nil, ErrNextFail
	}

	switch p.lex.Token() {
	case "number", "id":
		root.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))
		return root, nil
	case "(":
		arithmExprNode, err := p.arithmExpr()
		if err != nil {
			_ = p.lex.Rollback()
			return nil, err
		}
		root.AddChild(arithmExprNode)
		if !p.lex.Next() || p.lex.Token() != ")" {
			_ = p.lex.Rollback()
			_ = p.lex.Rollback()
			return nil, fmt.Errorf("need \")\" err:%w", ErrWrongToken)
		}
		root.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))
		return root, nil
	default:
		_ = p.lex.Rollback()
		return nil, fmt.Errorf("<первичное выражение> err: %w", ErrWrongToken)
	}
}

func (p *parser) signAdd() (*Node, error) {
	p.log("signAdd()")
	root := NewNode("signAdd")

	if !p.lex.Next() || (p.lex.Token() != "+" && p.lex.Token() != "-") {
		_ = p.lex.Rollback()
		return nil, fmt.Errorf("need \"+\" or \"-\" err:%w", ErrWrongToken)
	}
	root.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))
	return root, nil
}

func (p *parser) signMulti() (*Node, error) {
	p.log("signMulti()")
	root := NewNode("signMulti")

	if !p.lex.Next() {
		return nil, ErrNextFail
	}
	switch p.lex.Token() {
	case "*", "/", "%":
		root.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))
		return root, nil
	default:
		_ = p.lex.Rollback()
		return nil, fmt.Errorf("need \"* | / | %% \" err:%w", ErrWrongToken)
	}
}

func (p *parser) signCmp() (*Node, error) {
	p.log("signCmp()")
	root := NewNode("signCmp")

	if !p.lex.Next() {
		return nil, ErrNextFail
	}
	switch p.lex.Token() {
	case "<", "<=", "=", ">=", ">", "<>":

		root.AddChild(NewNode(fmt.Sprintf("%s_%d", p.lex.Token(), p.lex.CurPos())))
		return root, nil
	default:
		_ = p.lex.Rollback()
		return nil, fmt.Errorf("need \"< | <= | = | >= | > | <>\" err:%w", ErrWrongToken)
	}
}
