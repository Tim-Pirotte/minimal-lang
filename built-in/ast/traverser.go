package ast

import "fmt"

type TraverserSchema struct {
    s                 *ASTSchema
    UnexpectedEndNode NodeType
    IncorrectEndNode  NodeType
    MissingEndNode    NodeType
    MissingChildNodes NodeType
}

type traverser struct {
    scheme   *TraverserSchema
    visitor  Visitor
    ast      []Node
    position uint32
}

type Visitor interface {
    VisitNode(node Node, depth uint32)
}

type unexpectedEndNode struct {
    schema *ASTSchema
}

type incorrectEndNode struct {
    astSchema *ASTSchema
}

type missingEndNode struct {}
type missingChildNodes struct {}

func NewTraverserScheme(s *ASTSchema) *TraverserSchema {
    return &TraverserSchema{
        s,
        s.NewNodeType(&unexpectedEndNode{s}),
        s.NewNodeType(&incorrectEndNode{s}),
        s.NewNodeType(&missingEndNode{}),
        s.NewNodeType(&missingChildNodes{}),
    }
}

// Traverses the AST and calls VisitNode for every node.
// There are 4 instances where VisitNode is called with error nodes:
//   - Unexpected EndNode:  UnexpectedEndNode(Reference=EndNode.Reference)
//   - Incorrect EndNode:   IncorrectEndNode(Reference=EndNode.Reference)
//   - Missing EndNode:     MissingEndNode()
//   - Missing Child nodes: MissingChildNodes(Reference=amountMissing)
func (scheme *TraverserSchema) Traverse(visitor Visitor, ast []Node) {
    t := traverser{scheme, visitor, ast, 0}

    for int(t.position) != len(ast) {
        node := ast[t.position]
        t.position++

        if node.Type != EndNode {
            t.traverseNode(node, 0)
        } else {
            t.traverseNode(Node{scheme.UnexpectedEndNode, node.Reference}, 0)
        }
    }
}

func (t *traverser) traverseNode(node Node, depth uint32) {
    t.visitor.VisitNode(node, depth)

    childCount := t.scheme.s.GetNodeTypeMetadata(node.Type).GetChildCount()

    if childCount == VariableChildCount {
        t.traverseVariable(node, depth)
    } else {
        t.traverseFixed(childCount, depth)
    }
}

func (t *traverser) traverseVariable(parent Node, depth uint32) {
    for int(t.position) != len(t.ast) {
        child := t.ast[t.position]

        if child.Type == EndNode {
            if child.Reference == uint32(parent.Type) {
                t.position++

                return
            } else {
                t.visitor.VisitNode(Node{t.scheme.IncorrectEndNode, child.Reference}, depth)

                return
            }
        }

        t.position++
        t.traverseNode(child, depth + 1)
    }

    t.visitor.VisitNode(Node{t.scheme.MissingEndNode, 0}, depth)
}

func (t *traverser) traverseFixed(childCount uint8, depth uint32) {
    for i := range childCount {
        if int(t.position) == len(t.ast) {
            t.visitor.VisitNode(Node{t.scheme.MissingChildNodes, uint32(childCount - i)}, depth + 1)

            return
        }

        child := t.ast[t.position]
        t.position++

        if child.Type != EndNode {
            t.traverseNode(child, depth + 1)
        } else {
            t.visitor.VisitNode(Node{t.scheme.UnexpectedEndNode, child.Reference}, depth + 1)
        }
    }
}

func (*unexpectedEndNode) GetDisplayName(reference uint32) string {
    panic("GetDisplayName called on unexpectedEndNode")
}

func (u *unexpectedEndNode) GetDebugName(reference uint32) string {
    return fmt.Sprintf(
        "Unexpected EndNode Reference=%d%s",
        reference,
        convertReferenceToNodeTypeName(u.schema, reference),
    )
}

func (*unexpectedEndNode) GetChildCount() uint8 {
    return 0
}

func (*incorrectEndNode) GetDisplayName(reference uint32) string {
    panic("GetDisplayName called on incorrectEndNode")
}

func (i *incorrectEndNode) GetDebugName(reference uint32) string {
    return fmt.Sprintf(
        "Incorrect EndNode Reference=%d%s",
        reference,
        convertReferenceToNodeTypeName(i.astSchema, reference),
    )
}

func (*incorrectEndNode) GetChildCount() uint8 {
    return 0
}

func (*missingEndNode) GetDisplayName(reference uint32) string {
    panic("GetDisplayName called on missingEndNode")
}

func (*missingEndNode) GetDebugName(reference uint32) string {
    return "Missing EndNode"
}

func (*missingEndNode) GetChildCount() uint8 {
    return 0
}

func (*missingChildNodes) GetDisplayName(reference uint32) string {
    panic("GetDisplayName called on missingChildNodes")
}

func (*missingChildNodes) GetDebugName(reference uint32) string {
    return fmt.Sprintf("%d missing", reference)
}

func (*missingChildNodes) GetChildCount() uint8 {
    return 0
}

func convertReferenceToNodeTypeName(astSchema *ASTSchema, reference uint32) string {
    if int(reference) < len(astSchema.metadata) {
        return fmt.Sprintf(
            " NodeType(Reference)=%s",
            astSchema.GetNodeTypeMetadata(NodeType(reference)).GetDebugName(0),
        )
    }

    return ""
}
