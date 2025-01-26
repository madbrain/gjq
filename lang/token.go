package lang

type TokenKind = int

const (
	DOT TokenKind = iota
	LEFT_BRT
	RIGHT_BRT

	IDENT
	INTEGER
	EOF
)

type Token struct {
	span  Span
	kind  TokenKind
	value string
}
