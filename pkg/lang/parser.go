package lang

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

func Parse(expression string) (int, string, error) {
	token := Lex(expression)

	proxy := Nexter{token: token}

	// Read first value
	proxy.Init()
	value, description, err := expr(&proxy)
	if proxy.Peek().t != EOF {
		return 0, "", fmt.Errorf("Unexpected end of message")
	} else {
		return value, description, err
	}
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
		default:
			return lvalue, ldescr, nil
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
	if proxy.Peek().t == LPAREN {
		proxy.Pop()
		val, descr, err := expr(proxy)
		if err != nil {
			return 0, "", err
		}
		if proxy.Peek().t != RPAREN {
			return 0, "", fmt.Errorf("Expected ')', found '%s'", proxy.Peek().value)
		}
		proxy.Pop()
		return val, fmt.Sprintf("(%s)", descr), nil
	} else {
		var size *int = nil
		if proxy.Peek().t == NUMBER {
			val, _, err := number(proxy)
			if err != nil {
				return 0, "", err
			}
			size = new(int)
			*size = val
		}

		if proxy.Peek().t != DICE {
			if size == nil {
				return 0, "", fmt.Errorf("Expected number or dice expression, got '%v'", proxy.Peek().value)
			} else {
				return *size, fmt.Sprintf("%v", *size), nil
			}
		}

		if size == nil {
			size = new(int)
			*size = 1
		}

		number, _, err := dice(proxy)
		if err != nil {
			return 0, "", err
		}

		if proxy.Peek().t != KEEP {
			return roll(*size, number, number, false)
		}

		keep, _, err := keep(proxy)
		if err != nil {
			return 0, "", err
		}

		return roll(*size, number, keep, false)
	}
}

func number(proxy *Nexter) (int, string, error) {
	if proxy.Peek().t != NUMBER {
		return 0, "", fmt.Errorf("Expected number, got '%v'", proxy.Peek().value)
	}
	value, err := strconv.Atoi(proxy.Pop().value)
	if err != nil {
		return 0, "", err
	}
	return value, fmt.Sprintf("%v", value), nil
}

func dice(proxy *Nexter) (int, string, error) {
	if proxy.Peek().t != DICE {
		return 0, "", fmt.Errorf("Expected 'd', got '%v'", proxy.Peek().value)
	}
	proxy.Pop()
	if proxy.Peek().t != NUMBER {
		return 0, "", fmt.Errorf("Expected number, found %s", proxy.Peek().value)
	}
	rtoken := proxy.Pop()
	rvalue, err := strconv.Atoi(rtoken.value)
	if err != nil {
		return 0, "", fmt.Errorf("Error while parsing a number")
	}
	return rvalue, "", nil
}

func keep(proxy *Nexter) (int, string, error) {
	if proxy.Peek().t != KEEP {
		return 0, "", fmt.Errorf("Expected k, found '%s'", proxy.Peek().value)
	}
	proxy.Pop()
	if proxy.Peek().t != NUMBER {
		return 0, "", fmt.Errorf("Expected number, found %s", proxy.Peek().value)
	}
	ktoken := proxy.Pop()
	kvalue, err := strconv.Atoi(ktoken.value)
	if err != nil {
		return 0, "", fmt.Errorf("Error while parsing keep number '%s'", ktoken.value)
	}
	return kvalue, "", nil
}

// Test this shit
func roll(number, size, keep int, formatted bool) (int, string, error) {
	if size < 1 {
		return 0, "", fmt.Errorf("Unexpected size '%v'", size)
	}

	results := make([]int, number)
	for i := 0; i < number; i++ {
		tmp := rand.Intn(size) + 1
		results[i] = tmp
	}
	sorted := sort.IntSlice(results)
	sort.Sort(sorted)

	repr := ""
	if number-keep > 0 {
		for _, v := range results[:number-keep] {
			if len(repr) > 0 {
				repr = fmt.Sprintf("%s+%v", repr, v)
			} else {
				repr = fmt.Sprintf("%v", v)
			}
		}
		repr = fmt.Sprintf("~~%s~~", repr)
	}
	sum := 0
	for _, v := range results[number-keep : number] {
		sum += v
		singleRepr := fmt.Sprintf("%v", v)
		if v == size {
			singleRepr = fmt.Sprintf("**%s**", singleRepr)
		}

		if len(repr) > 0 {
			repr = fmt.Sprintf("%s+%s", repr, singleRepr)
		} else {
			repr = fmt.Sprintf("%s", singleRepr)
		}
	}

	return sum, fmt.Sprintf("(%s)", repr), nil
}
