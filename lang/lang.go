package lang

import "fmt"

type Pos = int

type Span struct {
	Start Pos
	End   Pos
}

func (s Span) Contains(position int) bool {
	return s.Start <= position && position <= s.End
}

func (s Span) replaceIfNil(startSpan Span) Span {
	if s.IsNil() {
		return startSpan
	}
	return s
}

func (s Span) IsNil() bool {
	return s.Start < 0
}

func (s Span) Length() int {
	return s.End - s.Start
}

func (s Span) mergeSpan(a Span) Span {
	if s.IsNil() {
		return a
	}
	if a.IsNil() {
		return s
	}
	return Span{Start: min(s.Start, a.Start), End: max(s.End, a.End)}
}

// TODO pas moyen de trouver les fonction min/max sur des ints
func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}

type Reporter interface {
	Report(span Span, message string)
}

type Error struct {
	span    Span
	message string
}

type DefaultReporter struct {
	errors []Error
}

func (reporter *DefaultReporter) HasErrors() bool {
	return len(reporter.errors) > 0
}

func (reporter *DefaultReporter) Report(span Span, message string) {
	reporter.errors = append(reporter.errors, Error{span: span, message: message})
}

func (reporter DefaultReporter) DisplayErrors(content string) {
	for _, error := range reporter.errors {
		fmt.Println(content)
		for i := 0; i < error.span.Start; i += 1 {
			fmt.Print(" ")
		}
		for i := 0; i < max(error.span.Length(), 1); i += 1 {
			fmt.Print("^")
		}
		fmt.Println()
		fmt.Printf("%d:%d: %s\n", error.span.Start, error.span.End, error.message)
	}
}
