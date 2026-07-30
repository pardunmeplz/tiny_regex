package regex

import (
	"maps"
)

func eval(string string, nfa NFA) bool {
	currState := StateSet{}
	currState[nfa.Start] = struct{}{}

	for _, ch := range string {
		if len(currState) == 0 {
			return false
		}
		currState = flattenEpsilonsInSet(currState)
		currState = transition(currState, ch)
	}

	currState = flattenEpsilonsInSet(currState)
	_, win := currState[nfa.Accept]

	return win

}

func flattenEpsilonsInSet(set StateSet) StateSet {
	output := StateSet{}
	for state := range set {
		res := flattenEpsilonTransitions(state)
		// bruh
		maps.Copy(output, res)
	}
	return output
}

func transition(set StateSet, ch rune) StateSet {
	output := StateSet{}
	for state := range set {
		res, has := state.Transitions[ch]
		if has {
			output[res] = struct{}{}
		}
	}

	return output

}
