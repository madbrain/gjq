package lang

import "fmt"

type Expr interface {
	Span() Span
}

type Start struct {
	span Span
}

func (s *Start) Span() Span {
	return s.span
}

type FieldAccess struct {
	span  Span
	expr  Expr
	field Token
}

func (s *FieldAccess) Span() Span {
	return s.span
}

type ArrayAccess struct {
	span  Span
	expr  Expr
	index Expr
}

func (s *ArrayAccess) Span() Span {
	return s.span
}

type IntegerValue struct {
	span  Span
	value string
}

func (s *IntegerValue) Span() Span {
	return s.span
}

func displaySpan(s Span) string {
	return fmt.Sprintf("{%d:%d}", s.start, s.end)
}

func displayAst(e Expr) string {
	switch t := e.(type) {
	case *Start:
		return displaySpan(t.span)
	case *FieldAccess:
		return fmt.Sprintf("(%s%s.%s%s)", displaySpan(t.span), displayAst(t.expr), displaySpan(t.field.span), t.field.value)
	case *ArrayAccess:
		return fmt.Sprintf("(%s%s[%s])", displaySpan(t.span), displayAst(t.expr), displayAst(t.index))
	case *IntegerValue:
		return fmt.Sprintf("(%s%s)", displaySpan(t.span), t.value)
	default:
		panic("AST is unknown\n")
	}
}
