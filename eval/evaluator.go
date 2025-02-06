package eval

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"com.github/madbrain/gjq/lang"
)

type Evaluator struct {
	jsonContent any
	schema      lang.Type
}

type Context struct {
	root any
}

func (e *Evaluator) ReadJsonFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var parsed any
	err = json.Unmarshal(content, &parsed)
	if err != nil {
		return err // malformed input
	}
	e.jsonContent = parsed
	e.schema = InferSchema(parsed) // TODO store schema
	return nil
}

func (e *Evaluator) Evaluate(expr string) {
	reporter := &lang.DefaultReporter{}
	lexer := lang.NewLexer(expr, reporter)
	parser := lang.NewParser(lexer, reporter)
	ast := parser.Parse()

	if reporter.HasErrors() {
		reporter.DisplayErrors(expr)
		return
	}

	//fmt.Println(lang.DisplayAst(ast))

	result := e.evaluateExpr(ast, &Context{root: e.jsonContent})

	PrettyPrintJson(result)
}

func (e *Evaluator) evaluateExpr(ast lang.Expr, context *Context) any {
	switch a := ast.(type) {
	case *lang.Start:
		return context.root
	case *lang.IntegerValue:
		if value, err := strconv.Atoi(a.Value); err == nil {
			return value
		}
		return nil
	case *lang.FieldAccess:
		return getField(e.evaluateExpr(a.Expr, context), a.Field.Value)
	case *lang.ArrayAccess:
		if index, err := asInteger(e.evaluateExpr(a.Index, context)); err == nil {
			return getElement(e.evaluateExpr(a.Expr, context), index)
		}
		return nil
	case *lang.FunctionCall:
		var args []any = nil
		for _, arg := range a.Arguments {
			args = append(args, e.evaluateExpr(arg, context))
		}
		return e.evaluateFunc(e.evaluateExpr(a.Expr, context), a.Name.Value, args)
	default:
		return nil
	}
}

func (e *Evaluator) evaluateFunc(target any, name string, args []any) any {
	switch t := target.(type) {
	case map[string]any:
		if name == "keys" && len(args) == 0 {
			keys := make([]any, len(t))
			i := 0
			for k := range t {
				keys[i] = k
				i++
			}
			return keys
		} else {
			return nil
		}
	case []any:
		if name == "length" && len(args) == 0 {
			return len(t)
		} else {
			return nil
		}
	default:
		return nil
	}
}

func asInteger(v any) (int, error) {
	switch t := v.(type) {
	case int:
		return t, nil
	default:
		return 0, errors.New("not an integer")
	}
}

func getField(v any, fieldName string) any {
	switch m := v.(type) {
	case map[string]any:
		return m[fieldName]
	default:
		return nil
	}
}

func getElement(v any, index int) any {
	switch m := v.(type) {
	case []any:
		if index >= 0 && index < len(m) {
			return m[index]
		}
	}
	return nil
}
