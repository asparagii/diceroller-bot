package lang

import (
	"fmt"
	"strings"
	"testing"
)

func TestSimple_Expression(t *testing.T) {
	channel := Lex("1+1")
	Expect(t, <-channel, Token{t: NUMBER, value: "1"})
	Expect(t, <-channel, Token{t: SUM})
	Expect(t, <-channel, Token{t: NUMBER, value: "1"})
	Expect(t, <-channel, Token{t: EOF})
}

func TestMultiplication_Expression(t *testing.T) {
	channel := Lex("2*12-9")
	Expect(t, <-channel, Token{t: NUMBER, value: "2"})
	Expect(t, <-channel, Token{t: MUL})
	Expect(t, <-channel, Token{t: NUMBER, value: "12"})
	Expect(t, <-channel, Token{t: SUB})
	Expect(t, <-channel, Token{t: NUMBER, value: "9"})
	Expect(t, <-channel, Token{t: EOF})
}

func TestDice_Expression(t *testing.T) {
	channel := Lex("3d12+2")
	Expect(t, <-channel, Token{t: NUMBER, value: "3"})
	Expect(t, <-channel, Token{t: DICE})
	Expect(t, <-channel, Token{t: NUMBER, value: "12"})
	Expect(t, <-channel, Token{t: SUM})
	Expect(t, <-channel, Token{t: NUMBER, value: "2"})
	Expect(t, <-channel, Token{t: EOF})
}

func TestIgnoreWhitespace(t *testing.T) {
	channel := Lex("1d 1 + 2")
	Expect(t, <-channel, Token{t: NUMBER, value: "1"})
	Expect(t, <-channel, Token{t: DICE})
	Expect(t, <-channel, Token{t: NUMBER, value: "1"})
	Expect(t, <-channel, Token{t: SUM})
	Expect(t, <-channel, Token{t: NUMBER, value: "2"})
	Expect(t, <-channel, Token{t: EOF})
}

func TestIgnoreWhitespaceInArray(t *testing.T) {
	channel := Lex("[1,2, 3]")
	Expect(t, <-channel, Token{t: LARRAYPAREN})
	Expect(t, <-channel, Token{t: NUMBER, value: "1"})
	Expect(t, <-channel, Token{t: LISTSEPARATOR})
	Expect(t, <-channel, Token{t: NUMBER, value: "2"})
	Expect(t, <-channel, Token{t: LISTSEPARATOR})
	Expect(t, <-channel, Token{t: NUMBER, value: "3"})
	Expect(t, <-channel, Token{t: RARRAYPAREN})
	Expect(t, <-channel, Token{t: EOF})
}

func TestPipe(t *testing.T) {
	channel := Lex("[1,2,3]|$+2")
	Expect(t, <-channel, Token{t: LARRAYPAREN})
	Expect(t, <-channel, Token{t: NUMBER, value: "1"})
	Expect(t, <-channel, Token{t: LISTSEPARATOR})
	Expect(t, <-channel, Token{t: NUMBER, value: "2"})
	Expect(t, <-channel, Token{t: LISTSEPARATOR})
	Expect(t, <-channel, Token{t: NUMBER, value: "3"})
	Expect(t, <-channel, Token{t: RARRAYPAREN})
	Expect(t, <-channel, Token{t: PIPE})
	Expect(t, <-channel, Token{t: DOLLAR})
	Expect(t, <-channel, Token{t: SUM})
	Expect(t, <-channel, Token{t: NUMBER, value: "2"})
}

func Expect(t *testing.T, received, expected Token) {
	assert(t, received.t == expected.t, fmt.Sprintf("Token type did not match. Expected %v, found %v", expected.t, received.t))
	assert(t, strings.Compare(expected.value, received.value) == 0, fmt.Sprintf("Token type did not match. Expected %v, found %v", expected.value, received.value))
}
