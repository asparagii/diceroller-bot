package lang

import "fmt"

type Number struct {
	val  int
	repr string
}

func (n *Number) Init(val int) {
	n.val = val
	n.repr = fmt.Sprintf("%d", val)
}
func (n Number) String() string {
	return n.repr
}
func (a Number) Add(b Value) (Value, error) {
	if b.Type() != NUMBERVALUE {
		return a, fmt.Errorf("Cannot sum a number with a %v", b.Type())
	}
	repr := fmt.Sprintf("%v+%v", a, b)
	val := a.V() + b.(Number).V()
	return Number{val: val, repr: repr}, nil
}
func (n Number) Invert() (Value, error) {
	return Number{val: -n.V(), repr: fmt.Sprintf("-%v", n)}, nil
}
func (a Number) Mul(b Value) (Value, error) {
	if b.Type() != NUMBERVALUE {
		return a, fmt.Errorf("Cannot multiply a number with a %v", b.Type())
	}
	repr := fmt.Sprintf("%v*%v", a, b)
	val := a.V() * b.(Number).V()
	return Number{val: val, repr: repr}, nil
}
func (n Number) Type() ValueType {
	return NUMBERVALUE
}
func (n Number) V() int {
	return n.val
}
func (n *Number) Surround() {
	n.repr = fmt.Sprintf("(%s)", n.repr)
}
