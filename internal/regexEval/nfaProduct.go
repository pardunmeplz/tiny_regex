package regex

type Visited map[[2]*State]struct{}

func evaluateNfaProduct(a NFA, b NFA) bool {
	return evaluateProduct(a.Start, b.Start, a.Accept, b.Accept, Visited{})
}

func evaluateProduct(A *State, B *State, winA *State, winB *State, visited Visited) bool {

	// visited check to avoid infinite recursions
	if _, has := visited[[2]*State{A, B}]; has {
		return false
	}
	visited[[2]*State{A, B}] = struct{}{}

	aStates := flattenEpsilonTransitions(A)
	bStates := flattenEpsilonTransitions(B)

	// winCheck
	if _, successA := aStates[winA]; successA {
		if _, successB := bStates[winB]; successB {
			return true
		}
	}

	// find win further down the tree
	for nextA := range aStates {
		for nextB := range bStates {
			for t := range nextA.Transitions {
				if _, has := nextB.Transitions[t]; has && evaluateProduct(nextA.Transitions[t], nextB.Transitions[t], winA, winB, visited) {
					return true
				}
			}

		}
	}

	// no wins found
	return false
}

func flattenEpsilonTransitions(A *State) StateSet {
	out := StateSet{}
	out[A] = struct{}{}
	for {
		outLen := len(out)
		for s := range out {
			for t := range s.Epsilons {
				out[t] = struct{}{}
			}
		}
		if outLen == len(out) {
			break
		}
	}

	return out
}
