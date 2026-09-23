package lexer

import (
	"fmt"
	"minimal/minimal-lang/built-in/ansi"
	"minimal/minimal-lang/built-in/diff"
	"minimal/minimal-lang/built-in/substring"
	"strconv"
	"strings"
)

type Displayer struct {
    scheme      *LexerScheme
    OutputANSI  bool
    tokenColors map[TokenType]ansi.RGB
}

func NewDisplayer(scheme *LexerScheme) *Displayer {
    return &Displayer{scheme, true, map[TokenType]ansi.RGB{}}
}

func (d *Displayer) SetTokenTypeColor(tokenType TokenType, color ansi.RGB) {
    d.tokenColors[tokenType] = color
}

func (d *Displayer) Display(source string, tokens []Token) string {
    sb := &strings.Builder{}

    for _, token := range tokens {
        fmt.Fprintf(sb, "%s\n", d.StringifyToken(source, token))
    }

    return sb.String()
}

func compareTokens(a, b Token) bool {
    return a.Type == b.Type && a.Value == b.Value
}

// Tokens are considered the same if there types and values are equal.
// The range is deliberately ignored since a small change can change all the following ranges.
func (d *Displayer) DisplayDiff(source string, before, after []Token) string {
    return d.DisplayMultiSourceDiff(source, source, before, after)
}

func (d *Displayer) DisplayMultiSourceDiff(
    sourceBefore string, sourceAfter string,
    before, after []Token,
) string {
    tokenDiff := diff.GetDiff(before, after, compareTokens)
    sb := &strings.Builder{}

    for _, diffPart := range tokenDiff {
        source := sourceAfter
        prefix := "  "

        switch diffPart.Type {
        case diff.Insert:
            prefix = "+ "

            if d.OutputANSI {
                prefix = ansi.RGB{R: 121, G: 245, B: 5}.ToString() + prefix + ansi.Reset
            }
        case diff.Delete:
            prefix = "- "
            source = sourceBefore

            if d.OutputANSI {
                prefix = ansi.RGB{R: 245, G: 5, B: 61}.ToString() + prefix + ansi.Reset
            }
        }

        fmt.Fprintf(sb, "%s%s\n", prefix, d.StringifyToken(source, diffPart.Value))
    }

    return sb.String()
}

func (d *Displayer) StringifyToken(source string, token Token) string {
    name := d.scheme.GetTokenTypeMetadata(token.Type).DebugName
    paddedName := fmt.Sprintf("%-20s", name)

    if color, ok := d.tokenColors[token.Type]; ok && d.OutputANSI {
        paddedName = color.ToString() + paddedName + ansi.Reset
    }

    if !substring.IsSubString(source, token.Value) {
        return fmt.Sprintf(
            "%s %-20s     not from source",
            paddedName,
            strconv.Quote(token.Value),
        )
    }

    start := uint(substring.GetStringPtr(token.Value) - substring.GetStringPtr(source))
    length := uint(len(token.Value))

    return fmt.Sprintf(
        "%s %-20s %6d..%-6d (%d)",
        paddedName,
        strconv.Quote(token.Value),
        start,
        start+length,
        length,
    )
}
