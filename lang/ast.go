package lang

import (
	"errors"
	"fmt"
	"strings"
)

type node interface {
	evaluate(Memory) (Object, error)
	String() string
}

type Memory struct {
	dollar *Object
}

type ValueNode Object

func (n *ValueNode) evaluate(memory Memory) (Object, error) {
	return Object(*n), nil
}

func (n ValueNode) String() string {
	return Object(n).String()
}

type dollarNode struct {
	evaluatedValue Object
}

func (d *dollarNode) evaluate(memory Memory) (Object, error) {
	if memory.dollar == nil {
		return Object{}, errors.New("trying to read $, but it was never assigned")
	}
	d.evaluatedValue = *memory.dollar
	return *memory.dollar, nil
}

func (d dollarNode) String() string {
	return d.evaluatedValue.String()
}

type arrayNode []*node

func (arr *arrayNode) evaluate(m Memory) (Object, error) {
	results := make([]Object, len(*arr))
	for i, v := range *arr {
		tmp, err := (*v).evaluate(m)
		if err != nil {
			return tmp, err
		}
		results[i] = tmp
	}
	return Array(results), nil
}

func (arr arrayNode) String() string {
	var builder strings.Builder
	builder.WriteRune('[')
	if len(arr) > 0 {
		builder.WriteString((*arr[0]).String())
	}
	for _, v := range arr[1:] {
		builder.WriteRune(',')
		builder.WriteString((*v).String())
	}
	builder.WriteRune(']')
	return builder.String()
}

type PipeNode struct {
	left                node
	right               node
	leftRepresentation  string
	rightRepresentation []string
}

func (p *PipeNode) evaluate(m Memory) (Object, error) {
	leftResult, err := p.left.evaluate(m)
	if err != nil {
		return Object{}, err
	}
	p.leftRepresentation = p.left.String()
	var newMemory Memory
	switch leftResult.Type {
	case NUMBERVALUE:
		tmpdollar := Object(leftResult)
		newMemory.dollar = &tmpdollar
		result, err := p.right.evaluate(newMemory)
		p.rightRepresentation = []string{p.right.String()}
		return result, err
	case ARRAYVALUE:
		result := make([]Object, len(leftResult.Value.([]Object)))
		representations := make([]string, len(leftResult.Value.([]Object)))
		for i, v := range leftResult.Value.([]Object) {
			newMemory.dollar = &v
			singleResult, err := p.right.evaluate(newMemory)
			if err != nil {
				return Object{}, err
			}
			result[i] = singleResult
			representations[i] = p.right.String()
		}
		p.rightRepresentation = representations
		return Array(result), nil
	default:
		return Object{}, fmt.Errorf("Unsupported operand of type %v for pipe operation", leftResult.Type)
	}
}

func (p PipeNode) String() string {
	var builder strings.Builder
	builder.WriteString(p.leftRepresentation)
	builder.WriteString(" => \n")
	builder.WriteString("[ ")
	if len(p.rightRepresentation) > 0 {
		builder.WriteString(p.rightRepresentation[0])
	}
	for _, v := range p.rightRepresentation[1:] {
		builder.WriteString(", ")
		builder.WriteString(v)
	}
	builder.WriteString(" ]")
	return builder.String()
}

type SumNode struct {
	left   node
	right  node
	result Object
}

func (n *SumNode) evaluate(m Memory) (Object, error) {
	lval, err := n.left.evaluate(m)
	if err != nil {
		return Object{}, err
	}
	rval, err := n.right.evaluate(m)
	if err != nil {
		return Object{}, err
	}
	return Add(lval, rval)
}

func (n SumNode) String() string {
	lrepr := n.left.String()
	rrepr := n.right.String()
	return fmt.Sprintf("%s+%s", lrepr, rrepr)
}

type SubtractNode struct {
	left   node
	right  node
	result Object
}

func (n *SubtractNode) evaluate(m Memory) (Object, error) {
	lval, err := n.left.evaluate(m)
	if err != nil {
		return Object{}, err
	}
	rval, err := n.right.evaluate(m)
	if err != nil {
		return Object{}, err
	}
	return Subtract(lval, rval)
}

