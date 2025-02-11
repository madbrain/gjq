package lang

type TokenKind = int

const (
	DOT TokenKind = iota
	COMA
	LEFT_BRT
	RIGHT_BRT
	LEFT_PAR
	RIGHT_PAR

	IDENT
	INTEGER
	EOF
)

type Token struct {
	span  Span
	kind  TokenKind
	value string
}
