package lang

type Object struct {
	Value interface{}
	Type  ValueType
}

type ValueType int

const (
	NUMBERVALUE ValueType = iota
	BOOLEANVALUE
	ARRAYVALUE
)

func Number(value int) Object {
	return Object{Value: value, Type: NUMBERVALUE}
}

func Array(elements []Object) Object {
	return Object{Value: elements, Type: ARRAYVALUE}
}

/*
type ArrayNumber []Value

func (n ArrayNumber) String() string {
	representations := make([]string, len(n))
	for i, v := range n {
		representations[i] = v.String()
	}
	return fmt.Sprintf("{ %s }", strings.Join(representations, ", "))
}

func (n ArrayNumber) Map(f Expression) (ArrayNumber, error) {
	result := make(ArrayNumber, len(n))
	for i, v := range n {
		result[i] = f.Eval(v)
	}
	return result, nil
}

func (n ArrayNumber) Filter(f Expression) (ArrayNumber, error) {
	var result ArrayNumber
	for _, v := range n {
		evaluation := f.Eval(v)
		if evaluation.Type() != BOOLEANVALUE {
			return n, fmt.Errorf("Filter expression did not return a boolean value")
		}

		if bool(evaluation.(Boolean)) {
			result = append(result, v)
		}
	}
	return result, nil
}
*/
