package lang

type Type interface {
}

type ErrorType struct{}

type NumberType struct{}

type StringType struct{}

type BooleanType struct{}

type ArrayType struct {
	Element Type
}

type ObjectType struct {
	Fields map[string]Type
}

type InferContext struct {
	RootType Type
}

func InferType(e Expr, context *InferContext) Type {
	switch a := e.(type) {
	case *Start:
		a.t = context.RootType
		return a.t
	case *IntegerValue:
		return &NumberType{}
	case *ArrayAccess:
		InferType(a.Index, context)
		rt := InferType(a.Expr, context)
		switch b := rt.(type) {
		case *ArrayType:
			a.t = b.Element
		}
		return a.t
	case *FieldAccess:
		rt := InferType(a.Expr, context)
		switch b := rt.(type) {
		case *ObjectType:
			a.t = b.Fields[a.Field.Value]
		}
		return a.t
	case *BadFieldAccess:
		InferType(a.Expr, context)
		return nil
	}
	return nil
}
