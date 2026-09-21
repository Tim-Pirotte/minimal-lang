package ast

import (
	"fmt"
	"minimal/minimal-lang/built-in/ansi"
	"minimal/minimal-lang/built-in/diff"
	"strings"
)

const spacesPerLevel = 2

type Displayer struct {
    schema         *ASTSchema
    SpacesPerLevel uint32
    OutputANSI     bool
    displayers     map[NodeType]NodeDisplayer
}

type DisplayState struct {
    displayer Displayer
    sb        strings.Builder
}

type tagTraverser struct {
    result []taggedNode
}

type taggedNode struct {
    node  Node
    depth uint32
    valid bool
}

type NodeDisplayer interface {
    GetNodeType() NodeType
    ToString(reference uint32) string
}

func NewDisplayer(s *ASTSchema) *Displayer {
    return &Displayer{s, spacesPerLevel, true, map[NodeType]NodeDisplayer{}}
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

func (d *Displayer) DisplayDiff(before, after []Node) string {
    t := tagTraverser{[]taggedNode{}}

    Traverse(d.schema, &t, before)
    taggedBefore := t.result
    t.result = []taggedNode{}

    Traverse(d.schema, &t, after)
    taggedAfter := t.result

    astDiff := diff.GetDiff(taggedBefore, taggedAfter, compareNodes)

    sb := strings.Builder{}

    for _, part := range astDiff {
        fmt.Fprint(&sb, d.getIndentation(part.Value.depth))

        switch part.Type {
        case diff.Equal:
            fmt.Fprint(&sb, "  ")
        case diff.Insert:
            prefix := "+ "

            if d.OutputANSI {
                prefix = ansi.RGB{R: 121, G: 245, B: 5}.ToString() + prefix + ansi.Reset
            }

            fmt.Fprint(&sb, prefix)
        case diff.Delete:
            prefix := "- "

            if d.OutputANSI {
                prefix = ansi.RGB{R: 245, G: 5, B: 61}.ToString() + prefix + ansi.Reset
            }

            fmt.Fprint(&sb, prefix)
        }

        fmt.Fprintf(&sb, "%s\n", d.displayNode(part.Value.node))
    }

    return sb.String()
}

func compareNodes(a, b taggedNode) bool {
    return a.node.Type == b.node.Type && a.node.Reference == b.node.Reference
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

func (t *tagTraverser) VisitNode(node Node, depth uint32) {
    t.result = append(t.result, taggedNode{node, depth, true})
}

func (t *tagTraverser) HandleUnexpectedEndNode(reference, depth uint32) {
    t.result = append(t.result, taggedNode{Node{EndNode, reference}, depth, false})
}

func (t *tagTraverser) HandleIncorrectEndNode(reference, depth uint32) {
    t.result = append(t.result, taggedNode{Node{EndNode, reference}, depth, false})
}

// TODO insert a proxy node for these two cases
func (t *tagTraverser) HandleMissingEndNode(depth uint32) {}
func (t *tagTraverser) HandleMissingChildNodes(count uint8, depth uint32) {}
