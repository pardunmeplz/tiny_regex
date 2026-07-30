package regex

type State struct {
	// the reason we don't need a list of state is that one transition leading to multiple states is always handled by epsilons
	// in case o rune leads to multiple states, you end up having it transition to one intermediate state with many epsilons transitions instead
	Transitions TransitionSet
	Epsilons    StateSet
}

type StateSet map[*State]struct{}
type TransitionSet map[rune]*State

type NFA struct {
	Start  *State
	Accept *State
}

func literalNfa(value rune) NFA {
	out := NFA{&State{TransitionSet{}, StateSet{}}, &State{TransitionSet{}, StateSet{}}}
	out.Start.Transitions[value] = out.Accept
	return out
}

func concatinationNfa(left NFA, right NFA) NFA {
	out := NFA{left.Start, right.Accept}
	left.Accept.Epsilons[right.Start] = struct{}{}
	return out
}

func alterationNfa(left NFA, right NFA) NFA {
	out := NFA{&State{TransitionSet{}, StateSet{}}, &State{TransitionSet{}, StateSet{}}}
	out.Start.Epsilons[left.Start] = struct{}{}
	out.Start.Epsilons[right.Start] = struct{}{}

	left.Accept.Epsilons[out.Accept] = struct{}{}
	right.Accept.Epsilons[out.Accept] = struct{}{}
	return out
}

func repeatNfa(repeat NFA) NFA {
	out := NFA{&State{TransitionSet{}, StateSet{}}, &State{TransitionSet{}, StateSet{}}}
	out.Start.Epsilons[out.Accept] = struct{}{}
	out.Start.Epsilons[repeat.Start] = struct{}{}

	repeat.Accept.Epsilons[repeat.Start] = struct{}{}
	repeat.Accept.Epsilons[out.Accept] = struct{}{}

	return out
}

func blankEpsilonNfa() NFA {
	out := NFA{&State{TransitionSet{}, StateSet{}}, &State{TransitionSet{}, StateSet{}}}
	out.Start.Epsilons[out.Accept] = struct{}{}
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
	case EPSILON:
		return blankEpsilonNfa()
	}

	return NFA{}
}
