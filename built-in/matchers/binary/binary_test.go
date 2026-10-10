package binarymatcher

import (
	"minimal/minimal-lang/built-in/lexer"
	"testing"
)

func getLexer() (*lexer.LexerScheme, lexer.TokenType, lexer.TokenType) {
	l := lexer.NewScheme()

	intType := l.NewTokenType(
		lexer.TokenTypeMetadata{NounPhrase: "a binary integer literal", DebugName: "BinaryInt"},
	)

	floatType := l.NewTokenType(
		lexer.TokenTypeMetadata{NounPhrase: "a binary float literal", DebugName: "BinaryFloat"},
	)

	l.AddMatcher(New(intType, floatType))

	return l, intType, floatType
}

func TestValidInt(t *testing.T) {
	source := "0b10001000101110"

	l, intType, _ := getLexer()
	expected := []lexer.Token{{Type: intType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestValidFloat(t *testing.T) {
	source := "0b101111110.011110111"

	l, _, floatType := getLexer()
	expected := []lexer.Token{{Type: floatType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestIntUnderscores(t *testing.T) {
	source := "0b0_1_1_00010___1__"

	l, intType, _ := getLexer()
	expected := []lexer.Token{{Type: intType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatUnderscores(t *testing.T) {
	source := "0b__1010__1__.__1_100__"

	l, _, floatType := getLexer()
	expected := []lexer.Token{{Type: floatType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingDigits(t *testing.T) {
    source := "0b"

	l, intType, _ := getLexer()
	expected := []lexer.Token{{Type: intType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingFractionalPart(t *testing.T) {
	source := "0b101101."

	l, _, floatType := getLexer()
	expected := []lexer.Token{{Type: floatType, Value: source}}

	lexer.CheckTokens(t, l, expected, source)
}

func TestDigitBounds(t *testing.T) {
    source := "0b2" + "0b/" + "0b.2" + "0b./"

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
