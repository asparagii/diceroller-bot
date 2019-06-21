package lang

func Interpret(tree node) (Object, error) {
	mem := Memory{dollar: &Object{}}
	return tree.evaluate(mem)
}
