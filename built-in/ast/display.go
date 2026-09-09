package ast

import (
	"fmt"
	"io"
	"minimal/minimal-lang/built-in/messenger"
	"strings"
)

type Displayer struct {
    messenger      *messenger.Messenger
    schema         *ASTSchema
    output         io.Writer
    writeFailed    bool
    SpacesPerLevel uint
}

type cursor struct {
    ast      []Node
    position uint
}

func NewDisplayer(m *messenger.Messenger, s *ASTSchema, output io.Writer) *Displayer {
    return &Displayer{m, s, output, false, 2}
}

func (d *Displayer) Display(ast []Node) {
    d.writeFailed = false
    c := cursor{ast, 0}

    for int(c.position) != len(c.ast) {
        node := c.ast[c.position]
        c.position++

        if node.Type != EndNode {
            d.displayNode(&c, node, 0)
        } else {
            d.write("EndNode %s not inside a Node\n", d.getEndNodeName(node))
        }
    }
}

func (d *Displayer) DisplayDiff() {

}

func (d *Displayer) displayNode(c *cursor, node Node, depth int) {
    metadata := d.schema.GetNodeTypeMetadata(node.Type)
    indentationSpace := strings.Repeat(" ", int(d.SpacesPerLevel) * depth)
    debugName := metadata.GetDebugName(node.Reference)

    d.write("%s%s\n", indentationSpace, debugName)

    childCount := metadata.GetChildCount()

    if childCount == VariableChildCount {
        for int(c.position) != len(c.ast) {
            nextNode := c.ast[c.position]

            if nextNode.Type != EndNode || nextNode.Reference == uint32(node.Type) {
                c.position++

                if nextNode.Type == EndNode {
                    return
                }

                d.displayNode(c, nextNode, depth + 1)
            } else {
                d.write("%sIncorrect EndNode %s\n", indentationSpace, d.getEndNodeName(nextNode))

                return
            }
        }

        d.write("%sMissing EndNode\n", indentationSpace)

        return
    }

    childIndentationSpace := strings.Repeat(" ", int(d.SpacesPerLevel) * (depth + 1))

    for i := range childCount {
        if int(c.position) == len(c.ast) {
            d.write("%s%d missing\n", childIndentationSpace, childCount - i)

            return
        }

        nextNode := c.ast[c.position]
        c.position++

        if nextNode.Type == EndNode {
            d.write("%sEndNode %s in fixed childcount Node\n", childIndentationSpace, d.getEndNodeName(nextNode))
        } else {
            d.displayNode(c, nextNode, depth + 1)
        }
    }
}

func (d *Displayer) write(format string, args ...any) {
    if d.writeFailed {
        return
    }

    _, err := fmt.Fprintf(d.output, format, args...)

    if err != nil {
        d.writeFailed = true

        d.messenger.Send(
            messenger.Message{
                Message: "AST debugger output write failed",
                Severity: messenger.Error,
            },
        )
    }
}

func (d *Displayer) getEndNodeName(endNode Node) string {
    if int(endNode.Reference) < len(d.schema.metadata) {
        return d.schema.GetNodeTypeMetadata(NodeType(endNode.Reference)).GetDebugName(0)
    }

    return fmt.Sprintf("UNKNOWN Reference=%d", endNode.Reference)
}
