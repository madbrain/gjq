package lang

import (
	"testing"
)

func checkLexer(lexer *Lexer, reporter *DefaultReporter, expected []Token, expectedMessages []Error, t *testing.T) {
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

	if len(reporter.errors) != len(expectedMessages) {
		t.Fatalf("bad message count, expecting %d got %d", len(expectedMessages), len(reporter.errors))
	}

	for i, message := range reporter.errors {
		if message.span != expectedMessages[i].span || message.message != expectedMessages[i].message {
			t.Fatalf("bad message %+v", message)
		}
	}
}

func TestLexingPath(t *testing.T) {
	var reporter = &DefaultReporter{}
	var lexer = NewLexer(".tutu\t[1] &*.toto", reporter)

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
	var expectedMessages = []Error{
		{span: Span{Start: 10, End: 12}, message: "Unrecognized character(s)"},
	}

	checkLexer(lexer, reporter, expected, expectedMessages, t)
}

func TestLexingFunctionCall(t *testing.T) {
	var reporter = &DefaultReporter{}
	var lexer = NewLexer(".tutu(1,2).toto", reporter)

	var expected = []Token{
		{span: Span{Start: 0, End: 1}, kind: DOT},
		{span: Span{Start: 1, End: 5}, kind: IDENT, value: "tutu"},
		{span: Span{Start: 5, End: 6}, kind: LEFT_PAR},
		{span: Span{Start: 6, End: 7}, kind: INTEGER, value: "1"},
		{span: Span{Start: 7, End: 8}, kind: COMA},
		{span: Span{Start: 8, End: 9}, kind: INTEGER, value: "2"},
		{span: Span{Start: 9, End: 10}, kind: RIGHT_PAR},
		{span: Span{Start: 10, End: 11}, kind: DOT},
		{span: Span{Start: 11, End: 15}, kind: IDENT, value: "toto"},
		{span: Span{Start: 15, End: 15}, kind: EOF},
	}

	checkLexer(lexer, reporter, expected, nil, t)
}
