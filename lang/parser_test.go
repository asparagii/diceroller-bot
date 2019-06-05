package lang

import (
	"math/rand"
	"strings"
	"testing"
)

func TestAtomSingleNumber(t *testing.T) {
	channel := make(chan Token)
	proxy := &Nexter{token: channel}

	go func() {
		channel <- Token{t: NUMBER, value: "12"}
		channel <- Token{t: EOF}
	}()

	proxy.Init()
	val, descr, err := atom(proxy)
	assert(t, err == nil, "Expected no error, got %v", err)
	assert(t, val == 12, "Expected 12, got %v", val)
	assert(t, strings.Compare("12", descr) == 0, "Expected '12', got %v", descr)
}

func TestAtomDice(t *testing.T) {
	rand.Seed(int64(0))
	channel := make(chan Token)
	proxy := &Nexter{token: channel}

	go func() {
		channel <- Token{t: NUMBER, value: "10"}
		channel <- Token{t: DICE}
		channel <- Token{t: NUMBER, value: "8"}
		channel <- Token{t: EOF}
	}()

	proxy.Init()
	val, descr, err := atom(proxy)
	assert(t, err == nil, "Expected no error, got '%v'", err)
	assert(t, strings.Count(descr, "+") == 9, "Expected 9 `+` runes, got %v", strings.Count(descr, "+"))
	assert(t, descr[0] == '(' && descr[len(descr)-1] == ')', "Expected `(` and `)`, found `%v` and `%v`", descr[0], descr[len(descr)-1])
	assert(t, val > 0, "Expected `val` to be greater than zero, got '%v'", val)
}

func TestTerm(t *testing.T) {
	rand.Seed(int64(0))
	channel := make(chan Token)
	proxy := &Nexter{token: channel}

	go func() {
		channel <- Token{t: NUMBER, value: "6"}
		channel <- Token{t: MUL}
		channel <- Token{t: NUMBER, value: "7"}
		channel <- Token{t: EOF}
	}()

	proxy.Init()
	val, _, err := term(proxy)
	assert(t, err == nil, "Expected no error, got '%v'", err)
	assert(t, val == 42, "Expected `val` to equal %v, got %v", 42, val)
}

func TestExprPrecedence(t *testing.T) {
	channel := make(chan Token)
	proxy := &Nexter{token: channel}

	go func() {
		channel <- Token{t: NUMBER, value: "12"}
		channel <- Token{t: SUM}
		channel <- Token{t: NUMBER, value: "4"}
		channel <- Token{t: MUL}
		channel <- Token{t: NUMBER, value: "9"}
		channel <- Token{t: EOF}
	}()

	proxy.Init()
	val, _, err := expr(proxy)
	assert(t, err == nil, "Expected no error, got '%v'", err)
	assert(t, val == 48, "Expected `val` to equal 48, got %v", val)
}

func TestParser(t *testing.T) {
	input := "13+5*2-2"
	val, _, err := Parse(input)

	assert(t, err == nil, "Expected no error, got '%v'", err)
	assert(t, val == 21, "Expected `val` to equal 21, got %v", val)
}

func TestParserUnexpectedToken(t *testing.T) {
	input := "13*+"
	_, _, err := Parse(input)
	assert(t, err != nil, "Expected error, got no error")
}
