package lang

type Expr interface {
	Span() Span
}

type BadExpr struct {
	span Span
}

func (s *BadExpr) Span() Span {
	return s.span
}

type Start struct {
	span Span
}

func (s *Start) Span() Span {
	return s.span
}

type FieldAccess struct {
	span  Span
	Expr  Expr
	Field Identifier
}

func (s *FieldAccess) Span() Span {
	return s.span
}

type BadFieldAccess struct {
	span Span
	expr Expr
}

func (s *BadFieldAccess) Span() Span {
	return s.span
}

type ArrayAccess struct {
	span  Span
	Expr  Expr
	Index Expr
}

func (s *ArrayAccess) Span() Span {
	return s.span
}

type IntegerValue struct {
	span  Span
	Value string
}

func (s *IntegerValue) Span() Span {
	return s.span
}

type Identifier struct {
	span  Span
	Value string
}
