package lang

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	var reporter = TestReporter{}
	var lexer = NewLexer(".tutu\t[1] &*.toto", &reporter)
	var parser = NewParser(lexer, &reporter)

	var expectedAst = &FieldAccess{
		span: Span{start: 0, end: 17},
		expr: &ArrayAccess{
			span: Span{start: 0, end: 9},
			expr: &FieldAccess{
				span:  Span{start: 0, end: 5},
				expr:  &Start{span: Span{start: 0, end: 0}},
				field: Token{span: Span{start: 1, end: 5}, kind: IDENT, value: "tutu"},
			},
			index: &IntegerValue{span: Span{start: 7, end: 8}, value: "1"},
		},
		field: Token{span: Span{start: 13, end: 17}, kind: IDENT, value: "toto"},
	}

	var ast = parser.Parse()

	fmt.Println(displayAst(ast))

	if ast != expectedAst {
		t.Fatalf("error in parsing %+v", ast)
	}
}
