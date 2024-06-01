package lexer

import (
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

func Test_Tokenize(t *testing.T) {
	tests := []struct {
		name       string
		reader     io.Reader
		wantTokens []Token
	}{
		{
			name: "ok",
			reader: strings.NewReader(`
{
    id = number = id;
}
`),
			wantTokens: []Token{"{", "id", "=", "number", "=", "id", ";", "}"},
		},
		{
			name: "ok",
			reader: strings.NewReader(`
{
    id = (number) <> number;
	id = (number) <= number;
}
`),
			wantTokens: []Token{"{", "id", "=", "(", "number", ")", "<>", "number", ";",
				"id", "=", "(", "number", ")", "<=", "number", ";", "}"},
		},
		{
			name: "ok",
			reader: strings.NewReader(`
{
    id = number = id;
    id = number > id;
    id = (number + id ^ (number - id)) <= (number * id + number);
    id = (number) <> number
}
`),
			wantTokens: []Token{
				"{", "id", "=", "number", "=", "id", ";",
				"id", "=", "number", ">", "id", ";",
				"id", "=", "(", "number", "+", "id", "^", "(", "number", "-", "id", ")", ")", "<=", "(", "number", "*", "id", "+", "number", ")", ";",
				"id", "=", "(", "number", ")", "<>", "number", "}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.reader)
			var actTokens []Token
			for token, ok := r.Next(); ok; token, ok = r.Next() {
				if token == TokenEnd {
					break
				}
				actTokens = append(actTokens, token)
			}
			require.Equal(t, tt.wantTokens, actTokens)
		})
	}
}
