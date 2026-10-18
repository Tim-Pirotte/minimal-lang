package octalmatcher

import (
	"minimal/minimal-lang/built-in/lexer"
	"testing"
)

func getLexer() (*lexer.LexerScheme, lexer.TokenType, lexer.TokenType) {
	l := lexer.NewScheme()

	intType := l.NewTokenType(
		lexer.TokenTypeMetadata{NounPhrase: "an octal integer literal", DebugName: "OctalInt"},
	)

	floatType := l.NewTokenType(
		lexer.TokenTypeMetadata{NounPhrase: "an octal float literal", DebugName: "OctalFloat"},
	)

	l.AddMatcher(New(intType, floatType))

	return l, intType, floatType
}

func TestValidInt(t *testing.T) {
	source := "0o01234567"

	l, intType, _ := getLexer()
	expected := []lexer.Token{{Type: intType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestValidFloat(t *testing.T) {
	source := "0o01234567.01234567"

	l, _, floatType := getLexer()
	expected := []lexer.Token{{Type: floatType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestIntUnderscores(t *testing.T) {
	source := "0o0_1_2_3456___7__"

	l, intType, _ := getLexer()
	expected := []lexer.Token{{Type: intType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatUnderscores(t *testing.T) {
	source := "0o__123__4__.__5_67__"

	l, _, floatType := getLexer()
	expected := []lexer.Token{{Type: floatType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingDigits(t *testing.T) {
    source := "0o"

	l, intType, _ := getLexer()
	expected := []lexer.Token{{Type: intType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingFractionalPart(t *testing.T) {
	source := "0o1234567."

	l, _, floatType := getLexer()
	expected := []lexer.Token{{Type: floatType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestDigitBounds(t *testing.T) {
    source := "0o8" + "0o/" + "0o.8" + "0o./"

	l, intType, floatType := getLexer()
    
	expected := []lexer.Token{
        {Type: intType, Value: source[:2]},
        {Type: lexer.UNKNOWN, Value: source[2:3]},
        {Type: intType, Value: source[3:5]},
        {Type: lexer.UNKNOWN, Value: source[5:6]},
        {Type: floatType, Value: source[6:9]},
        {Type: lexer.UNKNOWN, Value: source[9:10]},
        {Type: floatType, Value: source[10:13]},
        {Type: lexer.UNKNOWN, Value: source[13:14]},
    }

	lexer.CheckTokens(t, l, expected, source)
}
