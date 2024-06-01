package parser

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/Astemirdum/CC-lab4/internal/lexer"
	"github.com/stretchr/testify/require"
)

func Test_parser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		wantRes []string
		wantErr bool
	}{
		{
			name: "ok",
			r: strings.NewReader(`
    id = number = id
`),
			wantRes: []string{"id", "number", "=", "id", "="},
		},
		{
			name: "ok",
			r: strings.NewReader(`
    id = number + number * number - id
`),
			wantRes: []string{"id", "number", "number", "number", "*", "+", "id", "-", "="},
		},
		{
			name: "ok",
			r: strings.NewReader(`
    (number - id) >= number
`),
			wantRes: []string{"number", "id", "-", "number", ">="},
		},
		{
			name: "ok",
			r: strings.NewReader(`
    id <= (number + number) < number
`),
			wantRes: []string{"id", "number", "number", "+", "<=", "number", "<"},
		},
		{
			name: "ok",
			r: strings.NewReader(`
	id = (number + (number ^ number)) > number
`),
			wantRes: []string{"id", "number", "number", "number", "^", "+", "number", ">", "="},
		},
		{
			name: "ok",
			r: strings.NewReader(`
{
    id = asd
}
`),
			wantErr: true,
		},
		{
			name: "ok",
			r: strings.NewReader(`
    id = (number + id ^ (number - id)) <= (number * id + number)
`),
			wantRes: []string{"id", "number", "id", "number", "id", "-", "^", "+", "number", "id", "*", "number", "+", "<=", "="},
		},
	}

	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(lexer.New(tt.r))
			rpn, err := p.Parse()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantRes, rpn)
			}
		})
	}
}
