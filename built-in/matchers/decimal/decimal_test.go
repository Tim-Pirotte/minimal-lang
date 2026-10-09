package decimalmatcher

import (
	"minimal/minimal-lang/built-in/lexer"
	"testing"
)

func getLexer() (
    *lexer.LexerScheme, 
    lexer.TokenType, 
    lexer.TokenType, 
    lexer.TokenType, 
    lexer.TokenType, 
    lexer.TokenType, 
    lexer.TokenType,
) {
	l := lexer.NewScheme()

	intType := l.NewTokenType(
        lexer.TokenTypeMetadata{NounPhrase: "a decimal integer literal", DebugName: "DecimalInt"},
    )

    wholeFloatExpType := l.NewTokenType(
        lexer.TokenTypeMetadata{
            NounPhrase: "a decimal whole float literal with exponent", 
            DebugName: "DecimalWholeFloatExponent",
        },
    )

    wholeFloatNegExpType := l.NewTokenType(
        lexer.TokenTypeMetadata{
            NounPhrase: "a decimal whole float literal with negative exponent", 
            DebugName: "DecimalWholeFloatNegativeExponent",
        },
    )

    floatType := l.NewTokenType(
        lexer.TokenTypeMetadata{NounPhrase: "a decimal float literal", DebugName: "DecimalFloat"},
    )

    floatExpType := l.NewTokenType(
        lexer.TokenTypeMetadata{
            NounPhrase: "a decimal float literal with exponent", 
            DebugName: "DecimalFloatExponent",
        },
    )

    floatNegativeExpType := l.NewTokenType(
        lexer.TokenTypeMetadata{
            NounPhrase: "a decimal float literal with negative exponent", 
            DebugName: "DecimalFloatNegativeExponent",
        },
    )

	l.AddMatcher(
        New(
            intType, 
            wholeFloatExpType, 
            wholeFloatNegExpType, 
            floatType, 
            floatExpType, 
            floatNegativeExpType,
        ),
    )

	return l, intType, wholeFloatExpType, wholeFloatNegExpType, floatType, floatExpType, floatNegativeExpType
}

func TestValidInt(t *testing.T) {
	source := "0123456789"

	l, intType, _, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: intType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestValidFloat(t *testing.T) {
	source := "0123456789.0123456789"

	l, _, _, _, floatType, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestWholeFloatExponent(t *testing.T) {
    source := "0123456789e0123456789"

	l, _, wholeFloatExpType, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestWholeFloatNegativeExponent(t *testing.T) {
    source := "0123456789e-0123456789"

	l, _, _, wholeFloatNegExpType, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatExponent(t *testing.T) {
    source := "0123456789.0123456789e0123456789"

	l, _, _, _, _, floatExpType, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatNegativeExponent(t *testing.T) {
    source := "0123456789.0123456789e-0123456789"

	l, _, _, _, _, _, floatNegExpType := getLexer()

	expected := []lexer.Token{
		{Type: floatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestIntUnderscores(t *testing.T) {
	source := "_0__1_2_3_4_5_6_7_8_9_"

	l, intType, _, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: intType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatUnderscores(t *testing.T) {
	source := "_0__1_2_3_4_5_6_7_8_9.0_1_2_3_4_5_6_7_8_9_"

	l, _, _, _, floatType, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestWholeFloatExponentUnderscores(t *testing.T) {
    source := "_0__123456789__e_0_1_2_3_4__5_6_7_8_9__"

	l, _, wholeFloatExpType, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestWholeFloatNegativeExponentUnderscores(t *testing.T) {
    source := "0123456789__e-__0__123456789__"

	l, _, _, wholeFloatNegExpType, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatExponentUnderscores(t *testing.T) {
    source := "012345678__9__.__0123456789__e__0_123456789__"

	l, _, _, _, _, floatExpType, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatNegativeExponentUnderscores(t *testing.T) {
    source := "0123456789.0123456789__e-__0_123456789__"

	l, _, _, _, _, _, floatNegExpType := getLexer()

	expected := []lexer.Token{
		{Type: floatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingFractionalPart(t *testing.T) {
	source := "0123456789."

	l, _, _, _, floatType, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingWholeFloatExponent(t *testing.T) {
	source := "0123456789e"

	l, _, wholeFloatExpType, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingWholeFloatNegativeExponent(t *testing.T) {
	source := "0123456789e-"

	l, _, _, wholeFloatNegExpType, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingFloatExponent(t *testing.T) {
	source := "0123456789.0123456789e"

	l, _, _, _, _, floatExpType, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingFloatNegativeExponent(t *testing.T) {
	source := "0123456789.0123456789e-"

	l, _, _, _, _, _, floatNegExpType := getLexer()

	expected := []lexer.Token{
		{Type: floatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}
