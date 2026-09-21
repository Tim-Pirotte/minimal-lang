package lexer

import (
	"minimal/minimal-lang/built-in/ansi"
	"testing"
)

func TestDisplay(t *testing.T) {
    source := "abc"

    s := NewScheme()
    l := s.Lex(source)

    d := NewDisplayer(s)

    tokens := []Token{}

    for l.Peek(0).Type != END {
        tokens = append(tokens, l.Peek(0))
        l.Advance()
    }

    expected := "UNKNOWN              \"a\"                       0..1      (1)\n" +
                "UNKNOWN              \"b\"                       1..2      (1)\n" +
                "UNKNOWN              \"c\"                       2..3      (1)\n"

    actual := d.Display(source, tokens)

    if actual != expected {
        t.Errorf("\nExpected:\n%s\nGot:\n%s", expected, actual)
    }
}

func TestColor(t *testing.T) {
    source := "a"

    s := NewScheme()
    l := s.Lex(source)

    d := NewDisplayer(s)
    d.SetTokenTypeColor(UNKNOWN, ansi.RGB{R: 197, G: 255, B: 23})

    tokens := []Token{}

    for l.Peek(0).Type != END {
        tokens = append(tokens, l.Peek(0))
        l.Advance()
    }

    expected := "\x1b[38;2;197;255;23mUNKNOWN             \x1b[0m \"a\"                       0..1      (1)\n"

    actual := d.Display(source, tokens)

    if actual != expected {
        t.Errorf("\nExpected:\n%s\nGot:\n%s", expected, actual)
    }
}

func TestDiff(t *testing.T) {
    source := "ab"

    s := NewScheme()

    l1 := s.Lex(source[:1])
    l2 := s.Lex(source[1:])

    d := NewDisplayer(s)
    d.OutputANSI = false

    tokens1 := []Token{}

    for l1.Peek(0).Type != END {
        tokens1 = append(tokens1, l1.Peek(0))
        l1.Advance()
    }

    tokens2 := []Token{}

    for l2.Peek(0).Type != END {
        tokens2 = append(tokens2, l2.Peek(0))
        l2.Advance()
    }

    expected := " - UNKNOWN              \"a\"                       0..1      (1)\n" +
                " + UNKNOWN              \"b\"                       1..2      (1)\n"

    actual := d.DisplayDiff(source, tokens1, tokens2)

    if actual != expected {
        t.Errorf("\nExpected:\n%s\nGot:\n%s", expected, actual)
    }
}

func TestColoredDiff(t *testing.T) {
    source := "ab"

    s := NewScheme()

    l1 := s.Lex(source[:1])
    l2 := s.Lex(source[1:])

    d := NewDisplayer(s)

    tokens1 := []Token{}

    for l1.Peek(0).Type != END {
        tokens1 = append(tokens1, l1.Peek(0))
        l1.Advance()
    }

    tokens2 := []Token{}

    for l2.Peek(0).Type != END {
        tokens2 = append(tokens2, l2.Peek(0))
        l2.Advance()
    }

    expected := "\x1b[38;2;245;5;61m - \x1b[0mUNKNOWN              \"a\"                       0..1      (1)\n" +
                "\x1b[38;2;121;245;5m + \x1b[0mUNKNOWN              \"b\"                       1..2      (1)\n"

    actual := d.DisplayDiff(source, tokens1, tokens2)

    if actual != expected {
        t.Errorf("\nExpected:\n%s\nGot:\n%s", expected, actual)
    }
}

func TestMultiSourceDiff(t *testing.T) {
    source1 := "aa"
    source2 := "ba"

    s := NewScheme()

    l1 := s.Lex(source1)
    l2 := s.Lex(source2)

    d := NewDisplayer(s)
    d.OutputANSI = false

    tokens1 := []Token{}

    for l1.Peek(0).Type != END {
        tokens1 = append(tokens1, l1.Peek(0))
        l1.Advance()
    }

    tokens2 := []Token{}

    for l2.Peek(0).Type != END {
        tokens2 = append(tokens2, l2.Peek(0))
        l2.Advance()
    }

    expected := " - UNKNOWN              \"a\"                       0..1      (1)\n" +
                " + UNKNOWN              \"b\"                       0..1      (1)\n" +
                "   UNKNOWN              \"a\"                       1..2      (1)\n"

    actual := d.DisplayMultiSourceDiff(source1, source2, tokens1, tokens2)

    if actual != expected {
        t.Errorf("\nExpected:\n%s\nGot:\n%s", expected, actual)
    }
}
