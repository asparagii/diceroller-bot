package lang

func Number(value int) Object {
	return Object{Value: value, Type: NUMBERVALUE}
}
