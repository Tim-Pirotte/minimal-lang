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
            d.displayNodeTree(&c, node, 0)
        } else {
            d.showEndNodeOutsideNode(node)
        }
    }
}

func (d *Displayer) DisplayDiff() {

}

func (d *Displayer) displayNodeTree(c *cursor, node Node, depth int) {
    indentSpace := strings.Repeat(" ", int(d.SpacesPerLevel) * depth)

    d.write("%s", indentSpace)
    d.displayNode(node)
    d.write("\n")

    metadata := d.schema.GetNodeTypeMetadata(node.Type)
    childCount := metadata.GetChildCount()

    if childCount == VariableChildCount {
        d.displayVariableChildCount(c, node, indentSpace, depth)
    } else {
        d.displayFixedChildCount(c, int(childCount), depth)
    }
}

func (d *Displayer) displayVariableChildCount(c *cursor, node Node, indentSpace string, depth int) {
    for int(c.position) != len(c.ast) {
        nextNode := c.ast[c.position]

        if nextNode.Type == EndNode {
            if nextNode.Reference != uint32(node.Type) {
                d.showIncorrectEndNode(indentSpace, nextNode)

                return
            }

            c.position++

            return
        }

        c.position++
        d.displayNodeTree(c, nextNode, depth + 1)
    }

    d.showMissingEndNode(indentSpace)
}

func (d *Displayer) displayFixedChildCount(c *cursor, childCount, depth int) {
    childIndentSpace := strings.Repeat(" ", int(d.SpacesPerLevel) * (depth + 1))

    for i := range childCount {
        if int(c.position) == len(c.ast) {
            d.showMissingChildNode(childIndentSpace, childCount - i)

            return
        }

        nextNode := c.ast[c.position]
        c.position++

        if nextNode.Type == EndNode {
            d.showEndNodeInFixedChildCountNode(childIndentSpace, nextNode)
        } else {
            d.displayNodeTree(c, nextNode, depth + 1)
        }
    }
}

func (d *Displayer) displayNode(node Node) {
    metadata := d.schema.GetNodeTypeMetadata(node.Type)
    debugName := metadata.GetDebugName(node.Reference)

    d.write("%s", debugName)
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

func (d *Displayer) showEndNodeOutsideNode(node Node) {
    d.write("EndNode %s not inside a Node\n", d.getEndNodeName(node))
}

func (d *Displayer) showIncorrectEndNode(indentationSpace string, node Node) {
    d.write("%sIncorrect EndNode %s\n", indentationSpace, d.getEndNodeName(node))
}

func (d *Displayer) showMissingEndNode(indentationSpace string) {
    d.write("%sMissing EndNode\n", indentationSpace)
}

func (d *Displayer) showMissingChildNode(childIndentationSpace string, number int) {
    d.write("%s%d missing\n", childIndentationSpace, number)
}

func (d *Displayer) showEndNodeInFixedChildCountNode(indentationSpace string, node Node) {
    d.write(
        "%sEndNode %s in fixed childcount Node\n",
        indentationSpace,
        d.getEndNodeName(node),
    )
}
