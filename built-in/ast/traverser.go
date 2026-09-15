package ast

type traverser struct {
    s        *ASTSchema
    hooks    TraverserHooks
    ast      []Node
    position uint32
}

type TraverserHooks interface {
    VisitNode(node Node, depth uint32)
    HandleUnexpectedEndNode(reference, depth uint32)
    HandleIncorrectEndNode(reference, depth uint32)
    HandleMissingEndNode(depth uint32)
    HandleMissingChildNodes(count uint8, depth uint32)
}

func Traverse(s *ASTSchema, hooks TraverserHooks, ast []Node) {
    t := traverser{s, hooks, ast, 0}

    for int(t.position) != len(ast) {
        node := ast[t.position]
        t.position++

        if node.Type != EndNode {
            t.traverseNode(node, 0)
        } else {
            hooks.HandleUnexpectedEndNode(node.Reference, 0)
        }
    }
}

func (t *traverser) traverseNode(node Node, depth uint32) {
    t.hooks.VisitNode(node, depth)

    childCount := t.s.GetNodeTypeMetadata(node.Type).GetChildCount()

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
                t.hooks.HandleIncorrectEndNode(child.Reference, depth)

                return
            }
        }

        t.position++
        t.traverseNode(child, depth + 1)
    }

    t.hooks.HandleMissingEndNode(depth)
}

func (t *traverser) traverseFixed(childCount uint8, depth uint32) {
    for i := range childCount {
        if int(t.position) == len(t.ast) {
            t.hooks.HandleMissingChildNodes(childCount - i, depth + 1)

            return
        }

        child := t.ast[t.position]
        t.position++

        if child.Type != EndNode {
            t.traverseNode(child, depth + 1)
        } else {
            t.hooks.HandleUnexpectedEndNode(child.Reference, depth + 1)
        }
    }
}
