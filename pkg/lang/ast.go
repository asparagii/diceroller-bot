package lang

import (
	"errors"
	"fmt"
)

type node interface {
	evaluate(Memory) (Object, error)
}

type Memory struct {
	dollar *Object
}

type valueNode Object

func (n valueNode) evaluate(memory Memory) (Object, error) {
	return Object(n), nil
}

type dollarNode struct{}

func (d dollarNode) evaluate(memory Memory) (Object, error) {
	if memory.dollar == nil {
		return Object{}, errors.New("trying to read $, but it was never assigned")
	}
	return *memory.dollar, nil
}

type uniNode struct {
	operation UnaryOperation
	left      *node
}

type biNode struct {
	operation BinaryOperation
	left      *node
	right     *node
}

type triNode struct {
	operation TernaryOperation
	left      *node
	middle    *node
	right     *node
}

type arrayNode []*node

func (arr arrayNode) evaluate(m Memory) (Object, error) {
	results := make([]Object, len(arr))
	for i, v := range arr {
		tmp, err := (*v).evaluate(m)
		if err != nil {
			return tmp, err
		}
		results[i] = tmp
	}
	return Array(results), nil
}

func (n uniNode) evaluate(memory Memory) (Object, error) {
	val, err := (*n.left).evaluate(memory)
	if err != nil {
		return val, err
	}
	return n.operation(val)
}

func (n biNode) evaluate(memory Memory) (Object, error) {
	leftVal, err := (*n.left).evaluate(memory)
	if err != nil {
		return leftVal, err
	}
	rightVal, err := (*n.right).evaluate(memory)
	if err != nil {
		return rightVal, err
	}
	return n.operation(leftVal, rightVal)
}

func (n triNode) evaluate(memory Memory) (Object, error) {
	leftVal, err := (*n.left).evaluate(memory)
	if err != nil {
		return leftVal, err
	}
	middleVal, err := (*n.middle).evaluate(memory)
	if err != nil {
		return middleVal, err
	}
	rightVal, err := (*n.right).evaluate(memory)
	if err != nil {
		return rightVal, err
	}
	return n.operation(leftVal, middleVal, rightVal)
}

type pipeNode struct {
	left  *node
	right *node
}

func PipeNode(lnode, rnode node) *pipeNode {
	return &pipeNode{
		left:  &lnode,
		right: &rnode,
	}
}

func (p pipeNode) evaluate(m Memory) (Object, error) {
	leftResult, err := (*p.left).evaluate(m)
	if err != nil {
		return Object{}, err
	}
	var newMemory Memory
	switch leftResult.Type {
	case NUMBERVALUE:
		tmpdollar := Object(leftResult)
		newMemory.dollar = &tmpdollar
		result, err := (*p.right).evaluate(newMemory)
		return result, err
	case ARRAYVALUE:
		result := make([]Object, len(leftResult.Value.([]Object)))
		for i, v := range leftResult.Value.([]Object) {
			newMemory.dollar = &v
			singleResult, err := (*p.right).evaluate(newMemory)
			if err != nil {
				return Object{}, err
			}
			result[i] = singleResult
		}
		return Array(result), nil
	default:
		return Object{}, fmt.Errorf("Unsupported operand of type %v for pipe operation", leftResult.Type)
	}
}
