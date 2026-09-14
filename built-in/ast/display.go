package ast

import (
	"fmt"
	"strings"
)

const spacesPerLevel = 2

type Displayer struct {
    schema         *ASTSchema
    SpacesPerLevel uint
    displayers     map[NodeType]NodeDisplayer
}

type NodeDisplayer interface {
    GetNodeType() NodeType
    ToString(reference uint32) string
}

type displayState struct {
    ast      []Node
    position uint
    sb       strings.Builder
}

func NewDisplayer(s *ASTSchema) *Displayer {
    return &Displayer{s, spacesPerLevel, map[NodeType]NodeDisplayer{}}
}

func (d *Displayer) AddDisplayer(n NodeDisplayer) bool {
    nodeType := n.GetNodeType()

    if _, ok := d.displayers[nodeType]; ok {
        return false
    }

    d.displayers[nodeType] = n

    return true
}

func (d *Displayer) Display(ast []Node) string {
    s := displayState{ast, 0, strings.Builder{}}

    for int(s.position) != len(s.ast) {
        node := s.ast[s.position]
        s.position++

        if node.Type != EndNode {
            d.displayNodeTree(&s, node, 0)
        } else {
            d.showEndNodeOutsideNode(&s, node)
        }
    }

    return s.sb.String()
}

func (d *Displayer) DisplayDiff() {

}

func (d *Displayer) displayNodeTree(s *displayState, node Node, depth int) {
    indentSpace := strings.Repeat(" ", int(d.SpacesPerLevel) * depth)

    fmt.Fprintf(&s.sb, "%s%s\n", indentSpace, d.displayNode(node))

    metadata := d.schema.GetNodeTypeMetadata(node.Type)
    childCount := metadata.GetChildCount()

    if childCount == VariableChildCount {
        d.displayVariableChildCount(s, node, indentSpace, depth)
    } else {
        d.displayFixedChildCount(s, int(childCount), depth)
    }
}

func (d *Displayer) displayVariableChildCount(s *displayState, node Node, indentSpace string, depth int) {
    for int(s.position) != len(s.ast) {
        nextNode := s.ast[s.position]

        if nextNode.Type == EndNode {
            if nextNode.Reference != uint32(node.Type) {
                d.showIncorrectEndNode(s, indentSpace, nextNode)

                return
            }

            s.position++

            return
        }

        s.position++
        d.displayNodeTree(s, nextNode, depth + 1)
    }

    d.showMissingEndNode(s, indentSpace)
}

func (d *Displayer) displayFixedChildCount(s *displayState, childCount, depth int) {
    childIndentSpace := strings.Repeat(" ", int(d.SpacesPerLevel) * (depth + 1))

    for i := range childCount {
        if int(s.position) == len(s.ast) {
            d.showMissingChildNode(s, childIndentSpace, childCount - i)

            return
        }

        nextNode := s.ast[s.position]
        s.position++

        if nextNode.Type == EndNode {
            d.showEndNodeInFixedCountNode(s, childIndentSpace, nextNode)
        } else {
            d.displayNodeTree(s, nextNode, depth + 1)
        }
    }
}

func (d *Displayer) displayNode(node Node) string {
    if displayer, ok := d.displayers[node.Type]; ok {
        return displayer.ToString(node.Reference)
    }

    return d.schema.GetNodeTypeMetadata(node.Type).GetDebugName(node.Reference)
}

func (d *Displayer) getEndNodeName(endNode Node) string {
    if int(endNode.Reference) < len(d.schema.metadata) {
        return d.schema.GetNodeTypeMetadata(NodeType(endNode.Reference)).GetDebugName(0)
    }

    return fmt.Sprintf("UNKNOWN Reference=%d", endNode.Reference)
}

func (d *Displayer) showEndNodeOutsideNode(s *displayState, node Node) {
    fmt.Fprintf(&s.sb, "EndNode %s not inside a Node\n", d.getEndNodeName(node))
}

func (d *Displayer) showIncorrectEndNode(s *displayState, indentationSpace string, node Node) {
    fmt.Fprintf(&s.sb, "%sIncorrect EndNode %s\n", indentationSpace, d.getEndNodeName(node))
}

func (d *Displayer) showMissingEndNode(s *displayState, indentationSpace string) {
    fmt.Fprintf(&s.sb, "%sMissing EndNode\n", indentationSpace)
}

func (d *Displayer) showMissingChildNode(s *displayState, childIndentationSpace string, number int) {
    fmt.Fprintf(&s.sb, "%s%d missing\n", childIndentationSpace, number)
}

func (d *Displayer) showEndNodeInFixedCountNode(s *displayState, indentationSpace string, node Node) {
    fmt.Fprintf(&s.sb, "%sEndNode %s in fixed childcount Node\n", indentationSpace, d.getEndNodeName(node))
}
