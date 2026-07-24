package regex

import "fmt"

/*
Minimal Regex Language Grammero

Regex          ::= Alternation

Alternation    ::= Concatenation ( "|" Concatenation )*

Concatenation  ::= Repetition+

Repetition     ::= Primary ( "*" )*

Primary        ::= Literal
                 | "(" Regex ")"

Literal        ::= <any character except '(', ')', '|', '*'> | Epsilon

Epsilon        ::= ''

*/

type ParserState struct {
	input string
	index int
	error []string
}

func parse(input string) (Node, []string) {
	parserState := &ParserState{input, 0, []string{}}

	out := regex(parserState)
	if pending(parserState) {
		addError(parserState, EXPECTED_EOF)
	}
	return out, parserState.error
}

func regex(parserState *ParserState) Node {
	return alternation(parserState)
}

func alternation(parserState *ParserState) Node {
	out := concatination(parserState)
	for pending(parserState) && peek(parserState) == '|' {
		advance(parserState)
		right := concatination(parserState)
		left := out
		out = Node{ALT, ' ', &left, &right}
	}
	return out
}
func concatination(parserState *ParserState) Node {
	out := repetition(parserState)
	for pending(parserState) && peek(parserState) != '|' && peek(parserState) != ')' {
		right := repetition(parserState)
		left := out
		out = Node{CONCAT, ' ', &left, &right}
	}
	return out
}

func repetition(parserState *ParserState) Node {
	out := primary(parserState)
	for pending(parserState) && peek(parserState) == '*' {
		advance(parserState)
		left := out
		out = Node{REPEAT, ' ', &left, nil}
	}
	return out
}

func primary(parserState *ParserState) Node {

	if eof(parserState) || peekMatch(parserState, '|') {
		return Node{EPSILON, ' ', nil, nil}
	}

	char := advance(parserState)
	if char == '(' {
		grp := Node{GROUP, ' ', nil, nil}
		if peekMatch(parserState, ')') {
			grp.Left = &Node{EPSILON, ' ', nil, nil}
		} else {
			left := regex(parserState)
			grp.Left = &left
		}

		consume(parserState, ')', MISSING_GROUP_END)
		return grp
	}

	if char == ')' || char == '|' || char == '*' {
		addError(parserState, fmt.Sprintf(UNEXPECTED_CHAR, string(char)))
		return Node{}
	}

	return Node{LITERAL, char, nil, nil}
}

func peek(parserState *ParserState) rune {
	return rune(parserState.input[parserState.index])
}

func advance(parserState *ParserState) rune {
	parserState.index++
	return rune(parserState.input[parserState.index-1])
}

func eof(parserState *ParserState) bool {
	return parserState.index >= len(parserState.input)
}

func pending(parserState *ParserState) bool {
	return parserState.index < len(parserState.input)
}

func addError(parserState *ParserState, message string) {
	parserState.error = append(parserState.error, message)
}

func consume(parserState *ParserState, ch rune, messsage string) {
	if peek(parserState) != ch {
		addError(parserState, messsage)
	} else {
		advance(parserState)
	}
}

func peekMatch(parserState *ParserState, ch rune) bool {
	return ch == peek(parserState)
}
