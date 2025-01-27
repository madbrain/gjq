package lang

import (
	"errors"
	"slices"
)

type Parser struct {
	reporter Reporter
	lexer    *Lexer
	token    *Token
}

func NewParser(lexer *Lexer, reporter Reporter) *Parser {
	token := lexer.NextToken()
	return &Parser{reporter: reporter, lexer: lexer, token: token}
}

/*
Path ::= Atom ( '.' IDENT | '[' Atom ']' )*
*/
func (p *Parser) Parse() Expr {
	expr := p.parseRoot()
	for {
		if p.token.kind == DOT {
			dotSpan := p.token.span
			p.nextToken()
			if ident, err := p.expectIdent(); err == nil {
				expr = &FieldAccess{span: expr.Span().mergeSpan(ident.span), expr: expr, field: *ident}
			} else {
				span := p.skipTo([]TokenKind{DOT, LEFT_BRT, EOF}).mergeSpan(dotSpan)
				expr = &BadFieldAccess{span: expr.Span().mergeSpan(span), expr: expr}
			}
		} else if p.token.kind == LEFT_BRT {
			p.nextToken()
			index := p.parseAtom([]TokenKind{DOT, RIGHT_BRT, EOF})
			if p.token.kind != RIGHT_BRT {
				p.reporter.Report(p.token.span, "expecting ']'")
				span := p.skipTo([]TokenKind{DOT, LEFT_BRT, RIGHT_BRT, EOF}).mergeSpan(index.Span())
				expr = &ArrayAccess{span: expr.Span().mergeSpan(span), expr: expr, index: index}
			} else {
				endSpan := p.token.span
				p.nextToken()
				expr = &ArrayAccess{span: expr.Span().mergeSpan(endSpan), expr: expr, index: index}
			}
		} else if p.token.kind == EOF {
			break
		} else {
			span := p.skipTo([]TokenKind{DOT, LEFT_BRT, EOF})
			p.reporter.Report(span, "unexpected tokens")
		}
	}
	return expr
}

func (p *Parser) parseAtom(syncTokens []TokenKind) Expr {
	startSpan := Span{start: p.token.span.start, end: p.token.span.start}
	if p.token.kind == INTEGER {
		e := &IntegerValue{span: p.token.span, value: p.token.value}
		p.nextToken()
		return e
	}
	p.reporter.Report(p.token.span, "expecting integer")
	span := p.skipTo(syncTokens).replaceIfNil(startSpan)
	return &BadExpr{span: span}
}

func (p *Parser) parseRoot() Expr {
	return &Start{span: Span{start: p.token.span.start, end: p.token.span.start}}
}

func (p *Parser) expectIdent() (*Identifier, error) {
	if p.token.kind == IDENT {
		t := p.token
		p.nextToken()
		return &Identifier{span: t.span, value: t.value}, nil
	}
	p.reporter.Report(p.token.span, "expecting ident")
	return nil, errors.New("")
}

func (p *Parser) skipTo(expectedTokens []TokenKind) Span {
	span := Span{start: -1, end: -1}
	for {
		if slices.Contains(expectedTokens, p.token.kind) {
			break
		}
		span = span.mergeSpan(p.token.span)
		p.nextToken()
	}
	return span
}

func (p *Parser) nextToken() {
	p.token = p.lexer.NextToken()
}
