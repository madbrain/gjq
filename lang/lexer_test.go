package lang

import (
	"testing"
)

func TestLexer(t *testing.T) {
	var reporter = DefaultReporter{}
	var lexer = NewLexer(".tutu\t[1] &*.toto", &reporter)

	var expected = []Token{
		{span: Span{Start: 0, End: 1}, kind: DOT},
		{span: Span{Start: 1, End: 5}, kind: IDENT, value: "tutu"},
		{span: Span{Start: 6, End: 7}, kind: LEFT_BRT},
		{span: Span{Start: 7, End: 8}, kind: INTEGER, value: "1"},
		{span: Span{Start: 8, End: 9}, kind: RIGHT_BRT},
		{span: Span{Start: 12, End: 13}, kind: DOT},
		{span: Span{Start: 13, End: 17}, kind: IDENT, value: "toto"},
		{span: Span{Start: 17, End: 17}, kind: EOF},
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
	var expectedMessage = Error{span: Span{Start: 10, End: 12}, message: "Unrecognized character(s)"}
	if reporter.errors[0] != expectedMessage {
		t.Fatalf("bad message %+v", reporter.errors[0])
	}
}
