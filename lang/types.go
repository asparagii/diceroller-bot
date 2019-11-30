package lang

import (
	"fmt"
	"strings"
)

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

func Boolean(value bool) Object {
	return Object{
		Value: value,
		Type:  BOOLEANVALUE,
	}
}

func (o Object) String() string {
	switch o.Type {
	case NUMBERVALUE:
		return fmt.Sprintf("%d", o.Value)
	case ARRAYVALUE:
		elements := o.Value.([]Object)
		var builder strings.Builder
		builder.WriteString("[ ")
		if len(elements) > 0 {
			builder.WriteString((elements)[0].String())
		}
		for _, v := range elements[1:] {
			builder.WriteString(", ")
			builder.WriteString(v.String())
		}
		builder.WriteString(" ]")
		return builder.String()
	case BOOLEANVALUE:
		if o.Value.(bool) {
			return "true"
		} else {
			return "false"
		}
	default:
		return "Error: Not implemented!"
	}
}
