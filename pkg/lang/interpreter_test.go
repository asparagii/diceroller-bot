package lang

import (
	"strings"
	"testing"
)

func TestInterpretSum(t *testing.T) {
	val := ValueNode(Number(1))
	tree := &SumNode{
		left:  &val,
		right: &val,
	}

	result, representation, err := Interpret(tree)
	assert(t, err == nil, "Expected no error, got %v", err)
	assert(t, result.Value == Number(2).Value, "Expected number 2, got %v", result)
	assert(t, strings.Compare(representation, "1+1") == 0, "Expected '1+1', got %s", representation)
}

func TestInterpretSubtract(t *testing.T) {
	val := ValueNode(Number(1))
	tree := &SubtractNode{
		left:  &val,
		right: &val,
	}

	result, representation, err := Interpret(tree)
	assert(t, err == nil, "Expected no error, got %v", err)
	assert(t, result.Value == Number(0).Value, "Expected number 2, got %v", result)
	assert(t, strings.Compare(representation, "1-1") == 0, "Expected '1-1', got %s", representation)
}

func TestInterpretMultiply(t *testing.T) {
	val := ValueNode(Number(1))
	tree := &MultiplyNode{
		left:  &val,
		right: &val,
	}

	result, representation, err := Interpret(tree)
	assert(t, err == nil, "Expected no error, got %v", err)
	assert(t, result.Value == Number(1).Value, "Expected number 2, got %v", result)
	assert(t, strings.Compare(representation, "(1)*(1)") == 0, "Expected '(1)*(1)', got %s", representation)
}

func TestInterpretPipe(t *testing.T) {
	val := ValueNode(Array([]Object{Number(1), Number(2)}))
	right := dollarNode{}
	tree := &PipeNode{
		left:  &val,
		right: &right,
	}
	result, representation, err := Interpret(tree)
	assert(t, err == nil, "Expected no error, got %v", err)
	assert(t, result.Type == ARRAYVALUE, "Expected array result, got %v", result.Type)
	assertArrayEqual(t, result.Value.([]Object), []Object{Number(1), Number(2)})
	assert(t, strings.Compare(representation, "[ 1, 2 ] => \n[ 1, 2 ]") == 0, "Representation didn't match")
}

func assertArrayEqual(t *testing.T, a, b []Object) {
	assert(t, len(a) == len(b), "Sizes don't match")
	for i := range a {
		assert(t, a[i].Type == b[i].Type, "Types don't match at position %d", i)
		assert(t, a[i].Value == b[i].Value, "Values don't match at position %d", i)
	}
}
