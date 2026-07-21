package regex

/*
Minimal Regex Language Grammero

Regex          ::= Alternation

Alternation    ::= Concatenation ( "|" Concatenation )*

Concatenation  ::= Repetition+

Repetition     ::= Primary ( "*" )*

Primary        ::= Literal
                 | "(" Regex ")"

Literal        ::= <any character except '(', ')', '|', '*'>

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
		addError(parserState, "Invalid characters at end of regex")
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
	char := advance(parserState)
	if char == '(' {
		left := regex(parserState)
		grp := Node{GROUP, ' ', &left, nil}
		consume(parserState, ')', "Expected closing paranthesis ')'")
		return grp
	}

	if char == ')' || char == '|' || char == '*' {
		addError(parserState, "Unexpected character "+string(char))
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
