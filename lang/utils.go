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
		return fmt.Sprintf("(%s%s.%s%s)", displaySpan(t.span), DisplayAst(t.expr), displaySpan(t.field.span), t.field.value)
	case *BadFieldAccess:
		return fmt.Sprintf("(%s.@)", displaySpan(t.span))
	case *ArrayAccess:
		return fmt.Sprintf("(%s%s[%s])", displaySpan(t.span), DisplayAst(t.expr), DisplayAst(t.index))
	case *BadExpr:
		return fmt.Sprintf("(%s@)", displaySpan(t.span))
	case *IntegerValue:
		return fmt.Sprintf("(%s%s)", displaySpan(t.span), t.value)
	default:
		panic("AST is unknown\n")
	}
}
