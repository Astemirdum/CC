package parser

import (
	"fmt"
	"github.com/Astemirdum/CC-lab3/internal/lexer"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func Test_parser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		wantErr bool
	}{

		{
			name: "ok",
			r: strings.NewReader(`
{
    id = (number + id) < id
}
`),
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
{
    id = (number + id ^ (number - id)) <= (number * id + number)
}
`),
		},

		{
			name: "ok",
			r: strings.NewReader(`
{
    id = number = id;
    id = number > id;
    id = (number) <> number;
    id = (number + id ^ (number - id)) <= (number * id + number)
}`),
			wantErr: false,
		},
	}

	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(lexer.New(tt.r))
			dotData, err := p.Parse()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				fmt.Println(dotData)
				//err = os.WriteFile("graph.dot", []byte(dotData), 0644)
				//if err != nil {
				//	panic(err)
				//}
				//cmd := exec.Command("dot", "-Tpng", "graph.dot", "-o", "graph.png")
				//if err := cmd.Run(); err != nil {
				//	panic(err)
				//}
			}
		})
	}
}
