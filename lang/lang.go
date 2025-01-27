package lang

type Pos = int

type Span struct {
	start Pos
	end   Pos
}

func (s Span) replaceIfNil(startSpan Span) Span {
	if s.IsNil() {
		return startSpan
	}
	return s
}

func (s Span) IsNil() bool {
	return s.start < 0
}

func (s Span) Length() int {
	return s.end - s.start
}

func (s Span) mergeSpan(a Span) Span {
	if s.IsNil() {
		return a
	}
	if a.IsNil() {
		return s
	}
	return Span{start: min(s.start, a.start), end: max(s.end, a.end)}
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
