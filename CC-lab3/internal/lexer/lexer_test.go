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
			for r.Next() {
				actTokens = append(actTokens, r.Token())
			}
			require.Equal(t, tt.wantTokens, actTokens)
		})
	}
}
