package lang

type Nexter struct {
	token    chan Token
	internal Token
}

func (n *Nexter) Init() {
	n.internal = <-n.token
}

func (n *Nexter) Pop() Token {
	tmp := n.internal
	if n.internal.t != EOF {
		n.internal = <-n.token
	}
	return tmp
}

func (n *Nexter) Peek() Token {
	return n.internal
}
