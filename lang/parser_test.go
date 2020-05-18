package lang

import (
	"testing"
)

func TestParseParentesis(t *testing.T) {
	ast, err := Parse("1*(2+2)")
	assert(t, err == nil, "Expected no error, got %v", err)
	root, ok := ast.Value 
	assert(t, ok, "Expected root of type MultiplyNode, found %T", ast)
	_, ok = root.left.(*ValueNode)
	assert(t, ok, "Expected root of type ValueNode, found %T", root.left)
	right, ok := root.right.(*SumNode)
	assert(t, ok, "Expected root of type SumNode, found %T", root.right)
	_, ok = right.left.(*ValueNode)
	assert(t, ok, "Expected root of type ValueNode, found %T", right.left)
	_, ok = right.right.(*ValueNode)
	assert(t, ok, "Expected root of type ValueNode, found %T", right.right)
}

func TestParseDice(t *testing.T) {
	ast, err := Parse("(3*2)d(2*3d2)")
	assert(t, err == nil, "Expected no error, got %v", err)
	root, ok := ast.(*DiceNode)
	assert(t, ok, "Expected root of type DiceNode, found %T", ast)
	_, ok = root.number.(*MultiplyNode)
	assert(t, ok, "Expected left of type MultiplyNode, found %T", root.number)
	right, ok := root.size.(*MultiplyNode)
	assert(t, ok, "Expected right of type MultiplyNode, found %T", root.size)
	_, ok = right.left.(*ValueNode)
	assert(t, ok, "Expected right->left of type ValueNode, found %T", right.left)
	_, ok = right.right.(*DiceNode)
	assert(t, ok, "Expected right->right of type DiceNode, found %T", right.right)
}
