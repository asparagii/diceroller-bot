package lang

type Node struct {
	Value NodeType
	Left  *Node
	Right *Node
}

type NodeType int32

const (
	ValueNode NodeType = iota
	DollarNode
	ArrayNode
	PipeNode
	SumNode
	SubtractNode
	MultiplyNode
	DiceNode
	KeepNode
	ColonNode
	LessThanNode
	MoreThanNode
	EqualNode
)
