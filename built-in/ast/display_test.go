package ast

import (
	"bytes"
	"errors"
	"minimal/minimal-lang/built-in/messenger"
	"minimal/minimal-lang/built-in/outputs/test-output"
	"testing"
)

type testAST struct {
    a                      *Displayer
    schema                 *ASTSchema
    buf                    *bytes.Buffer
    messenger              *messenger.Messenger
    to                     *testoutput.TestOutput
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

    m := messenger.New()
    to := testoutput.New()
    m.AddOutput(to)

    var buf bytes.Buffer

    a := NewDisplayer(m, schema, &buf)

    return testAST{
        a,
        schema,
        &buf,
        m,
        to,
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

    ta.a.Display(ast)

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

    if ta.buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, ta.buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

func TestIncorrectFixedChildren(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.twoChildren, 0},
        {ta.oneChild, 0},
    }

    ta.a.Display(ast)

    expected := "Two\n" +
                "  One\n" +
                "    1 missing\n" +
                "  1 missing\n"

    if ta.buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, ta.buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

func TestMissingEndNode(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.firstVariableChildren, 0},
        {ta.zeroChildren, 0},
    }

    ta.a.Display(ast)

    expected := "Variable1\n" +
                "  Zero\n" +
                "Missing EndNode\n"

    if ta.buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, ta.buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
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

    ta.a.Display(ast)

    expected := "Variable1\n" +
                "  Zero\n" +
                "  Variable2\n" +
                "    Zero\n" +
                "    Zero\n" +
                "  Incorrect EndNode Variable1\n"

    if ta.buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, ta.buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

func TestEndNodeInFixedChildrenNode(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.oneChild, 0},
        {EndNode, uint32(ta.firstVariableChildren)},
    }

    ta.a.Display(ast)

    expected := "One\n" +
                "  EndNode Variable1 in fixed childcount Node\n"

    if ta.buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, ta.buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

func TestUnknownEndNodeReference(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.firstVariableChildren, 0},
        {EndNode, uint32(100)},
    }

    ta.a.Display(ast)

    expected := "Variable1\n" +
                "Incorrect EndNode UNKNOWN Reference=100\n" +
                "EndNode UNKNOWN Reference=100 not inside a Node\n"

    if ta.buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, ta.buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) {
    return 0, errors.New("")
}

func TestFailingWriter(t *testing.T) {
    ta := getTestAST()
    writer := failingWriter{}

    ast := []Node{{ta.zeroChildren, 0}}

    d := NewDisplayer(ta.messenger, ta.schema, writer)
    d.Display(ast)

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{{
        Message: "AST debugger output write failed",
        Severity: messenger.Error,
    }})
}
