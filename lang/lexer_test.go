package lang

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

func TestSimple_Expression(t *testing.T) {
	channel := Lex("1+1")
	Expect(t, <-channel, Token{t: NUMBER, value: "1"})
	Expect(t, <-channel, Token{t: SUM})
	Expect(t, <-channel, Token{t: NUMBER, value: "1"})
	Expect(t, <-channel, Token{t: EOF})
}

func TextMultiplication_Expression(t *testing.T) {
	channel := Lex("2*12-9")
	Expect(t, <-channel, Token{t: NUMBER, value: "2"})
	Expect(t, <-channel, Token{t: SUM})
	Expect(t, <-channel, Token{t: NUMBER, value: "12"})
	Expect(t, <-channel, Token{t: SUB})
	Expect(t, <-channel, Token{t: NUMBER, value: "9"})
	Expect(t, <-channel, Token{t: EOF})
}

func TextDice_Expression(t *testing.T) {
	channel := Lex("3d12+2")
	Expect(t, <-channel, Token{t: NUMBER, value: "3"})
	Expect(t, <-channel, Token{t: DICE})
	Expect(t, <-channel, Token{t: NUMBER, value: "12"})
	Expect(t, <-channel, Token{t: SUM})
	Expect(t, <-channel, Token{t: NUMBER, value: "2"})
	Expect(t, <-channel, Token{t: EOF})
}

func Expect(t *testing.T, received, expected Token) {
	assert(t, received.t == expected.t, fmt.Sprintf("Token type did not match. Expected %v, found %v", expected.t, received.t))
	assert(t, strings.Compare(expected.value, received.value) == 0, fmt.Sprintf("Token type did not match. Expected %v, found %v", expected.value, received.value))
}
