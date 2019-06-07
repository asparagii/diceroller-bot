package lang

import (
	"strings"
	"unicode"
)

type Token struct {
	t     TokenType
	value string
}

type TokenType int

const (
	NUMBER TokenType = iota
	EOF
	SUM
	SUB
	MUL
	DICE
	KEEP
	RPAREN
	LPAREN
	INVALID
)

func (t TokenType) String() string {
	switch t {
	case NUMBER:
		return "NUMBER"
	case SUM:
		return "SUM"
	case SUB:
		return "SUB"
	case DICE:
		return "DICE"
	case INVALID:
		return "INVALID"
	case MUL:
		return "MUL"
	case EOF:
		return "EOF"
	case KEEP:
		return "KEEP"
	case RPAREN:
		return "RPAREN"
	case LPAREN:
		return "LPAREN"
	default:
		return "<ERROR: Unexpected token type>"
	}
}

type Lexer struct {
	input string
	start int
	pos   int
	out   chan Token
}

func (l *Lexer) Next() rune {
	l.pos++
	if len(l.input) < l.pos {
		return rune(0)
	}
	return []rune(l.input)[l.pos-1]
}

func (l *Lexer) Back() {
	l.pos--
}

func (l *Lexer) Accept(accepted string) {
	for ch := l.Next(); strings.ContainsRune(accepted, ch); {
		ch = l.Next()
	}
	l.Back()
}

func (l *Lexer) Current() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) Emit(t Token) {
	l.out <- t
}

func (l *Lexer) Pop() string {
	ret := l.Current()
	l.start = l.pos
	return ret
}

func Lex(input string) chan Token {
	out := make(chan Token, 10)
	l := &Lexer{
		input: input,
		start: 0,
		pos:   0,
		out:   out,
	}

	go run(l)
	return out
}

type StateFn func(l *Lexer) StateFn

func lexNumber(l *Lexer) StateFn {
	l.Accept("0123456789")
	token := Token{
		t:     NUMBER,
		value: l.Pop(),
	}
	l.Emit(token)
	return lexNeutral
}

func lexNeutral(l *Lexer) StateFn {
	ch := l.Next()
	if ch == rune(0) {
		l.Emit(Token{t: EOF})
		return nil
	}
	switch true {
	case unicode.IsSpace(ch):
		l.IgnoreWhitespace()
		return lexNeutral
	case unicode.IsNumber(ch):
		return lexNumber
	case ch == '+':
		popEmitSingle(SUM, l)
		return lexNeutral
	case ch == '-':
		popEmitSingle(SUB, l)
		return lexNeutral
	case ch == '*':
		popEmitSingle(MUL, l)
		return lexNeutral
	case ch == 'd':
		popEmitSingle(DICE, l)
		return lexNeutral
	case ch == 'k':
		popEmitSingle(KEEP, l)
		return lexNeutral
	case ch == '(':
		popEmitSingle(LPAREN, l)
		return lexNeutral
	case ch == ')':
		popEmitSingle(RPAREN, l)
		return lexNeutral
	default:
		return lexUnexpected
	}
}

func lexUnexpected(l *Lexer) StateFn {
	l.Emit(Token{t: INVALID, value: l.Pop()})
	return nil
}

func (l *Lexer) IgnoreWhitespace() {
	for ch := l.Next(); unicode.IsSpace(ch); {
		ch = l.Next()
	}
	l.Back()
	l.Pop()
}

func popEmitSingle(t TokenType, l *Lexer) {
	l.Pop()
	token := Token{
		t: t,
	}
	l.Emit(token)
}

func run(l *Lexer) {
	for state := lexNeutral; state != nil; {
		state = state(l)
	}
}
