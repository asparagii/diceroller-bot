package lang

import (
	"fmt"
)

func Number(value int) Object {
	p := make(map[Property]interface{})
	p[VALUE] = value
	p[TYPE] = NUMBERVALUE

	m := make(map[Method]func(a, b Object) (Object, error))
	m[ADD] = AddNumber
	m[SUBTRACT] = SubtractNumber
	m[MULTIPLY] = MultiplyNumber

	return Object{properties: p, methods: m}
}

func AddNumber(self, other Object) (Object, error) {
	if other.properties[TYPE] != NUMBERVALUE {
		return self, fmt.Errorf("Invalid type: can't sum Number with %v", other.properties[TYPE])
	}
	return Number(self.properties[VALUE].(int) + other.properties[VALUE].(int)), nil
}

func SubtractNumber(self, other Object) (Object, error) {
	if other.properties[TYPE] != NUMBERVALUE {
		return self, fmt.Errorf("Invalid type: can't subtract Number with %v", other.properties[TYPE])
	}
	return Number(self.properties[VALUE].(int) - other.properties[VALUE].(int)), nil
}

func InvertNumber(self Object) (Object, error) {
	return Number(self.properties[VALUE].(int)), nil
}

func MultiplyNumber(self, other Object) (Object, error) {
	if other.properties[TYPE] != NUMBERVALUE {
		return self, fmt.Errorf("Invalid type: can't multiply Number with %v", other.properties[TYPE])
	}
	return Number(self.properties[VALUE].(int) * other.properties[VALUE].(int)), nil
}
