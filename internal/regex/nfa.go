package regex

type Transition struct {
	// its a poniter because epsilon transitions will have In as nil
	In  *rune
	Out *State
}

type State struct {
	Transitions []Transition
}

type NFA struct {
	Start  *State
	Accept *State
}

func literalNfa(value rune) NFA {
	out := NFA{&State{}, &State{}}
	out.Start.Transitions = []Transition{{&value, out.Accept}}
	return out
}

func concatinationNfa(left NFA, right NFA) NFA {
	out := NFA{left.Start, right.Accept}
	left.Accept.Transitions = append(left.Accept.Transitions, Transition{nil, right.Start})
	return out
}

func alterationNfa(left NFA, right NFA) NFA {
	out := NFA{&State{[]Transition{{nil, left.Start}, {nil, right.Start}}}, &State{}}
	left.Accept.Transitions = append(left.Accept.Transitions, Transition{nil, out.Accept})
	right.Accept.Transitions = append(right.Accept.Transitions, Transition{nil, out.Accept})
	return out
}

func repeatNfa(repeat NFA) NFA {
	out := NFA{&State{}, &State{}}
	out.Start.Transitions = []Transition{{nil, out.Accept}, {nil, repeat.Start}}
	repeat.Accept.Transitions = append(repeat.Accept.Transitions, Transition{nil, repeat.Start}, Transition{nil, out.Accept})

	return out
}
