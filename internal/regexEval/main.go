package regex

func HasUnion(regexA string, regexB string) (bool, []string) {
	astA, errA := parse(regexA)
	astB, errB := parse(regexB)
	if len(errA) > 0 || len(errB) > 0 {
		return false, append(errA, errB...)
	}

	nfaA := thompsonConstruction(astA)
	nfaB := thompsonConstruction(astB)

	return evaluateNfaProduct(nfaA, nfaB), nil
}

func Matches(regex string, input string) (bool, []string) {
	ast, err := parse(regex)
	if len(err) > 0 {
		return false, err
	}

	nfa := thompsonConstruction(ast)
	return eval(input, nfa), nil
}
