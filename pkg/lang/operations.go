package lang

import (
	"fmt"
	"math/rand"
	"sort"
)

type UnaryOperation func(a Object) (Object, error)

type BinaryOperation func(a, b Object) (Object, error)

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

func Roll(number, size Object) (Object, string, error) {
	return RollKeep(number, size, number)
}

func RollKeep(number, size, keep Object) (Object, string, error) {
	if number.Type != NUMBERVALUE || size.Type != NUMBERVALUE || keep.Type != NUMBERVALUE {
		return Object{}, "", fmt.Errorf("Wrong type for roll")
	}

	if size.Value.(int) < 1 {
		return Object{}, "", fmt.Errorf("Unexpected size '%v'", size)
	}

	results := make([]int, number.Value.(int))
	for i := 0; i < number.Value.(int); i++ {
		tmp := rand.Intn(size.Value.(int)) + 1
		results[i] = tmp
	}
	sorted := sort.IntSlice(results)
	sort.Sort(sorted)

	repr := ""
	if number.Value.(int)-keep.Value.(int) > 0 {
		for _, v := range results[:number.Value.(int)-keep.Value.(int)] {
			if len(repr) > 0 {
				repr = fmt.Sprintf("%s+%v", repr, v)
			} else {
				repr = fmt.Sprintf("%v", v)
			}
		}
		repr = fmt.Sprintf("~~%s~~", repr)
	}
	sum := 0
	lowerBound := number.Value.(int) - keep.Value.(int)
	if lowerBound < 0 {
		lowerBound = 0
	}
	for _, v := range results[lowerBound:number.Value.(int)] {
		sum += v
		singleRepr := fmt.Sprintf("%v", v)
		if v == size.Value.(int) {
			singleRepr = fmt.Sprintf("**%s**", singleRepr)
		}

		if len(repr) > 0 {
			repr = fmt.Sprintf("%s+%s", repr, singleRepr)
		} else {
			repr = fmt.Sprintf("%s", singleRepr)
		}
	}

	return Number(sum), fmt.Sprintf("(%s)", repr), nil
}

func Less(a, b Object) (Object, error) {
	if a.Type != NUMBERVALUE {
		return Object{}, fmt.Errorf("Unexpected type for operand '<': %v", a.Type)
	}

	if b.Type != NUMBERVALUE {
		return Object{}, fmt.Errorf("Unexpected type for operand '<': %v", b.Type)
	}

	left := a.Value.(int)
	right := b.Value.(int)

	return Boolean(left < right), nil
}

func More(a, b Object) (Object, error) {
	if a.Type != NUMBERVALUE {
		return Object{}, fmt.Errorf("Unexpected type for operand '>': %v", a.Type)
	}

	if b.Type != NUMBERVALUE {
		return Object{}, fmt.Errorf("Unexpected type for operand '>': %v", b.Type)
	}

	left := a.Value.(int)
	right := b.Value.(int)

	return Boolean(left > right), nil
}

func Equal(a, b Object) (Object, error) {
	if a.Type != NUMBERVALUE {
		return Object{}, fmt.Errorf("Unexpected type for operand '=': %v", a.Type)
	}

	if b.Type != NUMBERVALUE {
		return Object{}, fmt.Errorf("Unexpected type for operand '=': %v", b.Type)
	}

	left := a.Value.(int)
	right := b.Value.(int)

	return Boolean(left == right), nil
}
