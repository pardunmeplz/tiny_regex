package main

import (
	"fmt"
	rx "regex/internal/regexEval"
)

func main() {
	fmt.Print(rx.HasUnion("ab", "b|ab*"))
}
