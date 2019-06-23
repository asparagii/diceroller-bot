package lang

func Interpret(tree node) (Object, string, error) {
	mem := Memory{dollar: &Object{}}
	result, err := tree.evaluate(mem)
	return result, tree.String(), err
}
