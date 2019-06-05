package lang

import (
	"fmt"
	"math/rand"
	"strconv"
)

func Parse(expression string) (int, string, error) {
	token := Lex(expression)

	proxy := Nexter{token: token}

	// Read first value
	proxy.Init()
	return expr(&proxy)
}

func expr(proxy *Nexter) (int, string, error) {
	lvalue, ldescr, err := term(proxy)
	if err != nil {
		return lvalue, ldescr, err
	}

	for true {
		switch proxy.Peek().t {
		case SUM:
			proxy.Pop()
			rvalue, rdescr, err := term(proxy)
			if err != nil {
				return lvalue, ldescr, err
			}
			lvalue, ldescr = lvalue+rvalue, fmt.Sprintf("%s+%s", ldescr, rdescr)
		case SUB:
			proxy.Pop()
			rvalue, rdescr, err := term(proxy)
			if err != nil {
				return lvalue, ldescr, err
			}
			lvalue, ldescr = lvalue-rvalue, fmt.Sprintf("%s-%s", ldescr, rdescr)
		case EOF:
			return lvalue, ldescr, nil
		default:
			return lvalue, ldescr, fmt.Errorf("Expected operator or `EOF`, found '%v'", proxy.Peek())
		}
	}
	return lvalue, ldescr, fmt.Errorf("Unexpected end of the parser")
}

func term(proxy *Nexter) (int, string, error) {
	lvalue, ldescr, err := atom(proxy)

	if err != nil {
		return lvalue, ldescr, err
	}

	for true {
		switch proxy.Peek().t {
		case MUL:
			proxy.Pop()
			rvalue, rdescr, err := atom(proxy)
			if err != nil {
				return lvalue, ldescr, err
			}
			lvalue, ldescr = lvalue*rvalue, fmt.Sprintf("%s*%s", ldescr, rdescr)
		default:
			return lvalue, ldescr, nil
		}
	}
	return lvalue, ldescr, fmt.Errorf("Unexpected end of the parser")
}

func atom(proxy *Nexter) (int, string, error) {
	var lvalue int
	var lrepr string
	var err error
	if proxy.Peek().t == NUMBER {
		ltoken := proxy.Pop()
		lvalue, err = strconv.Atoi(ltoken.value)
		if err != nil {
			return 0, "", fmt.Errorf("Error while parsing a number")
		}
		lrepr = fmt.Sprintf("%v", lvalue)
	} else {
		lvalue = 1
		lrepr = "1"
	}

	if proxy.Peek().t != DICE {
		return lvalue, lrepr, nil
	} else {
		proxy.Pop()
		if proxy.Peek().t != NUMBER {
			return lvalue, lrepr, fmt.Errorf("Expected number, found %s", proxy.Peek().value)
		}
		rtoken := proxy.Pop()
		rvalue, err := strconv.Atoi(rtoken.value)
		if err != nil {
			return 0, "", fmt.Errorf("Error while parsing a number")
		}
		result, repr := roll(lvalue, rvalue)
		return result, repr, nil
	}
}

func roll(number, size int) (int, string) {
	sum := 0
	repr := ""
	for i := 0; i < number; i++ {
		tmp := rand.Intn(size) + 1
		sum += tmp
		if len(repr) == 0 {
			repr = fmt.Sprintf("%v", tmp)
		} else {
			repr = fmt.Sprintf("%s+%v", repr, tmp)
		}
	}
	return sum, fmt.Sprintf("(%s)", repr)
}

type Nexter struct {
	token    chan Token
	internal Token
}

func (n *Nexter) Init() {
	n.internal = <-n.token
}

func (n *Nexter) Pop() Token {
	tmp := n.internal
	if n.internal.t != EOF {
		n.internal = <-n.token
	}
	return tmp
}

func (n *Nexter) Peek() Token {
	return n.internal
}
