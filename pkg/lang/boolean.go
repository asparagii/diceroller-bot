package lang

/*
import "fmt"

type Boolean struct {
	val  bool
	repr string
}

func (b Boolean) String() string {
	return b.repr
}
func (a Boolean) Add(b Value) (Value, error) {
	if b.Type() != BOOLEANVALUE {
		return &a, fmt.Errorf("Cannot perform `add` operation between Boolean and non-Boolean values")
	}
	return Boolean(a.V() || b.(Boolean).V()), nil
}
func (a Boolean) Mul(b Value) (Value, error) {
	if b.Type() != BOOLEANVALUE {
		return &a, fmt.Errorf("Cannot perform `add` operation between Boolean and non-Boolean values")
	}
	return Boolean(a.V() && b.(Boolean).V()), nil
}
func (n Boolean) Invert() (Value, error) {
	return Boolean{!n.V(), fmt.Sprintf("%v", !n.V())}, nil
}
func (n Boolean) Type() ValueType {
	return BOOLEANVALUE
}
func (n *Boolean) Surround() {
	n.repr = fmt.Sprintf("(%s)", n.String())
}
func (n Boolean) V() interface{} {
	return n.val
}
*/
