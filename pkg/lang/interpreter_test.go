package lang

import (
	"testing"
)

func TestInterpretSum(t *testing.T) {
	val := valueNode(Number(1))
	tree := SumNode(&val, &val)
	result, err := Interpret(tree)
	assert(t, err == nil, "Expected no error, got %v", err)
	assert(t, result.Value == Number(2).Value, "Expected number 2, got %v", result)
}

func TestInterpretPipe(t *testing.T) {
	val := valueNode(Number(1))
	tree := PipeNode(&val, &val)
	result, err := Interpret(tree)
	assert(t, err == nil, "Expected no error, got %v", err)
	assert(t, result.Value == Number(1).Value, "Expected number 1, got %v", result)
}