func (n SubtractNode) String() string {
	lrepr := n.left.String()
	rrepr := n.right.String()
	return fmt.Sprintf("%s-%s", lrepr, rrepr)
}

type MultiplyNode struct {
	left   node
	right  node
	result Object
}

func (n MultiplyNode) evaluate(m Memory) (Object, error) {
	lval, err := n.left.evaluate(m)
	if err != nil {
		return Object{}, err
	}
	rval, err := n.right.evaluate(m)
	if err != nil {
		return Object{}, err
	}
	return Multiply(lval, rval)
}

func (n MultiplyNode) String() string {
	lrepr := n.left.String()
	rrepr := n.right.String()
	return fmt.Sprintf("(%s)*(%s)", lrepr, rrepr)
}

type DiceKeepNode struct {
	number         node
	size           node
	keep           node
	representation string
}

func (n *DiceKeepNode) evaluate(m Memory) (Object, error) {
	leftval, err := n.number.evaluate(m)
	if err != nil {
		return Object{}, err
	}
	midval, err := n.size.evaluate(m)
	if err != nil {
		return Object{}, err
	}
	rightval, err := n.keep.evaluate(m)
	if err != nil {
		return Object{}, err
	}

	result, repr, err := RollKeep(leftval, midval, rightval)
	n.representation = repr
	return result, err
}

func (n DiceKeepNode) String() string {
	return n.representation
}

type DiceNode struct {
	number         node
	size           node
	representation string
}

func (n *DiceNode) evaluate(m Memory) (Object, error) {
	number, err := n.number.evaluate(m)
	if err != nil {
		return Object{}, err
	}
	size, err := n.size.evaluate(m)
	if err != nil {
		return Object{}, err
	}

	result, repr, err := Roll(number, size)
	n.representation = repr
	return result, err
}

func (n DiceNode) String() string {
	return n.representation
}

type ColonNode struct {
	start  node
	end    node
	result Object
}

func (n *ColonNode) evaluate(m Memory) (Object, error) {
	start, err := n.start.evaluate(m)
	if err != nil {
		return start, err
	}
	if start.Type != NUMBERVALUE {
		return start, fmt.Errorf("Unexpected object of type %v as start of range", start.Type)
	}

	end, err := n.end.evaluate(m)
	if err != nil {
		return end, err
	}
	if end.Type != NUMBERVALUE {
		return end, fmt.Errorf("Unexpected object of type %v as end of range", end.Type)
	}

	var (
		s = start.Value.(int)
		e = end.Value.(int)
	)

	if e-s < 0 {
		return start, fmt.Errorf("Invalid range: 'start' cannot be less than 'end'")
	}

	result := make([]Object, e-s)
	for i := 0; i < e-s; i++ {
		result[i] = Number(s + i)
	}
	n.result = Array(result)
	return n.result, nil
}

func (n ColonNode) String() string {
	return n.result.String()
}

type ComparisonNode struct {
	left                    node
	right                   node
	operation               func(a, b Object) (Object, error)
	operationRepresentation string
	representation          string
}

func (n *ComparisonNode) evaluate(memory Memory) (Object, error) {
	left, err := n.left.evaluate(memory)
	if err != nil {
		return left, err
	}

	right, err := n.right.evaluate(memory)
	if err != nil {
		return right, err
	}

	result, err := n.operation(left, right)
	if err != nil {
		return result, err
	}

	n.representation = fmt.Sprintf("%v %s %v", n.left, n.operationRepresentation, n.right)
	return result, nil
}

func (n ComparisonNode) String() string {
	return n.representation
}

func LessThanNode(left, right node) ComparisonNode {
	return ComparisonNode{
		left:                    left,
		right:                   right,
		operation:               Less,
		operationRepresentation: "<",
	}
}
func GreaterThanNode(left, right node) ComparisonNode {
	return ComparisonNode{
		left:                    left,
		right:                   right,
		operation:               More,
		operationRepresentation: ">",
	}
}
func EqualToNode(left, right node) ComparisonNode {
	return ComparisonNode{
		left:                    left,
		right:                   right,
		operation:               Equal,
		operationRepresentation: "=",
	}
}
