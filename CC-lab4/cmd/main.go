package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/Astemirdum/CC-lab4/internal/lexer"
	"github.com/Astemirdum/CC-lab4/internal/parser"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("need args - file input")
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelInfo})))

	inp := os.Args[1]
	if err := run(inp); err != nil {
		log.Fatal("run", err)
	}
	log.Println("GOOD JOB")
}

func run(inp string) error {
	f, err := os.Open(inp)
	if err != nil {
		return err
	}
	defer f.Close()

	p := parser.New(lexer.New(f))
	dotData, err := p.Parse()
	if err != nil {
		return err
	}

	fmt.Println(dotData)
	return nil
}
