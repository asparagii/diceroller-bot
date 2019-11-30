package lang

import (
	"testing"
)

func TestParseParentesis(t *testing.T) {
	ast, err := Parse("1*(2+2)")
	assert(t, err == nil, "Expected no error, got %v", err)
	root, ok := ast.(*MultiplyNode) // Type assertion
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
