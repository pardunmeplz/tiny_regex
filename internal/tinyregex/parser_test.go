package tinyregex

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

var positiveTests = []struct {
	regex  string
	result *Node
}{
	{"AB|", &Node{ALT, ' ', &Node{CONCAT, ' ', &Node{LITERAL, 'A', nil, nil}, &Node{LITERAL, 'B', nil, nil}}, &Node{EPSILON, ' ', nil, nil}}},
	{"()", &Node{GROUP, ' ', &Node{EPSILON, ' ', nil, nil}, nil}},
	{"|BA", &Node{ALT, ' ', &Node{EPSILON, ' ', nil, nil}, &Node{CONCAT, ' ', &Node{LITERAL, 'B', nil, nil}, &Node{LITERAL, 'A', nil, nil}}}},
	{"AB|C", &Node{ALT, ' ', &Node{CONCAT, ' ', &Node{LITERAL, 'A', nil, nil}, &Node{LITERAL, 'B', nil, nil}}, &Node{LITERAL, 'C', nil, nil}}},
	{"A(B|C)", &Node{CONCAT, ' ', &Node{LITERAL, 'A', nil, nil}, &Node{GROUP, ' ', &Node{ALT, ' ', &Node{LITERAL, 'B', nil, nil}, &Node{LITERAL, 'C', nil, nil}}, nil}}},
	{"(AB|C)", &Node{GROUP, ' ', &Node{ALT, ' ', &Node{CONCAT, ' ', &Node{LITERAL, 'A', nil, nil}, &Node{LITERAL, 'B', nil, nil}}, &Node{LITERAL, 'C', nil, nil}}, nil}},
	{"AB*|C", &Node{ALT, ' ', &Node{CONCAT, ' ', &Node{LITERAL, 'A', nil, nil}, &Node{REPEAT, ' ', &Node{LITERAL, 'B', nil, nil}, nil}}, &Node{LITERAL, 'C', nil, nil}}},
	{"(A|B)*|C", &Node{ALT, ' ', &Node{REPEAT, ' ', &Node{GROUP, ' ', &Node{ALT, ' ', &Node{LITERAL, 'A', nil, nil}, &Node{LITERAL, 'B', nil, nil}}, nil}, nil}, &Node{LITERAL, 'C', nil, nil}}},
	{"(A|B)*|", &Node{ALT, ' ', &Node{REPEAT, ' ', &Node{GROUP, ' ', &Node{ALT, ' ', &Node{LITERAL, 'A', nil, nil}, &Node{LITERAL, 'B', nil, nil}}, nil}, nil}, &Node{EPSILON, ' ', nil, nil}}},
}

var negativeTests = []struct {
	regex string
	errs  []string
}{
	{")", []string{fmt.Sprintf(UNEXPECTED_CHAR, ")")}},
	{"(", []string{MISSING_GROUP_END}},
	{"c(a|b", []string{MISSING_GROUP_END}},
	{"c)a|b", []string{fmt.Sprintf(UNEXPECTED_CHAR, ")")}},
}

func TestPositive(t *testing.T) {
	for _, test := range positiveTests {
		out, err := parse(test.regex)
		if len(err) > 0 {
			t.Errorf("%s has invalid errors %s", test.regex, err)
		}
		if !matchNodes(test.result, &out) {
			t.Errorf("%s has invalid output %s expected %s", test.regex, printNode(&out), printNode(test.result))
		}
	}
}

func TestNegative(t *testing.T) {
	for _, test := range negativeTests {
		_, err := parse(test.regex)
		if !slices.EqualFunc(err, test.errs, func(a string, b string) bool { return a == b }) {
			t.Errorf("%s has invalid errors %s", test.regex, err)
		}
	}
}

func matchNodes(A *Node, B *Node) bool {
	if A == nil || B == nil {
		return A == B
	}
	return A.value == B.value && A.nodeType == B.nodeType && matchNodes(A.Left, B.Left) && matchNodes(A.Right, B.Right)
}

func printNode(A *Node) string {
	if A == nil {
		return "{}"
	}
	out := strings.Builder{}
	out.WriteString("{")
	switch A.nodeType {
	case LITERAL:
		out.WriteString("LITERAL")
		out.WriteString(",")
		out.WriteRune(A.value)
	case CONCAT:
		out.WriteString("CONCAT")
	case ALT:
		out.WriteString("ALT")
	case REPEAT:
		out.WriteString("REPEAT")
	case GROUP:
		out.WriteString("GROUP")
	}

	out.WriteString(",")
	out.WriteString(printNode(A.Left))
	out.WriteString(",")
	out.WriteString(printNode(A.Right))
	out.WriteString("}")

	return out.String()
}
