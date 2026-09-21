package ast

import (
	"testing"
)

type testAST struct {
    d                      *Displayer
    schema                 *ASTSchema
    zeroChildren           NodeType
    oneChild               NodeType
    twoChildren            NodeType

    firstVariableChildren  NodeType
    secondVariableChildren NodeType
}

func getTestAST() testAST {
    schema := NewSchema()

    zeroChildren := schema.NewNodeType(&StructNodeTypeMetadata{DebugName: "Zero"})
    oneChild := schema.NewNodeType(&StructNodeTypeMetadata{DebugName: "One", ChildCount: 1})
    twoChildren := schema.NewNodeType(&StructNodeTypeMetadata{DebugName: "Two", ChildCount: 2})

    firstVariableChildren := schema.NewNodeType(
        &StructNodeTypeMetadata{DebugName: "Variable1", ChildCount: VariableChildCount},
    )

    secondVariableChildren := schema.NewNodeType(
        &StructNodeTypeMetadata{DebugName: "Variable2", ChildCount: VariableChildCount},
    )

    a := NewDisplayer(schema)

    return testAST{
        a,
        schema,
        zeroChildren,
        oneChild,
        twoChildren,
        firstVariableChildren,
        secondVariableChildren,
    }
}

func TestCorrect(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.zeroChildren, 0},
        {ta.oneChild, 0},
        {ta.zeroChildren, 0},
        {ta.twoChildren, 0},
        {ta.oneChild, 0},
        {ta.zeroChildren, 0},
        {ta.firstVariableChildren, 0},
        {ta.zeroChildren, 0},
        {ta.secondVariableChildren, 0},
        {ta.zeroChildren, 0},
        {EndNode, uint32(ta.secondVariableChildren)},
        {ta.zeroChildren, 0},
        {EndNode, uint32(ta.firstVariableChildren)},
    }

    expected := "Zero\n" +
                "One\n" +
                "  Zero\n" +
                "Two\n" +
                "  One\n" +
                "    Zero\n" +
                "  Variable1\n" +
                "    Zero\n" +
                "    Variable2\n" +
                "      Zero\n" +
                "    Zero\n"

    actual := ta.d.Display(ast)

    if actual != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, actual)
    }
}

func TestIncorrectFixedChildren(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.twoChildren, 0},
        {ta.oneChild, 0},
    }

    expected := "Two\n" +
                "  One\n" +
                "    1 missing\n" +
                "  1 missing\n"

    actual := ta.d.Display(ast)

    if actual != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, actual)
    }
}

func TestMissingEndNode(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.firstVariableChildren, 0},
        {ta.zeroChildren, 0},
    }

    expected := "Variable1\n" +
                "  Zero\n" +
                "Missing EndNode\n"

    actual := ta.d.Display(ast)

    if actual != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, actual)
    }
}

func TestMissingEndNodeNested(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.firstVariableChildren, 0},
        {ta.zeroChildren, 0},
        {ta.secondVariableChildren, 0},
        {ta.zeroChildren, 0},
        {ta.zeroChildren, 0},
        {EndNode, uint32(ta.firstVariableChildren)},
    }

    expected := "Variable1\n" +
                "  Zero\n" +
                "  Variable2\n" +
                "    Zero\n" +
                "    Zero\n" +
                "  Incorrect EndNode Variable1\n"

    actual := ta.d.Display(ast)

    if actual != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, actual)
    }
}

func TestEndNodeInFixedChildrenNode(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.oneChild, 0},
        {EndNode, uint32(ta.firstVariableChildren)},
    }

    expected := "One\n" +
                "  EndNode Variable1 in fixed childcount Node\n"

    actual := ta.d.Display(ast)

    if actual != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, actual)
    }
}

func TestUnknownEndNodeReference(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.firstVariableChildren, 0},
        {EndNode, uint32(100)},
    }

    expected := "Variable1\n" +
                "Incorrect EndNode UNKNOWN Reference=100\n" +
                "EndNode UNKNOWN Reference=100 not inside a Node\n"

    actual := ta.d.Display(ast)

    if actual != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, actual)
    }
}

func TestCustomDisplay(t *testing.T) {

}

func TestDiff(t *testing.T) {
    ta := getTestAST()

    before := []Node{
        {ta.zeroChildren, 0},
        {ta.zeroChildren, 0},
        {ta.oneChild, 0},
            {ta.oneChild, 0},
                {ta.zeroChildren, 0},
        {ta.twoChildren, 0},
            {ta.zeroChildren, 0},
            {ta.zeroChildren, 0},
        {ta.twoChildren, 0},
            {ta.zeroChildren, 0},
            {ta.zeroChildren, 0},
    }

    after := []Node{
        {ta.zeroChildren, 0},
        {ta.zeroChildren, 1},
        {ta.oneChild, 0},
            {ta.oneChild, 0},
                {ta.zeroChildren, 1},
        {ta.twoChildren, 0},
            {ta.zeroChildren, 0},
            {ta.zeroChildren, 1},
        {ta.twoChildren, 0},
            {ta.zeroChildren, 1},
            {ta.zeroChildren, 0},
    }

    expected := "  Zero\n" +
                "- Zero\n" +
                "+ Zero Reference=1\n" +
                "  One\n" +
                "    One\n" +
                "    - Zero\n" +
                "    + Zero Reference=1\n" +
                "  Two\n" +
                "    Zero\n" +
                "  + Zero Reference=1\n" +
                "  - Zero\n" +
                "  Two\n" +
                "  - Zero\n" +
                "  + Zero Reference=1\n" +
                "    Zero\n"

    ta.d.OutputANSI = false
    actual := ta.d.DisplayDiff(before, after)

    if actual != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, actual)
    }
}

func TestColoredDiff(t *testing.T) {
    ta := getTestAST()

    before := []Node{
        {ta.zeroChildren, 0},
        {ta.zeroChildren, 1},
    }

    after := []Node{
        {ta.zeroChildren, 1},
        {ta.zeroChildren, 0},
    }

    expected := "\x1b[38;2;245;5;61m- \x1b[0mZero\n" +
                "  Zero Reference=1\n" +
                "\x1b[38;2;121;245;5m+ \x1b[0mZero\n"

    actual := ta.d.DisplayDiff(before, after)

    if actual != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, actual)
    }
}
