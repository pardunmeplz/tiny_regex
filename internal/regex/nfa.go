package regex

type State struct {
	// the reason we don't need a list of state is that one transition leading to multiple states is always handled by epsilons
	// in case o rune leads to multiple states, you end up having it transition to one intermediate state with many epsilons transitions instead
	Transitions map[rune]*State
	Epsilons    map[*State]struct{}
}

type NFA struct {
	Start  *State
	Accept *State
}

type Visited map[[2]*State]struct{}

func literalNfa(value rune) NFA {
	out := NFA{&State{}, &State{}}
	out.Start.Transitions[value] = out.Accept
	return out
}

func concatinationNfa(left NFA, right NFA) NFA {
	out := NFA{left.Start, right.Accept}
	left.Accept.Epsilons[right.Start] = struct{}{}
	return out
}

func alterationNfa(left NFA, right NFA) NFA {
	out := NFA{&State{}, &State{}}
	out.Start.Epsilons[left.Start] = struct{}{}
	out.Start.Epsilons[right.Start] = struct{}{}

	left.Accept.Epsilons[out.Accept] = struct{}{}
	right.Accept.Epsilons[out.Accept] = struct{}{}
	return out
}

func repeatNfa(repeat NFA) NFA {
	out := NFA{&State{}, &State{}}
	out.Start.Epsilons[out.Accept] = struct{}{}
	out.Start.Epsilons[repeat.Start] = struct{}{}

	repeat.Accept.Epsilons[repeat.Start] = struct{}{}
	repeat.Accept.Epsilons[out.Accept] = struct{}{}

	return out
}

func thompsonConstruction(node Node) NFA {
	switch node.nodeType {
	case LITERAL:
		return literalNfa(node.value)
	case CONCAT:
		return concatinationNfa(thompsonConstruction(*node.Left), thompsonConstruction(*node.Right))
	case ALT:
		return alterationNfa(thompsonConstruction(*node.Left), thompsonConstruction(*node.Right))
	case REPEAT:
		return repeatNfa(thompsonConstruction(*node.Left))
	case GROUP:
		return thompsonConstruction(*node.Left)
	}
	return NFA{}
}

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

func flattenEpsilonTransitions(A *State) map[*State]struct{} {
	out := map[*State]struct{}{}
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
