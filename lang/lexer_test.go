package lang

import (
	"testing"
)

func TestLexer(t *testing.T) {
	var reporter = DefaultReporter{}
	var lexer = NewLexer(".tutu\t[1] &*.toto", &reporter)

	var expected = []Token{
		{span: Span{start: 0, end: 1}, kind: DOT},
		{span: Span{start: 1, end: 5}, kind: IDENT, value: "tutu"},
		{span: Span{start: 6, end: 7}, kind: LEFT_BRT},
		{span: Span{start: 7, end: 8}, kind: INTEGER, value: "1"},
		{span: Span{start: 8, end: 9}, kind: RIGHT_BRT},
		{span: Span{start: 12, end: 13}, kind: DOT},
		{span: Span{start: 13, end: 17}, kind: IDENT, value: "toto"},
		{span: Span{start: 17, end: 17}, kind: EOF},
	}
	var index = 0

	for {
		var token = lexer.NextToken()
		if index >= len(expected) {
			t.Fatal("too much tokens")
		}
		if *token != expected[index] {
			t.Fatalf("[%d] expecting token %+v, got %+v", index, expected[index], *token)
		}
		if token.kind == EOF {
			break
		}
		index += 1
	}
	if len(reporter.errors) != 1 {
		t.Fatalf("expecting one error message")
	}
	var expectedMessage = Error{span: Span{start: 10, end: 12}, message: "Unrecognized character(s)"}
	if reporter.errors[0] != expectedMessage {
		t.Fatalf("bad message %+v", reporter.errors[0])
	}
}
