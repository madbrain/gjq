package eval

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"com.github/madbrain/gjq/lang"
	rl "github.com/nyaosorg/go-readline-ny"
)

type ExprCompletion struct {
}

func (C ExprCompletion) String() string {
	return "EXPR_COMPLETION"
}

type NullReporter struct{}

func (r *NullReporter) Report(span lang.Span, message string) {
	// ignore
}

type Candidate struct {
	replace int
	value   string
}

type CompletionResult struct {
	candidates      []Candidate
	alreadyComplete bool
}

var objectFunctions = [...]string{"keys"}
var arrayFunctions = [...]string{"length"}

func completeFunction(candidates []Candidate, prefix string, addDot bool, functions []string) []Candidate {
	for _, n := range functions {
		if len(n) >= len(prefix) && n[0:len(prefix)] == prefix {
			value := n + "()"
			if addDot {
				value = "." + value
			}
			candidates = append(candidates, Candidate{replace: len(prefix), value: value})
		}
	}
	return candidates
}

func completeObject(candidates []Candidate, t *lang.ObjectType, prefix string, addDot bool) CompletionResult {
	for n := range t.Fields {
		if len(n) >= len(prefix) && n[0:len(prefix)] == prefix {
			if n == prefix {
				return CompletionResult{alreadyComplete: true}
			}
			value := n
			if addDot {
				value = "." + value
			}
			candidates = append(candidates, Candidate{replace: len(prefix), value: value})
		}
	}
	candidates = completeFunction(candidates, prefix, addDot, objectFunctions[:])
	return CompletionResult{candidates: candidates, alreadyComplete: false}
}

func findCompletionsAt(ast lang.Expr, position int) []Candidate {
	var candidates []Candidate
	if ast.Span().Contains(position) {
		switch a := ast.(type) {
		case *lang.Start:
			switch t := a.Type().(type) {
			case *lang.ObjectType:
				result := completeObject(candidates, t, "", true)
				candidates = result.candidates
			}
		case *lang.FieldAccess:
			if a.Field.Span.Contains(position) {
				prefix := a.Field.Value[0 : position-a.Field.Span.Start]
				switch t := a.Expr.Type().(type) {
				case *lang.ObjectType:
					result := completeObject(candidates, t, prefix, false)
					if result.alreadyComplete {
						switch tt := a.Type().(type) {
						case *lang.ObjectType:
							result = completeObject(candidates, tt, "", true)
							candidates = result.candidates
						case *lang.ArrayType:
							candidates = completeFunction(candidates, "", true, arrayFunctions[:])
						}
					} else {
						candidates = result.candidates
					}
				}
			} else if a.Expr.Span().Contains(position) {
				return findCompletionsAt(a.Expr, position)
			}
		case *lang.BadFieldAccess:
			switch t := a.Expr.Type().(type) {
			case *lang.ObjectType:
				result := completeObject(candidates, t, "", false)
				candidates = result.candidates
			}
		}
	} else if position > ast.Span().End {
		switch tt := ast.Type().(type) {
		case *lang.ObjectType:
			result := completeObject(candidates, tt, "", true)
			candidates = result.candidates

		}
	}
	return candidates
}

// copied from go-readline-ny
func commonPrefix(list []string) string {
	if len(list) < 1 {
		return ""
	}
	common := list[0]
	var cr, fr *strings.Reader
	// to make English case of completed word to the shortest candidate
	minimumLength := len(list[0])
	minimumIndex := 0
	for index, f := range list[1:] {
		if cr != nil {
			cr.Reset(common)
		} else {
			cr = strings.NewReader(common)
		}
		if fr != nil {
			fr.Reset(f)
		} else {
			fr = strings.NewReader(f)
		}
		i := 0
		var buffer strings.Builder
		for {
			ch, _, err := cr.ReadRune()
			if err != nil {
				break
			}
			fh, _, err := fr.ReadRune()
			if err != nil || unicode.ToUpper(ch) != unicode.ToUpper(fh) {
				break
			}
			buffer.WriteRune(ch)
			i++
		}
		common = buffer.String()
		if len(f) < minimumLength {
			minimumLength = len(f)
			minimumIndex = index + 1
		}
	}
	return list[minimumIndex][:len(common)]
}

func complete(B *rl.Buffer, e *Evaluator) {
	reporter := &NullReporter{}
	lexer := lang.NewLexer(B.String(), reporter)
	parser := lang.NewParser(lexer, reporter)
	ast := parser.Parse()
	lang.InferType(ast, &lang.InferContext{RootType: e.schema})

	// fmt.Printf("\n%s @ %d\n", lang.DisplayAst(ast), B.Cursor)
	candidates := findCompletionsAt(ast, B.Cursor)

	if len(candidates) == 1 {
		candidate := candidates[0]
		B.ReplaceAndRepaint(B.Cursor-candidate.replace, candidate.value)
	} else if len(candidates) > 1 {
		list := make([]string, len(candidates))
		for i, candidate := range candidates {
			list[i] = candidate.value
		}
		prefix := commonPrefix(list)
		B.ReplaceAndRepaint(B.Cursor-candidates[0].replace, prefix)
		B.Out.WriteByte('\n')
		for _, candidate := range candidates {
			B.Out.WriteString(fmt.Sprintf("%s\n", candidate.value))
		}
		B.RepaintAll()
	}
}

func (C ExprCompletion) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	e, ok := ctx.Value("evaluator").(*Evaluator)
	if !ok {
		return rl.ABORT
	}
	complete(B, e)
	return rl.CONTINUE
}
