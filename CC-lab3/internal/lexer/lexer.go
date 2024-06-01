package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"unicode"
)

type Lexer interface {
	Next() bool
	Token() Token
	CurPos() int
	Rollback() error
}

type lexer struct {
	r      *bufio.Reader
	curPos int
	err    error
	tokens []Token
}

func New(r io.Reader) Lexer {
	read := &lexer{
		r: bufio.NewReader(r),
	}
	return read
}

func (l *lexer) CurPos() int {
	slog.Debug("tokens", slog.Any("tokens", l.tokens))
	return l.curPos
}

func (l *lexer) Rollback() error {
	l.curPos--
	if l.curPos < 0 {
		return fmt.Errorf("rollback tokens is empty; curPos=%d", l.curPos)
	}
	return nil
}

func (l *lexer) Next() (ok bool) {
	if l.curPos < len(l.tokens) {
		l.curPos++
		return true
	}
	if err := l.skipSpace(); err != nil {
		l.err = err
		slog.Warn("skipSpace", "err", err.Error())
		return false
	}
	sb := strings.Builder{}
	for {
		b, err := l.r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return false
			}
			l.err = fmt.Errorf("read byte %w", err)
			return false
		}
		if stop(b) {
			return false
		}
		sb.WriteByte(b)
		if _b, ok := l.oneMore(b); ok {
			sb.WriteByte(_b)
		}
		if isToken(sb.String()) {
			l.tokens = append(l.tokens, (Token)(sb.String()))
			l.curPos++
			return true
		}
	}
}

func stop(b byte) bool {
	return unicode.IsSpace(rune(b))
}

func (l *lexer) oneMore(b byte) (byte, bool) {
	if b == '<' || b == '>' {
		_b, err := l.r.ReadByte()
		if err != nil && !errors.Is(err, io.EOF) {
			slog.Warn("oneMore ReadByte", "err", err.Error())
			return 0, false
		}
		if isToken(string([]byte{b, _b})) {
			return _b, true
		}
		if err = l.r.UnreadByte(); err != nil {
			slog.Warn("oneMore UnreadByte()", "err", err.Error())
			return 0, false
		}
	}
	return 0, false
}

func (l *lexer) skipSpace() error {
	for {
		b, err := l.r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if !unicode.IsSpace(rune(b)) {
			if err := l.r.UnreadByte(); err != nil {
				return err
			}
			return nil
		}
	}
}

func (l *lexer) Token() Token {
	if l.curPos-1 < 0 {
		return "empty l.curPos-1"
	}
	return l.tokens[l.curPos-1]
}
