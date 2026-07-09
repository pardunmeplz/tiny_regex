package regex

type Transition struct {
	In  *rune
	Out *State
}

type State struct {
	Transitions []Transition
}

type NFA struct {
	Start  []*State
	Accept []*State
}
