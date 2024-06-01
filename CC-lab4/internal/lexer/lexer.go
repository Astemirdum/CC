package lexer

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"strings"
	"unicode"
)

type Lexer interface {
	Next() (Token, bool)
}

type lexer struct {
	r      *bufio.Reader
	curPos int
}

func New(r io.Reader) Lexer {
	read := &lexer{
		r: bufio.NewReader(r),
	}
	return read
}

func (l *lexer) Next() (Token, bool) {
	if err := l.skipSpace(); err != nil {
		slog.Warn("skipSpace", "err", err.Error())
		return "", false
	}
	sb := strings.Builder{}
	for {
		b, err := l.r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return "$", true
			}
			return "", false
		}
		if stop(b) {
			return "", false
		}
		sb.WriteByte(b)
		if _b, ok := l.oneMore(b); ok {
			sb.WriteByte(_b)
		}
		if isToken(sb.String()) {
			return (Token)(sb.String()), true
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
