package lang

type Lexer struct {
	reporter Reporter
	content  string
	index    int
	junkPos  Pos
	junk     string
}

func NewLexer(content string, reporter Reporter) *Lexer {
	return &Lexer{reporter: reporter, content: content, index: 0, junk: ""}
}

func (l *Lexer) NextToken() *Token {
	for {
		var start = l.index
		var c = l.getChar()
		if c == -1 {
			return l.newToken(start, EOF)
		}
		if isSpace(c) {
			continue
		}
		if c == '.' {
			return l.newToken(start, DOT)
		}
		if c == '[' {
			return l.newToken(start, LEFT_BRT)
		}
		if c == ']' {
			return l.newToken(start, RIGHT_BRT)
		}
		if isLetter(c) {
			return l.ident(start, c)
		}
		if isDigit(c) {
			return l.integer(start, c)
		}
		l.addJunk(start, c)
	}
}

func (l *Lexer) ident(start Pos, c int) *Token {
	var result = string(byte(c))
	for {
		c = l.getChar()
		if !isLetter(c) {
			l.unget(c)
			break
		}
		result += string(byte(c))
	}
	return &Token{span: Span{start: start, end: l.index}, kind: IDENT, value: result}
}

func (l *Lexer) integer(start int, c int) *Token {
	var result = string(byte(c))
	for {
		c = l.getChar()
		if !isDigit(c) {
			l.unget(c)
			break
		}
		result += string(byte(c))
	}
	return &Token{span: Span{start: start, end: l.index}, kind: INTEGER, value: result}
}

func isLetter(c int) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isDigit(c int) bool {
	return c >= '0' && c <= '9'
}

func isSpace(c int) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func (l *Lexer) newToken(start int, kind TokenKind) *Token {
	l.flushJunk()
	return &Token{span: Span{start: start, end: l.index}, kind: kind}
}

func (l *Lexer) flushJunk() {
	if l.junk != "" {
		l.reporter.Report(Span{start: l.junkPos, end: l.junkPos + len(l.junk)}, "Unrecognized character(s)")
		l.junk = ""
	}
}

func (l *Lexer) addJunk(start Pos, c int) {
	if l.junk == "" {
		l.junkPos = start
	}
	l.junk += string(byte(c))
}

func (l *Lexer) getChar() int {
	if l.index >= len(l.content) {
		return -1
	}
	var c = l.content[l.index]
	l.index += 1
	return int(c)
}

func (l *Lexer) unget(c int) {
	if c >= 0 {
		l.index -= 1
	}
}
