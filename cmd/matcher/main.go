package main

import (
	"fmt"
	"os"

	rx "regex/internal/regexEval"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <regex-a> <regex-b>\n", os.Args[0])
		os.Exit(2)
	}

	ok, errs := rx.HasUnion(os.Args[1], os.Args[2])
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, e)
		}
		os.Exit(1)
	}

	if ok {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}
