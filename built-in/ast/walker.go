package ast

import "fmt"

type WalkerSchema struct {
	astSchema     *ASTSchema
	checkers      map[NodeType]NodeChecker
	depthTrackers []ScopeTracker
}

type Walker struct {
	schema   WalkerSchema
	ast      []Node
	position uint32
}

type NodeChecker interface {
	GetNodeType() NodeType
	Check(w *Walker)
}

type ScopeTracker interface {
	Enter()
	Exit()
}

func NewWalkerSchema(schema *ASTSchema) *WalkerSchema {
	return &WalkerSchema{schema, map[NodeType]NodeChecker{}, []ScopeTracker{}}
}

func (w *WalkerSchema) AddNodeChecker(nodeChecker NodeChecker) bool {
	nodeType := nodeChecker.GetNodeType()

	if _, ok := w.checkers[nodeType]; ok {
		return false
	}

	w.checkers[nodeType] = nodeChecker

	return true
}

func (w *WalkerSchema) NewWalker(ast []Node) *Walker {
	return &Walker{*w, ast, 0}
}

func (w *Walker) Next() {
	if w.position >= uint32(len(w.ast)) {
		panic("Next called when at the end of the AST")
	}

	node := w.ast[w.position]
	w.position++

	if checker, ok := w.schema.checkers[node.Type]; ok {
		for _, t := range w.schema.depthTrackers {
			t.Enter()
		}

		checker.Check(w)

		for _, t := range w.schema.depthTrackers {
			t.Exit()
		}

		return
	}

	metadata := w.schema.astSchema.GetNodeTypeMetadata(node.Type)

	if metadata.GetChildCount() != VariableChildCount {
		w.position += uint32(metadata.GetChildCount())

		return
	}

	for w.position != uint32(len(w.ast)) {
		if w.ast[w.position].Type == EndNode {
			return
		}

		w.position++
	}

	panic(fmt.Sprintf("Node %+v does not have an EndNode", node))
}
