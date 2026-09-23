package ast

import (
	"fmt"
	"minimal/minimal-lang/built-in/ansi"
	"minimal/minimal-lang/built-in/diff"
	"strings"
)

const (
    spacesPerLevel = 2
    outputANSI = true
)

type Displayer struct {
    astSchema       *ASTSchema
    traverserSchema TraverserSchema
    SpacesPerLevel  uint32
    OutputANSI      bool
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
}

func NewDisplayer(s *ASTSchema) *Displayer {
    return &Displayer{s, *NewTraverserScheme(s), spacesPerLevel, outputANSI}
}

func (d *Displayer) Display(ast []Node) string {
    s := DisplayState{*d, strings.Builder{}}
    d.traverserSchema.Traverse(&s, ast)

    return s.sb.String()
}

func (d *Displayer) DisplayDiff(before, after []Node) string {
    t := tagTraverser{[]taggedNode{}}

    d.traverserSchema.Traverse(&t, before)
    taggedBefore := t.result

    t.result = []taggedNode{}
    d.traverserSchema.Traverse(&t, after)
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

func (d *Displayer) displayNode(node Node) string {
    return d.astSchema.GetNodeTypeMetadata(node.Type).GetDebugName(node.Reference)
}

func (d *Displayer) getIndentation(depth uint32) string {
    return strings.Repeat(" ", int(d.SpacesPerLevel * depth))
}

func (t *tagTraverser) VisitNode(node Node, depth uint32) {
    t.result = append(t.result, taggedNode{node, depth})
}
