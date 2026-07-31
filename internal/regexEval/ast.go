package regex

type NodeType int

const (
	LITERAL NodeType = 1
	CONCAT  NodeType = 2
	ALT     NodeType = 3
	REPEAT  NodeType = 4
	GROUP   NodeType = 5
	EPSILON NodeType = 6
)

type Node struct {
	nodeType NodeType
	value    rune
	Left     *Node
	Right    *Node
}

const (
	MISSING_GROUP_END = "Expected closing paranthesis ')'"
	UNEXPECTED_CHAR   = "Unexpected character %s"
)
