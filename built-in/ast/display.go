package ast

import (
	"fmt"
	"strings"
)

const spacesPerLevel = 2

type Displayer struct {
    schema         *ASTSchema
    SpacesPerLevel uint32
    displayers     map[NodeType]NodeDisplayer
}

type DisplayState struct {
    displayer Displayer
    sb        strings.Builder
}

type NodeDisplayer interface {
    GetNodeType() NodeType
    ToString(reference uint32) string
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
    s := DisplayState{*d, strings.Builder{}}

    Traverse(d.schema, &s, ast)

    return s.sb.String()
}

func (s *DisplayState) VisitNode(node Node, depth uint32) {
    fmt.Fprintf(&s.sb, "%s%s\n", s.displayer.getIndentation(depth), s.displayer.displayNode(node))
}

func (s *DisplayState) HandleUnexpectedEndNode(reference, depth uint32) {
    if depth == 0 {
        fmt.Fprintf(&s.sb, "EndNode %s not inside a Node\n", s.displayer.getEndNodeName(reference))
    } else {
        fmt.Fprintf(
            &s.sb,
            "%sEndNode %s in fixed childcount Node\n",
            s.displayer.getIndentation(depth),
            s.displayer.getEndNodeName(reference),
        )
    }
}

func (s *DisplayState) HandleIncorrectEndNode(reference, depth uint32) {
    fmt.Fprintf(
        &s.sb,
        "%sIncorrect EndNode %s\n",
        s.displayer.getIndentation(depth),
        s.displayer.getEndNodeName(reference),
    )
}

func (s *DisplayState) HandleMissingEndNode(depth uint32) {
    fmt.Fprintf(&s.sb, "%sMissing EndNode\n", s.displayer.getIndentation(depth))
}

func (s *DisplayState) HandleMissingChildNodes(count uint8, depth uint32) {
    fmt.Fprintf(&s.sb, "%s%d missing\n", s.displayer.getIndentation(depth), count)
}

func (d *Displayer) displayNode(node Node) string {
    if displayer, ok := d.displayers[node.Type]; ok {
        return displayer.ToString(node.Reference)
    }

    return d.schema.GetNodeTypeMetadata(node.Type).GetDebugName(node.Reference)
}

func (d *Displayer) getIndentation(depth uint32) string {
    return strings.Repeat(" ", int(d.SpacesPerLevel * depth))
}

func (d *Displayer) getEndNodeName(reference uint32) string {
    if int(reference) < len(d.schema.metadata) {
        return d.schema.GetNodeTypeMetadata(NodeType(reference)).GetDebugName(0)
    }

    return fmt.Sprintf("UNKNOWN Reference=%d", reference)
}

func (d *Displayer) DisplayDiff(before, after []Node) string {
    return ""
}

func compareNodes(a, b Node) bool {
    return a.Type == b.Type && a.Reference == b.Reference
}
