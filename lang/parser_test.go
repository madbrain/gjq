package lang

import (
	"fmt"
	"testing"
)

func compareSpan(a Span, b Span, t *testing.T) bool {
	if a.Start != b.Start || a.End != b.End {
		t.Fatalf("not same Span %+v != %+v\n", a, b)
		return false
	}
	return true
}

func compareIdentifier(a Identifier, b Identifier, t *testing.T) bool {
	if !compareSpan(a.Span, b.Span, t) || a.Value != b.Value {
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
			return compareSpan(t.span, u.span, x) && compareAst(t.Expr, u.Expr, x) && compareIdentifier(t.Field, u.Field, x)
		default:
			x.Fatalf("not same FieldAccess %+v != %+v\n", t, u)
			return false
		}
	case *BadFieldAccess:
		switch u := b.(type) {
		case *BadFieldAccess:
			return compareSpan(t.span, u.span, x) && compareAst(t.Expr, u.Expr, x)
		default:
			x.Fatalf("not same BadFieldAccess %+v != %+v\n", t, u)
			return false
		}
	case *ArrayAccess:
		switch u := b.(type) {
		case *ArrayAccess:
			return compareSpan(t.span, u.span, x) && compareAst(t.Expr, u.Expr, x) && compareAst(t.Index, u.Index, x)
		default:
			x.Fatalf("not same ArrayAccess %+v != %+v\n", t, u)
			return false
		}
	case *IntegerValue:
		switch u := b.(type) {
		case *IntegerValue:
			return compareSpan(t.span, u.span, x) && t.Value == u.Value
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
	var reporter = DefaultReporter{}
	var lexer = NewLexer(content, &reporter)
	var parser = NewParser(lexer, &reporter)

	var expectedAst = &FieldAccess{
		span: Span{Start: 0, End: 17},
		Expr: &ArrayAccess{
			span: Span{Start: 0, End: 9},
			Expr: &FieldAccess{
				span:  Span{Start: 0, End: 5},
				Expr:  &Start{span: Span{Start: 0, End: 0}},
				Field: Identifier{Span: Span{Start: 1, End: 5}, Value: "tutu"},
			},
			Index: &IntegerValue{span: Span{Start: 7, End: 8}, Value: "1"},
		},
		Field: Identifier{Span: Span{Start: 13, End: 17}, Value: "toto"},
	}

	var ast = parser.Parse()

	if !compareAst(ast, expectedAst, t) {
		t.Fatalf("error in parsing")
	}
}

func TestParserRecovery(t *testing.T) {
	var content = ".tutu foo.toto"
	var reporter = DefaultReporter{}
	var lexer = NewLexer(content, &reporter)
	var parser = NewParser(lexer, &reporter)

	var expectedAst = &FieldAccess{
		span: Span{Start: 0, End: 14},
		Expr: &FieldAccess{
			span:  Span{Start: 0, End: 5},
			Expr:  &Start{span: Span{Start: 0, End: 0}},
			Field: Identifier{Span: Span{Start: 1, End: 5}, Value: "tutu"},
		},
		Field: Identifier{Span: Span{Start: 10, End: 14}, Value: "toto"},
	}

	var ast = parser.Parse()

	reporter.DisplayErrors(content)

	if !compareAst(ast, expectedAst, t) {
		t.Fatalf("error in parsing")
	}
}

func TestParserRecoverFieldAccess(t *testing.T) {
	var content = ".[10]"
	var reporter = DefaultReporter{}
	var lexer = NewLexer(content, &reporter)
	var parser = NewParser(lexer, &reporter)

	var expectedAst = &ArrayAccess{
		span: Span{Start: 0, End: 5},
		Expr: &BadFieldAccess{
			span: Span{Start: 0, End: 1},
			Expr: &Start{span: Span{Start: 0, End: 0}},
		},
		Index: &IntegerValue{span: Span{Start: 2, End: 4}, Value: "10"},
	}

	var ast = parser.Parse()

	fmt.Println(DisplayAst(ast))

	reporter.DisplayErrors(content)

	if !compareAst(ast, expectedAst, t) {
		t.Fatalf("error in parsing")
	}
}

func TestParserRecoverArrayAccess(t *testing.T) {
	var content = "[]"
	var reporter = DefaultReporter{}
	var lexer = NewLexer(content, &reporter)
	var parser = NewParser(lexer, &reporter)

	var expectedAst = &ArrayAccess{
		span:  Span{Start: 0, End: 2},
		Expr:  &Start{span: Span{Start: 0, End: 0}},
		Index: &BadExpr{span: Span{Start: 1, End: 1}},
	}

	var ast = parser.Parse()

	fmt.Println(DisplayAst(ast))

	reporter.DisplayErrors(content)

	if !compareAst(ast, expectedAst, t) {
		t.Fatalf("error in parsing")
	}
}

func TestParserRecoverUnterminatedArrayAccess(t *testing.T) {
	var content = "[10"
	var reporter = DefaultReporter{}
	var lexer = NewLexer(content, &reporter)
	var parser = NewParser(lexer, &reporter)

	var expectedAst = &ArrayAccess{
		span:  Span{Start: 0, End: 3},
		Expr:  &Start{span: Span{Start: 0, End: 0}},
		Index: &IntegerValue{span: Span{Start: 1, End: 3}, Value: "10"},
	}

	var ast = parser.Parse()

	fmt.Println(DisplayAst(ast))

	reporter.DisplayErrors(content)

	if !compareAst(ast, expectedAst, t) {
		t.Fatalf("error in parsing")
	}
}
