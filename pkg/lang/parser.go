package lang

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

func Parse(expression string) (Value, error) {
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

func expr(proxy *Nexter) (Value, error) {
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
			lvalue, err = lvalue.Add(rvalue)
			if err != nil {
				return lvalue, err
			}
		case SUB:
			proxy.Pop()
			rvalue, err := term(proxy)
			if err != nil {
				return lvalue, err
			}
			invertedvalue, err := rvalue.Invert()
			if err != nil {
				return invertedvalue, err
			}
			lvalue, err = lvalue.Add(invertedvalue)
			if err != nil {
				return lvalue, err
			}
		default:
			return lvalue, nil
		}
	}
	return lvalue, fmt.Errorf("Unexpected end of the parser")
}

func term(proxy *Nexter) (Value, error) {
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
			lvalue, err = lvalue.Mul(rvalue)
			if err != nil {
				return lvalue, err
			}
		default:
			return lvalue, nil
		}
	}
	return lvalue, fmt.Errorf("Unexpected end of the parser")
}

func atom(proxy *Nexter) (Value, error) {
	rvalue, err := literal(proxy)
	if err != nil {
		return rvalue, err
	} else if rvalue.Type() != NUMBERVALUE {
		return rvalue, fmt.Errorf("Expected number, got %v", rvalue.Type())
	}

	if proxy.Peek().t != DICE {
		return rvalue, nil
	}

	proxy.Pop()
	lvalue, err := literal(proxy)
	if err != nil {
		return lvalue, err
	} else if rvalue.Type() != NUMBERVALUE {
		return rvalue, fmt.Errorf("Expected number, got %v", rvalue.Type())
	}
	if proxy.Peek().t != KEEP {
		return roll(rvalue.(Number), lvalue.(Number), rvalue.(Number))
	}
	proxy.Pop()
	kvalue, err := literal(proxy)
	if err != nil {
		return kvalue, err
	} else if kvalue.Type() != NUMBERVALUE {
		return kvalue, fmt.Errorf("Expected number, got %v", kvalue.Type())
	}

	return roll(rvalue.(Number), lvalue.(Number), kvalue.(Number))
}

func literal(proxy *Nexter) (Value, error) {
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
		if val.Type() != NUMBERVALUE {
			return val, fmt.Errorf("Expected number, found %v", val.Type())
		}
		return Number{val: val.(Number).V(), repr: fmt.Sprintf("(%v)", val)}, nil
	case NUMBER:
		val, err := number(proxy)
		return val, err
	default:
		return Number{}, fmt.Errorf("Expected number or expression, found '%s'", proxy.Peek().t)
	}

}

func number(proxy *Nexter) (Number, error) {
	if proxy.Peek().t != NUMBER {
		return Number{}, fmt.Errorf("Expected number, got '%v'", proxy.Peek().value)
	}
	value, err := strconv.Atoi(proxy.Pop().value)
	if err != nil {
		return Number{}, err
	}
	var ret Number
	ret.Init(value)
	return ret, nil
}

func roll(number, size, keep Number) (Number, error) {
	if size.V() < 1 {
		return number, fmt.Errorf("Unexpected size '%v'", size)
	}

	results := make([]int, number.V())
	for i := 0; i < number.V(); i++ {
		tmp := rand.Intn(size.V()) + 1
		results[i] = tmp
	}
	sorted := sort.IntSlice(results)
	sort.Sort(sorted)

	repr := ""
	if number.V()-keep.V() > 0 {
		for _, v := range results[:number.V()-keep.V()] {
			if len(repr) > 0 {
				repr = fmt.Sprintf("%s+%v", repr, v)
			} else {
				repr = fmt.Sprintf("%v", v)
			}
		}
		repr = fmt.Sprintf("~~%s~~", repr)
	}
	sum := 0
	lowerBound := number.V() - keep.V()
	if lowerBound < 0 {
		lowerBound = 0
	}
	for _, v := range results[lowerBound:number.V()] {
		sum += v
		singleRepr := fmt.Sprintf("%v", v)
		if v == size.V() {
			singleRepr = fmt.Sprintf("**%s**", singleRepr)
		}

		if len(repr) > 0 {
			repr = fmt.Sprintf("%s+%s", repr, singleRepr)
		} else {
			repr = fmt.Sprintf("%s", singleRepr)
		}
	}

	return Number{val: sum, repr: fmt.Sprintf("(%s)", repr)}, nil
}
