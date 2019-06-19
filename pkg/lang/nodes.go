package lang

func SumNode(lnode, rnode node) *biNode {
	ret := new(biNode)
	ret.operation = Add
	ret.left = &lnode
	ret.right = &rnode
	return ret
}

func SubtractNode(lnode, rnode node) *biNode {
	return &biNode{
		operation: Subtract,
		left:      &lnode,
		right:     &rnode,
	}
}

func MultiplyNode(lnode, rnode node) *biNode {
	return &biNode{
		operation: Multiply,
		left:      &lnode,
		right:     &rnode,
	}
}

func DiceKeepNode(number, size, keep node) *triNode {
	return &triNode{
		operation: RollKeep,
		left:      &number,
		middle:    &size,
		right:     &keep,
	}
}

func DiceNode(number, size node) *biNode {
	return &biNode{
		operation: Roll,
		left:      &number,
		right:     &size,
	}
}
