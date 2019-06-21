package lang

import (
	"fmt"
)

type Representation interface {
	String() string
	Negate() Representation
	Append(string) Representation
}

type DiceRepresentation struct {
	internal string
	negated  bool
}

func NewDiceRepresentation(s string) string {
	return DiceRepresentation{internal: s, negated: false}
}

func (r DiceRepresentation) String() string {
	if r.negated {
		return fmt.Sprintf("~~%s~~", r.internal)
	} else {
		return r.internal
	}
}

func (r DiceRepresentation) Negate() DiceRepresentation {
	return DiceRepresentation{
		internal: r.internal,
		negated:  false,
	}
}

func (r DiceRepresentation) Append(s string) DiceRepresentation {
	return DiceRepresentation{
		internal: fmt.Sprintf("%s%s", r.internal, s),
		negated:  r.negated,
	}
}
