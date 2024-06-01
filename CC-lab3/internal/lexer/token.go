package lexer

type Token string

var tokens = map[Token]struct{}{
	"id":     {},
	"number": {},
	"{":      {},
	"}":      {},
	"*":      {},
	"+":      {},
	"-":      {},
	"/":      {},
	"%":      {},
	"=":      {},
	"<":      {},
	"<=":     {},
	">=":     {},
	">":      {},
	"<>":     {},
	")":      {},
	"(":      {},
	"^":      {},
	";":      {},
}

func isToken(t string) bool {
	_, ok := tokens[Token(t)]
	return ok
}
