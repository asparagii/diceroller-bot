package lang

import (
	"math/rand"
	//"strings"
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
	val, err := atom(proxy)
	assert(t, err == nil, "Expected no error, got %v", err)
	assert(t, val.Type == NUMBERVALUE, "Expected valber, got %v", val.Type)
	assert(t, val.Value.(int) == 12, "Expected 12, got %v", val.Value.(int))
	//assert(t, strings.Compare("12", val.String()) == 0, "Expected '12', got %v", val)
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
	val, err := atom(proxy)
	assert(t, err == nil, "Expected no error, got '%v'", err)
	//assert(t, strings.Count(val.String(), "+") == 9, "Expected 9 `+` runes, got %v", strings.Count(val.String(), "+"))
	//assert(t, val.String()[0] == '(' && val.String()[len(val.String())-1] == ')', "Expected `(` and `)`, found `%v` and `%v`", val.String()[0], val.String()[len(val.String())-1])
	assert(t, val.Type == NUMBERVALUE, "Expected valber, got %v", val.Type)
	assert(t, val.Value.(int) > 0, "Expected `val` to be greater than zero, got '%v'", val)
}

func TestAtomDiceKeep(t *testing.T) {
	rand.Seed(int64(0))
	channel := make(chan Token)
	proxy := &Nexter{token: channel}

	go func() {
		channel <- Token{t: NUMBER, value: "10"}
		channel <- Token{t: DICE}
		channel <- Token{t: NUMBER, value: "8"}
		channel <- Token{t: KEEP}
		channel <- Token{t: NUMBER, value: "4"}
		channel <- Token{t: EOF}
	}()

	proxy.Init()
	val, err := atom(proxy)
	assert(t, err == nil, "Expected no error, got '%v'", err)
	assert(t, val.Type == NUMBERVALUE, "Expected valber, got %v", val.Type)
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
	val, err := term(proxy)
	assert(t, err == nil, "Expected no error, got '%v'", err)
	assert(t, val.Value.(int) == 42, "Expected `val` to equal %v, got %v", 42, val)
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
	val, err := expr(proxy)
	assert(t, err == nil, "Expected no error, got '%v'", err)
	assert(t, val.Type == NUMBERVALUE, "Expected valber, got %v", val.Type)
	assert(t, val.Value.(int) == 48, "Expected `val` to equal 48, got %v", val.Value.(int))
}

func TestParser(t *testing.T) {
	input := "13+5*2-2"
	val, err := Parse(input)

	assert(t, err == nil, "Expected no error, got '%v'", err)
	assert(t, val.Type == NUMBERVALUE, "Expected valber, got %v", val.Type)
	assert(t, val.Value.(int) == 21, "Expected `val` to equal 21, got %v", val.Value.(int))
}

func TestParserUnexpectedToken(t *testing.T) {
	input := "13*+"
	_, err := Parse(input)
	assert(t, err != nil, "Expected error, got no error")
}

func TestParserKeepDice(t *testing.T) {
	input := "4d8k3"
	val, err := Parse(input)
	assert(t, err == nil, "Expected no error, got '%v'", err)
	assert(t, val.Type == NUMBERVALUE, "Expected valber, got %v", val.Type)
	assert(t, val.Value.(int) > 0, "Expected value to be greater than 0, but it was not")
}
