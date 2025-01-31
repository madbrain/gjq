package lang

type Expr interface {
	Span() Span
	Type() Type
}

type BadExpr struct {
	span Span
}

func (s *BadExpr) Span() Span {
	return s.span
}

func (s *BadExpr) Type() Type {
	return nil
}

type Start struct {
	span Span
	t    Type
}

func (s *Start) Span() Span {
	return s.span
}

func (s *Start) Type() Type {
	return s.t
}

type FieldAccess struct {
	span  Span
	Expr  Expr
	Field Identifier
	t     Type
}

func (s *FieldAccess) Span() Span {
	return s.span
}

func (s *FieldAccess) Type() Type {
	return s.t
}

type BadFieldAccess struct {
	span Span
	Expr Expr
}

func (s *BadFieldAccess) Span() Span {
	return s.span
}

func (s *BadFieldAccess) Type() Type {
	return nil
}

type ArrayAccess struct {
	span  Span
	Expr  Expr
	Index Expr
	t     Type
}

func (s *ArrayAccess) Span() Span {
	return s.span
}

func (s *ArrayAccess) Type() Type {
	return s.t
}

type IntegerValue struct {
	span  Span
	Value string
}

func (s *IntegerValue) Span() Span {
	return s.span
}

func (s *IntegerValue) Type() Type {
	return &NumberType{}
}

type Identifier struct {
	Span  Span
	Value string
}
