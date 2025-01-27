package lang

import (
	"fmt"
	"testing"
)

func compareSpan(a Span, b Span, t *testing.T) bool {
	if a.start != b.start || a.end != b.end {
		t.Fatalf("not same Span %+v != %+v\n", a, b)
		return false
	}
	return true
}

func compareIdentifier(a Identifier, b Identifier, t *testing.T) bool {
	if !compareSpan(a.span, b.span, t) || a.value != b.value {
		t.Fatalf("not same Identifier %+v != %+v\n", a, b)
		return false
	}
	return true
}

func compareAst(a Expr, b Expr, x *testing.T) bool {
	switch t := a.(type) {
	case *Start:
		switch u := b.(type) {
		case *Start:
			return compareSpan(t.span, u.span, x)
		default:
			x.Fatalf("not same Start %+v != %+v\n", t, u)
			return false
		}
	case *FieldAccess:
		switch u := b.(type) {
		case *FieldAccess:
			return compareSpan(t.span, u.span, x) && compareAst(t.expr, u.expr, x) && compareIdentifier(t.field, u.field, x)
		default:
			x.Fatalf("not same FieldAccess %+v != %+v\n", t, u)
			return false
		}
	case *BadFieldAccess:
		switch u := b.(type) {
		case *BadFieldAccess:
			return compareSpan(t.span, u.span, x) && compareAst(t.expr, u.expr, x)
		default:
			x.Fatalf("not same BadFieldAccess %+v != %+v\n", t, u)
			return false
		}
	case *ArrayAccess:
		switch u := b.(type) {
		case *ArrayAccess:
			return compareSpan(t.span, u.span, x) && compareAst(t.expr, u.expr, x) && compareAst(t.index, u.index, x)
		default:
			x.Fatalf("not same ArrayAccess %+v != %+v\n", t, u)
			return false
		}
	case *IntegerValue:
		switch u := b.(type) {
		case *IntegerValue:
			return compareSpan(t.span, u.span, x) && t.value == u.value
		default:
			x.Fatalf("not same IntegerValue %+v != %+v\n", t, u)
			return false
		}
	case *BadExpr:
		switch u := b.(type) {
		case *BadExpr:
			return compareSpan(t.span, u.span, x)
		default:
			x.Fatalf("not same BadExpr %+v != %+v\n", t, u)
			return false
		}
	default:
		x.Fatalf("AST type is unknown %+v\n", a)
		return false
	}
}

func TestParser(t *testing.T) {
	var content = ".tutu\t[1] &*.toto"
	var reporter = TestReporter{}
	var lexer = NewLexer(content, &reporter)
	var parser = NewParser(lexer, &reporter)

	var expectedAst = &FieldAccess{
		span: Span{start: 0, end: 17},
		expr: &ArrayAccess{
			span: Span{start: 0, end: 9},
			expr: &FieldAccess{
				span:  Span{start: 0, end: 5},
				expr:  &Start{span: Span{start: 0, end: 0}},
				field: Identifier{span: Span{start: 1, end: 5}, value: "tutu"},
			},
			index: &IntegerValue{span: Span{start: 7, end: 8}, value: "1"},
		},
		field: Identifier{span: Span{start: 13, end: 17}, value: "toto"},
	}

	var ast = parser.Parse()

	if !compareAst(ast, expectedAst, t) {
		t.Fatalf("error in parsing")
	}
}

func TestParserRecovery(t *testing.T) {
	var content = ".tutu foo.toto"
	var reporter = TestReporter{}
	var lexer = NewLexer(content, &reporter)
	var parser = NewParser(lexer, &reporter)

	var expectedAst = &FieldAccess{
		span: Span{start: 0, end: 14},
		expr: &FieldAccess{
			span:  Span{start: 0, end: 5},
			expr:  &Start{span: Span{start: 0, end: 0}},
			field: Identifier{span: Span{start: 1, end: 5}, value: "tutu"},
		},
		field: Identifier{span: Span{start: 10, end: 14}, value: "toto"},
	}

	var ast = parser.Parse()

	//reporter.Display(content)

	if !compareAst(ast, expectedAst, t) {
		t.Fatalf("error in parsing")
	}
}

func TestParserRecoverFieldAccess(t *testing.T) {
	var content = ".[10]"
	var reporter = TestReporter{}
	var lexer = NewLexer(content, &reporter)
	var parser = NewParser(lexer, &reporter)

	var expectedAst = &ArrayAccess{
		span: Span{start: 0, end: 5},
		expr: &BadFieldAccess{
			span: Span{start: 0, end: 1},
			expr: &Start{span: Span{start: 0, end: 0}},
		},
		index: &IntegerValue{span: Span{start: 2, end: 4}, value: "10"},
	}

	var ast = parser.Parse()

	fmt.Println(displayAst(ast))

	reporter.Display(content)

	if !compareAst(ast, expectedAst, t) {
		t.Fatalf("error in parsing")
	}
}

func TestParserRecoverArrayAccess(t *testing.T) {
	var content = "[]"
	var reporter = TestReporter{}
	var lexer = NewLexer(content, &reporter)
	var parser = NewParser(lexer, &reporter)

	var expectedAst = &ArrayAccess{
		span:  Span{start: 0, end: 2},
		expr:  &Start{span: Span{start: 0, end: 0}},
		index: &BadExpr{span: Span{start: 1, end: 1}},
	}

	var ast = parser.Parse()

	fmt.Println(displayAst(ast))

	reporter.Display(content)

	if !compareAst(ast, expectedAst, t) {
		t.Fatalf("error in parsing")
	}
}

func TestParserRecoverUnterminatedArrayAccess(t *testing.T) {
	var content = "[10"
	var reporter = TestReporter{}
	var lexer = NewLexer(content, &reporter)
	var parser = NewParser(lexer, &reporter)

	var expectedAst = &ArrayAccess{
		span:  Span{start: 0, end: 3},
		expr:  &Start{span: Span{start: 0, end: 0}},
		index: &IntegerValue{span: Span{start: 1, end: 3}, value: "10"},
	}

	var ast = parser.Parse()

	fmt.Println(displayAst(ast))

	reporter.Display(content)

	if !compareAst(ast, expectedAst, t) {
		t.Fatalf("error in parsing")
	}
}
