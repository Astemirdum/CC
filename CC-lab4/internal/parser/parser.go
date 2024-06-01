package parser

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/Astemirdum/CC-lab4/internal/lexer"
	"github.com/Astemirdum/CC-lab4/pkg"
)

type OperatorPrecedenceParser interface {
	Parse() ([]string, error)
}

func New(lex lexer.Lexer) OperatorPrecedenceParser {
	return &operatorPrecedenceParser{
		lex:   lex,
		stack: pkg.NewStack[lexer.Token](),
	}
}

type operatorPrecedenceParser struct {
	lex   lexer.Lexer
	stack pkg.Stack[lexer.Token]
	res   []string
}

var (
	ErrNextToken           = errors.New("next token")
	ErrUnknownOpPrecedence = errors.New("unknown op precedence")
)

func (p *operatorPrecedenceParser) Parse() ([]string, error) {
	var err error
	token, ok := p.lex.Next()
	if !ok {
		return nil, ErrNextToken
	}
	p.stack.Push(lexer.TokenEnd)

	for p.stack.Len() > 1 || token != lexer.TokenEnd {
		slog.Debug("", slog.String("stack", fmt.Sprintf("%v", p.stack)[3:p.stack.Len()+3]))
		slog.Debug("", slog.Any("top", p.stack.Top()), slog.Any("token", token), slog.Any("op", lexer.OpTable(p.stack.Top(), token)))

		switch lexer.OpTable(p.stack.Top(), token) {
		case lexer.OpPrecedenceLess, lexer.OpPrecedenceEqual:
			token, err = p.shift(token)
			if err != nil {
				return nil, fmt.Errorf("shift token %s %v", token, err)
			}
		case lexer.OpPrecedenceMore:
			if err := p.reduce(); err != nil {
				return nil, fmt.Errorf("reduce token %s %v", token, err)
			}
		default:
			return nil, ErrUnknownOpPrecedence
		}
	}

	return p.res, nil
}

func (p *operatorPrecedenceParser) shift(token lexer.Token) (lexer.Token, error) {
	p.stack.Push(token)
	token, ok := p.lex.Next()
	if !ok {
		return "", ErrNextToken
	}
	return token, nil
}

func (p *operatorPrecedenceParser) reduce() error {
	for p.stack.Len() > 0 {
		token := p.stack.Pop()
		if token != lexer.TokenBracketOpen && token != lexer.TokenBracketClose {
			p.res = append(p.res, (string)(token))
		}
		if lexer.OpTable(p.stack.Top(), token) == lexer.OpPrecedenceLess {
			break
		}
	}
	return nil
}
