package lang

import "fmt"

func displaySpan(s Span) string {
	return fmt.Sprintf("{%d:%d}", s.start, s.end)
}

func DisplayAst(e Expr) string {
	switch t := e.(type) {
	case *Start:
		return displaySpan(t.span)
	case *FieldAccess:
		return fmt.Sprintf("(%s%s.%s%s)", displaySpan(t.span), DisplayAst(t.Expr), displaySpan(t.Field.span), t.Field.Value)
	case *BadFieldAccess:
		return fmt.Sprintf("(%s.@)", displaySpan(t.span))
	case *ArrayAccess:
		return fmt.Sprintf("(%s%s[%s])", displaySpan(t.span), DisplayAst(t.Expr), DisplayAst(t.Index))
	case *BadExpr:
		return fmt.Sprintf("(%s@)", displaySpan(t.span))
	case *IntegerValue:
		return fmt.Sprintf("(%s%s)", displaySpan(t.span), t.Value)
	default:
		panic("AST is unknown\n")
	}
}
