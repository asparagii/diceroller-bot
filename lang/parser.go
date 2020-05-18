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

func Parse(expression string) (*Node, error) {
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

///
// Grammar definition:
//
//	PROGRAM := COMP (\| COMP)*
//	COMP := EXPR ([\<\>\=] EXPR)?
//	EXPR := TERM ([\+\-] TERM)*
//	TERM := ATOM (\* ATOM)*
//	ATOM := LITERAL ( d LITERAL ( k LITERAL )? )?
//	LITERAL := [( \(COMP\) )( \[LIST\] )(NUMBER)(DOLLAR)]
//	LIST := COMP ( \, COMP )*
//	NUMBER := [1234567890]+
//	DOLLAR := \$
//
///

func program(proxy *Nexter) (*Node, error) {
	ret, err := comp(proxy)
	if err != nil {
		return ret, err
	}

	for {
		switch proxy.Peek().t {
		case PIPE:
			proxy.Pop()
			right, err := comp(proxy)
			if err != nil {
				return right, err
			}
			ret = &Node{
				Value: PipeNode,
				Left:  ret,
				Right: right,
			}
		default:
			return ret, nil
		}
	}
}

func comp(proxy *Nexter) (*Node, error) {
	lvalue, err := expr(proxy)
	if err != nil {
		return lvalue, err
	}
	switch proxy.Peek().t {
	case LESS:
		proxy.Pop()
		rvalue, err := expr(proxy)
		if err != nil {
			return rvalue, err
		}
		return &Node{
			Value: LessThanNode,
			Left:  lvalue,
			Right: rvalue,
		}, nil
	case MORE:
		proxy.Pop()
		rvalue, err := expr(proxy)
		if err != nil {
			return rvalue, err
		}

		return &Node{
			Value: MoreThanNode,
			Left:  lvalue,
			Right: rvalue,
		}, nil
	case EQUAL:
		proxy.Pop()
		rvalue, err := expr(proxy)
		if err != nil {
			return rvalue, err
		}
		return &Node{
			Value: EqualNode,
			Left:  lvalue,
			Right: rvalue,
		}, nil
	default:
		return lvalue, nil
	}
}

func expr(proxy *Nexter) (*Node, error) {
	lvalue, err := term(proxy)
	if err != nil {
		return lvalue, err
	}

	for {
		switch proxy.Peek().t {
		case SUM:
			proxy.Pop()
			rvalue, err := term(proxy)
			if err != nil {
				return lvalue, err
			}
			lvalue = &Node{
				Value: SumNode,
				Right: rvalue,
				Left:  lvalue,
			}
		case SUB:
			proxy.Pop()
			rvalue, err := term(proxy)
			if err != nil {
				return lvalue, err
			}
			lvalue = &Node{
				Value: SubtractNode,
				Left:  lvalue,
				Right: rvalue,
			}
		default:
			return lvalue, nil
		}
	}
}

func term(proxy *Nexter) (*Node, error) {
	lvalue, err := atom(proxy)

	if err != nil {
		return lvalue, err
	}

	for {
		switch proxy.Peek().t {
		case MUL:
			proxy.Pop()
			rvalue, err := atom(proxy)
			if err != nil {
				return lvalue, err
			}
			lvalue = &Node{
				Value: MultiplyNode,
				Left:  lvalue,
				Right: rvalue,
			}
		default:
			return lvalue, nil
		}
	}
}

func atom(proxy *Nexter) (*Node, error) {
	lvalue, err := literal(proxy)
	if err != nil {
		return lvalue, err
	}

	switch proxy.Peek().t {
	case DICE:
		proxy.Pop()
		rvalue, err := literal(proxy)
		if err != nil {
			return rvalue, err
		}
		return &Node{
			Value: DiceNode,
			Left:  lvalue,
			Right: rvalue,
		}, nil
	case COLON:
		proxy.Pop()
		rvalue, err := literal(proxy)
		if err != nil {
			return rvalue, err
		}
		return &Node{
			Value: ColonNode,
			Left:  lvalue,
			Right: rvalue,
		}, nil
	default:
		return lvalue, nil
	}
}

func literal(proxy *Nexter) (*Node, error) {
	switch proxy.Peek().t {
	case LPAREN:
		proxy.Pop()
		val, err := comp(proxy)
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
		return &Node{Value: DollarNode}, nil
	default:
		return &Node{}, fmt.Errorf("Expected number or expression, found '%s'", proxy.Peek().t)
	}
}

func list(proxy *Nexter) (*Node, error) {
	// what shall we do with the array node?
	// The goal here is to have a binary tree as parser result
	// Maybe we need to introduce an array terminator of some sort,
	// representing a nil element or something
	// Yup that seems the way to go

	nodes := make([]*Node, 0)
	if proxy.Peek().t == RARRAYPAREN {
		ret := arrayNode(nodes)
		return &ret, nil
	} else {
		val, err := comp(proxy)
		if err != nil {
			return &ValueNode{}, err
		}
		nodes = append(nodes, &val)

		for proxy.Peek().t != RARRAYPAREN {
			if proxy.Peek().t != LISTSEPARATOR {
				return val, fmt.Errorf("Expected ',' or ']', found '%v'", proxy.Peek().t)
			}
			proxy.Pop()
			val, err := comp(proxy)
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
