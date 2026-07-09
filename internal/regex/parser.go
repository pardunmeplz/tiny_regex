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

func parse(input string) Node {
	parserState := &ParserState{input, 0, []string{}}

	out := regex(parserState)
	if pending(parserState) {
		addError(parserState, "Invalid characters at end of regex")
	}
	return out
}

func regex(parserState *ParserState) Node {
	return alternation(parserState)
}

func alternation(parserState *ParserState) Node {
	left := concatination(parserState)
	for pending(parserState) && peek(parserState) == '|' {
		advance(parserState)
		right := concatination(parserState)
		if left.nodeType == ALT {
			left.children = append(left.children, right)
		} else {
			left = Node{ALT, ' ', []Node{left, right}}
		}
	}
	return left
}
func concatination(parserState *ParserState) Node {
	left := repetition(parserState)
	for pending(parserState) && peek(parserState) != '|' && peek(parserState) != ')' {
		right := repetition(parserState)
		if left.nodeType == CONCAT {
			left.children = append(left.children, right)
		} else {
			left = Node{CONCAT, ' ', []Node{left, right}}
		}
	}
	return left
}

func repetition(parserState *ParserState) Node {
	left := primary(parserState)
	for pending(parserState) && peek(parserState) == '*' {
		advance(parserState)
		left = Node{REPEAT, ' ', []Node{left}}
	}
	return left
}

func primary(parserState *ParserState) Node {
	char := advance(parserState)
	if char == '(' {
		grp := Node{GROUP, ' ', []Node{regex(parserState)}}
		consume(parserState, ')', "Expected closing paranthesis ')'")
		return grp
	}

	if char == ')' || char == '|' || char == '*' {
		addError(parserState, "Unexpected character "+string(char))
		return Node{}
	}

	return Node{LITERAL, char, nil}
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
