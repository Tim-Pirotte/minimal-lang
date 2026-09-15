package ast

import (
	"fmt"
	"minimal/minimal-lang/built-in/diff"
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
    ast        []Node
    position   uint
    sb         strings.Builder
    prefixHook prefixHook
}

type prefixHook interface {
    getPrefix(*displayState) string
}

type diffPrefix struct {
    nodeDiff []diff.DiffPart[Node]
}

func (d *diffPrefix) getPrefix(s *displayState) string {
    part := d.nodeDiff[s.position]

    prefix := "  "

    switch part.Type {
    case diff.Insert:
        prefix = "+ "
    case diff.Delete:
        prefix = "- "
    }

    return prefix
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
    return d.display(ast, nil)
}

func (d *Displayer) DisplayDiff(before, after []Node) string {
    nodeDiff := diff.GetDiff(before, after, compareNodes)
    ast := make([]Node, len(nodeDiff))

    for i, part := range nodeDiff {
        ast[i] = part.Value
    }

    return d.display(ast, &diffPrefix{nodeDiff})
}

func compareNodes(a, b Node) bool {
    return a.Type == b.Type && a.Reference == b.Reference
}

func (d *Displayer) display(ast []Node, prefixHook prefixHook) string {
    s := displayState{ast, 0, strings.Builder{}, prefixHook}

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

func (d *Displayer) displayNodeTree(s *displayState, node Node, depth int) {
    prefix := ""

    if s.prefixHook != nil {
        prefix = s.prefixHook.getPrefix(s)
    }

    fmt.Fprintf(&s.sb, "%s%s%s\n", d.getIndentation(depth), prefix, d.displayNode(node))

    metadata := d.schema.GetNodeTypeMetadata(node.Type)
    childCount := metadata.GetChildCount()

    if childCount == VariableChildCount {
        d.displayVariableChildCount(s, node, depth)
    } else {
        d.displayFixedChildCount(s, int(childCount), depth)
    }
}

func (d *Displayer) displayVariableChildCount(s *displayState, node Node, depth int) {
    for int(s.position) != len(s.ast) {
        nextNode := s.ast[s.position]

        if nextNode.Type == EndNode {
            if nextNode.Reference != uint32(node.Type) {
                d.showIncorrectEndNode(s, depth, nextNode)

                return
            }

            s.position++

            return
        }

        s.position++
        d.displayNodeTree(s, nextNode, depth + 1)
    }

    d.showMissingEndNode(s, depth)
}

func (d *Displayer) displayFixedChildCount(s *displayState, childCount, depth int) {
    for i := range childCount {
        if int(s.position) == len(s.ast) {
            d.showMissingChildNode(s, depth + 1, childCount - i)

            return
        }

        nextNode := s.ast[s.position]
        s.position++

        if nextNode.Type == EndNode {
            d.showEndNodeInFixedCountNode(s, depth + 1, nextNode)
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

func (d *Displayer) getIndentation(depth int) string {
    return strings.Repeat(" ", int(d.SpacesPerLevel) * depth)
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

func (d *Displayer) showIncorrectEndNode(s *displayState, depth int, node Node) {
    fmt.Fprintf(&s.sb, "%sIncorrect EndNode %s\n", d.getIndentation(depth), d.getEndNodeName(node))
}

func (d *Displayer) showMissingEndNode(s *displayState, depth int) {
    fmt.Fprintf(&s.sb, "%sMissing EndNode\n", d.getIndentation(depth))
}

func (d *Displayer) showMissingChildNode(s *displayState, depth, number int) {
    fmt.Fprintf(&s.sb, "%s%d missing\n", d.getIndentation(depth), number)
}

func (d *Displayer) showEndNodeInFixedCountNode(s *displayState, depth int, node Node) {
    fmt.Fprintf(
        &s.sb,
        "%sEndNode %s in fixed childcount Node\n",
        d.getIndentation(depth),
        d.getEndNodeName(node),
    )
}
