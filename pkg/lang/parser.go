package lang

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

func Parse(expression string) (Object, error) {
	token := Lex(expression)

	proxy := Nexter{token: token}

	// Read first value
	proxy.Init()
	value, err := expr(&proxy)
	if proxy.Peek().t != EOF {
		return value, fmt.Errorf("Unexpected end of message")
	} else {
		return value, err
	}
}

func expr(proxy *Nexter) (Object, error) {
	lvalue, err := term(proxy)
	if err != nil {
		return lvalue, err
	}

	for true {
		switch proxy.Peek().t {
		case SUM:
			proxy.Pop()
			rvalue, err := term(proxy)
			if err != nil {
				return lvalue, err
			}
			lvalue, err = Add(lvalue, rvalue)
			if err != nil {
				return lvalue, err
			}
		case SUB:
			proxy.Pop()
			rvalue, err := term(proxy)
			if err != nil {
				return lvalue, err
			}
			lvalue, err = Subtract(lvalue, rvalue)
			if err != nil {
				return lvalue, err
			}
		default:
			return lvalue, nil
		}
	}
	return lvalue, fmt.Errorf("Unexpected end of the parser")
}

func term(proxy *Nexter) (Object, error) {
	lvalue, err := atom(proxy)

	if err != nil {
		return lvalue, err
	}

	for true {
		switch proxy.Peek().t {
		case MUL:
			proxy.Pop()
			rvalue, err := atom(proxy)
			if err != nil {
				return lvalue, err
			}
			lvalue, err = Multiply(lvalue, rvalue)
			if err != nil {
				return lvalue, err
			}
		default:
			return lvalue, nil
		}
	}
	return lvalue, fmt.Errorf("Unexpected end of the parser")
}

func atom(proxy *Nexter) (Object, error) {
	rvalue, err := literal(proxy)
	if err != nil {
		return rvalue, err
	} else if rvalue.Type != NUMBERVALUE {
		return rvalue, fmt.Errorf("Expected number, got %v", rvalue.Type)
	}

	if proxy.Peek().t != DICE {
		return rvalue, nil
	}

	proxy.Pop()
	lvalue, err := literal(proxy)
	if err != nil {
		return lvalue, err
	} else if rvalue.Type != NUMBERVALUE {
		return rvalue, fmt.Errorf("Expected number, got %v", rvalue.Type)
	}
	if proxy.Peek().t != KEEP {
		return roll(rvalue.Value.(int), lvalue.Value.(int), rvalue.Value.(int))
	}
	proxy.Pop()
	kvalue, err := literal(proxy)
	if err != nil {
		return kvalue, err
	} else if kvalue.Type != NUMBERVALUE {
		return kvalue, fmt.Errorf("Expected number, got %v", kvalue.Type)
	}

	return roll(rvalue.Value.(int), lvalue.Value.(int), kvalue.Value.(int))
}

func literal(proxy *Nexter) (Object, error) {
	switch proxy.Peek().t {
	case LPAREN:
		proxy.Pop()
		val, err := expr(proxy)
		if err != nil {
			return val, err
		}
		if proxy.Peek().t != RPAREN {
			return val, fmt.Errorf("Expected ')', found '%s'", proxy.Peek().value)
		}
		proxy.Pop()
		if val.Type != NUMBERVALUE {
			return val, fmt.Errorf("Expected number, found %v", val.Type)
		}
		return val, nil
	case NUMBER:
		val, err := number(proxy)
		return val, err
	default:
		return Object{}, fmt.Errorf("Expected number or expression, found '%s'", proxy.Peek().t)
	}

}

func number(proxy *Nexter) (Object, error) {
	if proxy.Peek().t != NUMBER {
		return Object{}, fmt.Errorf("Expected number, got '%v'", proxy.Peek().value)
	}
	value, err := strconv.Atoi(proxy.Pop().value)
	if err != nil {
		return Object{}, err
	}
	return Number(value), nil
}

func roll(number, size, keep int) (Object, error) {
	if size < 1 {
		return Object{}, fmt.Errorf("Unexpected size '%v'", size)
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
	lowerBound := number - keep
	if lowerBound < 0 {
		lowerBound = 0
	}
	for _, v := range results[lowerBound:number] {
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
	_ = fmt.Sprintf("(%s)", repr)
	return Number(sum), nil
}
