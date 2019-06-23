package lang

import (
	"fmt"
	"strconv"
)

func Start(input string) (Object, string, error) {
	ast, err := Parse(input)
	if err != nil {
		return Object{}, "", err
	}
	return Interpret(ast)
}

func Parse(expression string) (node, error) {
	token := Lex(expression)

	proxy := Nexter{token: token}

	// Read first value
	proxy.Init()
	tree, err := program(&proxy)
	if proxy.Peek().t != EOF {
		return tree, fmt.Errorf("Unexpected end of message")
	}
	return tree, err
}

func program(proxy *Nexter) (node, error) {
	ret, err := expr(proxy)
	if err != nil {
		return ret, err
	}

	for true {
		switch proxy.Peek().t {
		case PIPE:
			proxy.Pop()
			right, err := expr(proxy)
			if err != nil {
				return right, err
			}
			ret = &PipeNode{
				left:  ret,
				right: right,
			}
		default:
			return ret, nil
		}
	}
	return ret, fmt.Errorf("This shouldn't happen")
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
			lvalue = &SumNode{
				left:  lvalue,
				right: rvalue,
			}
		case SUB:
			proxy.Pop()
			rvalue, err := term(proxy)
			if err != nil {
				return lvalue, err
			}
			lvalue = &SubtractNode{
				left:  lvalue,
				right: rvalue,
			}
		default:
			return lvalue, nil
		}
	}
	return lvalue, fmt.Errorf("Unexpected end of the message")
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
			lvalue = &MultiplyNode{
				left:  lvalue,
				right: rvalue,
			}
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
		return &DiceNode{
			number: rvalue,
			size:   lvalue,
		}, nil
	}
	proxy.Pop()
	kvalue, err := literal(proxy)
	if err != nil {
		return kvalue, err
	}

	return &DiceKeepNode{
		number: rvalue,
		size:   lvalue,
		keep:   kvalue,
	}, nil
}

func literal(proxy *Nexter) (node, error) {
	switch proxy.Peek().t {
	case LPAREN:
		proxy.Pop()
		val, err := expr(proxy)
		if err != nil {
			return val, err
		}
		if proxy.Peek().t != RPAREN {
			return val, fmt.Errorf("Unmaching parentesis")
		}
		proxy.Pop()
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
		return &dollarNode{}, nil
	default:
		return &ValueNode{}, fmt.Errorf("Expected number or expression, found '%s'", proxy.Peek().t)
	}
}

func list(proxy *Nexter) (node, error) {
	nodes := make([]*node, 0)
	if proxy.Peek().t == RARRAYPAREN {
		ret := arrayNode(nodes)
		return &ret, nil
	} else {
		val, err := expr(proxy)
		if err != nil {
			return &ValueNode{}, err
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
			nodes = append(nodes, &val)
		}
		ret := arrayNode(nodes)
		return &ret, nil
	}
}

func number(proxy *Nexter) (node, error) {
	if proxy.Peek().t != NUMBER {
		return &ValueNode{}, fmt.Errorf("Expected number, got '%v'", proxy.Peek().value)
	}
	value, err := strconv.Atoi(proxy.Pop().value)
	if err != nil {
		return &ValueNode{}, err
	}
	ret := ValueNode(Number(value))
	return &ret, nil
}
