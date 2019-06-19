package lang

import (
	"fmt"
	"strconv"
)

func Parse(expression string) (Object, error) {
	token := Lex(expression)

	proxy := Nexter{token: token}

	// Read first value
	proxy.Init()
	value, err := expr(&proxy)
	if err != nil {
		return Object{}, err
	}
	return Interpret(value)
	//	if proxy.Peek().t != EOF {
	//		return value, fmt.Errorf("Unexpected end of message")
	//	} else {
	//		return value, err
	//	}
}

func Interpret(tree node) (Object, error) {
	mem := Memory{dollar: &Object{}}
	return tree.evaluate(mem)
}

func expr(proxy *Nexter) (node, error) {
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
			lvalue = SumNode(lvalue, rvalue)
		case SUB:
			proxy.Pop()
			rvalue, err := term(proxy)
			if err != nil {
				return lvalue, err
			}
			lvalue = SubtractNode(lvalue, rvalue)
		default:
			return lvalue, nil
		}
	}
	return lvalue, fmt.Errorf("Unexpected end of the parser")
}

func term(proxy *Nexter) (node, error) {
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
			lvalue = MultiplyNode(lvalue, rvalue)
		default:
			return lvalue, nil
		}
	}
	return lvalue, fmt.Errorf("Unexpected end of the parser")
}

func atom(proxy *Nexter) (node, error) {
	rvalue, err := literal(proxy)
	if err != nil {
		return rvalue, err
	}

	if proxy.Peek().t != DICE {
		return rvalue, nil
	}

	proxy.Pop()
	lvalue, err := literal(proxy)
	if err != nil {
		return lvalue, err
	}
	if proxy.Peek().t != KEEP {
		return DiceNode(rvalue, lvalue), nil
	}
	proxy.Pop()
	kvalue, err := literal(proxy)
	if err != nil {
		return kvalue, err
	}

	return DiceKeepNode(rvalue, lvalue, kvalue), nil
}

func literal(proxy *Nexter) (node, error) {
	switch proxy.Peek().t {
	case LPAREN:
		proxy.Pop()
		val, err := expr(proxy)
		if err != nil {
			return val, err
		}
		return val, nil
	case LARRAYPAREN:
		proxy.Pop()
		val, err := list(proxy)
		if err != nil {
			return val, err
		}
		if proxy.Peek().t != RARRAYPAREN {
			return val, fmt.Errorf("Error: expected end of list, found '%v'", val)
		}
		proxy.Pop()
		return val, nil
	case NUMBER:
		val, err := number(proxy)
		return val, err
	case DOLLAR:
		proxy.Pop()
		return dollarNode{}, nil //TODO change
	default:
		return valueNode{}, fmt.Errorf("Expected number or expression, found '%s'", proxy.Peek().t)
	}
}

func list(proxy *Nexter) (node, error) {
	nodes := make([]*node, 0)
	if proxy.Peek().t == RARRAYPAREN {
		return arrayNode(nodes), nil
	} else {
		val, err := expr(proxy)
		if err != nil {
			return valueNode{}, err
		}
		nodes = append(nodes, &val)

		for proxy.Peek().t != RARRAYPAREN {
			if proxy.Peek().t != LISTSEPARATOR {
				return val, fmt.Errorf("Expected ',' or ']', found '%v'", proxy.Peek().t)
			}
			proxy.Pop()
			val, err := expr(proxy)
			if err != nil {
				return val, err
			}
			nodes = append(nodes, &val) //TODO change to tree building
		}
		return arrayNode(nodes), nil
	}
}

func number(proxy *Nexter) (node, error) {
	if proxy.Peek().t != NUMBER {
		return valueNode{}, fmt.Errorf("Expected number, got '%v'", proxy.Peek().value)
	}
	value, err := strconv.Atoi(proxy.Pop().value)
	if err != nil {
		return valueNode{}, err
	}
	return valueNode(Number(value)), nil
}
