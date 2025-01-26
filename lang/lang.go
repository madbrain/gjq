package lang

type Pos = int

type Span struct {
	start Pos
	end   Pos
}

func (s Span) mergeSpan(a Span) Span {
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
