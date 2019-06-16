package lang

import (
	"fmt"
)

type BinaryOperation func(a, b Object) (Object, error)

type UnaryOperation func(a Object) (Object, error)

type TernaryOperation func(a, b, c Object) (Object, error)

func Add(a, b Object) (Object, error) {
	if a.Type == b.Type {
		switch a.Type {
		case NUMBERVALUE:
			return Number(a.Value.(int) + b.Value.(int)), nil
		default:
			return Object{}, fmt.Errorf("Error: unknown type")
		}
	}

	return Object{}, fmt.Errorf("Error: unknown type")
}

func Subtract(a, b Object) (Object, error) {
	if a.Type == b.Type {
		switch a.Type {
		case NUMBERVALUE:
			return Number(a.Value.(int) - b.Value.(int)), nil
		default:
			return Object{}, fmt.Errorf("Error: unknown type")
		}
	}
	return Object{}, fmt.Errorf("Error: unknown type")
}

func Multiply(a, b Object) (Object, error) {
	switch a.Type {
	case NUMBERVALUE:
		switch b.Type {
		case NUMBERVALUE:
			return Number(a.Value.(int) * b.Value.(int)), nil
		default:
			return Object{}, fmt.Errorf("Error: unknown type for right operand")
		}
	default:
		return Object{}, fmt.Errorf("Error: unknown type")
	}
	return Object{}, fmt.Errorf("Error: unknown type")
}
