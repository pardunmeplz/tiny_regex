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
