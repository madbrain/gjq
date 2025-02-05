package eval

import (
	"fmt"

	"com.github/madbrain/gjq/lang"
)

func InferSchema(v any) lang.Type {
	switch t := v.(type) {
	case bool:
		return &lang.BooleanType{}
	case float64:
		return &lang.NumberType{}
	case string:
		return &lang.StringType{}
	case []any:
		result := &lang.ArrayType{}
		if len(t) > 0 {
			result.Element = InferSchema(t[0]) // TODO handle tuple ?
		}
		return result
	case map[string]any:
		result := &lang.ObjectType{Fields: make(map[string]lang.Type)}
		for name, value := range t {
			result.Fields[name] = InferSchema(value)
		}
		return result
	default:
		fmt.Printf("unknown type for %+v\n", t)
		return nil
	}
}
