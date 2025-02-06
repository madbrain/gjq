package lang

import "fmt"

func displaySpan(s Span) string {
	return fmt.Sprintf("{%d:%d}", s.Start, s.End)
}

func DisplayExprs(exprs []Expr) string {
	result := ""
	for i, expr := range exprs {
		if i > 0 {
			result += ", "
		}
		result += DisplayAst(expr)
	}
	return result
}

func DisplayAst(e Expr) string {
	switch t := e.(type) {
	case *Start:
		return displaySpan(t.span)
	case *FieldAccess:
		return fmt.Sprintf("(%s%s.%s%s)", displaySpan(t.span), DisplayAst(t.Expr), displaySpan(t.Field.Span), t.Field.Value)
	case *FunctionCall:
		return fmt.Sprintf("(%s%s.%s%s(%s))", displaySpan(t.span), DisplayAst(t.Expr), displaySpan(t.Name.Span), t.Name.Value, DisplayExprs(t.Arguments))
	case *BadFieldAccess:
		return fmt.Sprintf("(%s%s.@)", displaySpan(t.span), DisplayAst(t.Expr))
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
