package regex

type NodeType int

const (
	LITERAL NodeType = 1
	CONCAT  NodeType = 2
	ALT     NodeType = 3
	REPEAT  NodeType = 4
	GROUP   NodeType = 5
)

type Node struct {
	nodeType NodeType
	value    rune
	children []Node
}
