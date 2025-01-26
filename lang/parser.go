package lang

import "fmt"

type Parser struct {
	reporter Reporter
	lexer    *Lexer
	token    *Token
}

func NewParser(lexer *Lexer, reporter Reporter) *Parser {
	var token = lexer.NextToken()
	return &Parser{reporter: reporter, lexer: lexer, token: token}
}

func (p *Parser) Parse() Expr {
	var expr = p.parseAtom()
	for {
		if p.token.kind == DOT {
			p.nextToken()
			var ident = p.expectIdent()
			expr = &FieldAccess{span: expr.Span().mergeSpan(ident.span), expr: expr, field: *ident}
		} else if p.token.kind == LEFT_BRT {
			p.nextToken()
			var index = p.parseAtom()
			if p.token.kind != RIGHT_BRT {
				panic(fmt.Sprintf("expecting ']', got %d", p.token.kind))
			}
			var endSpan = p.token.span
			p.nextToken()
			expr = &ArrayAccess{span: expr.Span().mergeSpan(endSpan), expr: expr, index: index}
		} else {
			break
		}
	}
	if p.token.kind != EOF {
		panic(fmt.Sprintf("Junk at end %+v", p.token))
	}
	return expr
}

func (p *Parser) parseAtom() Expr {
	if p.token.kind == INTEGER {
		var e = &IntegerValue{span: p.token.span, value: p.token.value}
		p.nextToken()
		return e
	}
	return &Start{span: Span{start: p.token.span.start, end: p.token.span.start}}
}

func (p *Parser) expectIdent() *Token {
	if p.token.kind == IDENT {
		var t = p.token
		p.nextToken()
		return t
	}
	panic(fmt.Sprintf("expecting ident, got %d", p.token.kind))
}

func (p *Parser) nextToken() {
	p.token = p.lexer.NextToken()
}
