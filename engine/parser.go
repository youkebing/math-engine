package engine

import (
	"errors"
	"fmt"
)

const (
	Identifier = iota
	Literal
	Operator
)

type Token struct {
	Tok  string
	Type int
	Flag int

	Offset int
}

type Parser struct {
	Source string
	ch     byte
	offset int

	err error
}

func Parse(s string) ([]*Token, error) {
	p := &Parser{
		Source: s,
		err:    nil,
		ch:     s[0],
	}
	toks := p.parse()
	if p.err != nil {
		return nil, p.err
	}
	return toks, nil
}

func (p *Parser) parse() []*Token {
	toks := make([]*Token, 0)
	for {
		tok := p.nextTok()
		if tok == nil {
			break
		}
		toks = append(toks, tok)
	}
	return toks
}

func (p *Parser) nextTok() *Token {
	if p.offset >= len(p.Source) || p.err != nil {
		return nil
	}
	var err error
	for p.isWhitespace(p.ch) && err == nil {
		err = p.nextCh()
	}
	start := p.offset
	var tok *Token
	switch p.ch {
	case
		'(',
		')',
		'+',
		'-',
		'*',
		'/',
		'^',
		'%':
		tok = &Token{
			Tok:  string(p.ch),
			Type: Operator,
		}
		tok.Offset = start
		err = p.nextCh()

	case
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9':
		for p.isDigitNum(p.ch) && p.nextCh() == nil {
		}
		tok = &Token{
			Tok:  p.Source[start:p.offset],
			Type: Literal,
		}
		tok.Offset = start

	default:
		if p.ch != ' ' {
			s := fmt.Sprintf("symbol error: unkown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
		}
	}
	return tok
}

func (p *Parser) nextCh() error {
	p.offset++
	if p.offset < len(p.Source) {
		p.ch = p.Source[p.offset]
		return nil
	}
	return errors.New("EOF")
}

func (p *Parser) isWhitespace(c byte) bool {
	return c == ' ' ||
		c == '\t' ||
		c == '\n' ||
		c == '\v' ||
		c == '\f' ||
		c == '\r'
}

func (p *Parser) isDigitNum(c byte) bool {
	return '0' <= c && c <= '9' || c == '.'
}