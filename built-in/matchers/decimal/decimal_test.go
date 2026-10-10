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
	source := "0__1_2_3_4_5_6_7_8_9_"

	l, intType, _, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: intType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatUnderscores(t *testing.T) {
	source := "0__1_2_3_4_5_6_7_8_9.0_1_2_3_4_5_6_7_8_9_"

	l, _, _, _, floatType, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestWholeFloatExponentUnderscores(t *testing.T) {
    source := "0__123456789__e_0_1_2_3_4__5_6_7_8_9__"

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

func TestDigitBounds(t *testing.T) {
    source := "/" + ":" + "0e/" + "0e:" + "0e-/" + "0e-:" + "0./" + 
			  "0.:" + "0.e/" + "0.e:" + "0.e-/" + "0.e-:"

	l, _, intExpType, intNegExpType, floatType, floatExpType, floatNegExpType := getLexer()

	expected := []lexer.Token{
        {Type: lexer.UNKNOWN, Value: source[:1]},
        {Type: lexer.UNKNOWN, Value: source[1:2]},
		{Type: intExpType, Value: source[2:4]},
        {Type: lexer.UNKNOWN, Value: source[4:5]},
		{Type: intExpType, Value: source[5:7]},
        {Type: lexer.UNKNOWN, Value: source[7:8]},
		{Type: intNegExpType, Value: source[8:11]},
        {Type: lexer.UNKNOWN, Value: source[11:12]},
		{Type: intNegExpType, Value: source[12:15]},
        {Type: lexer.UNKNOWN, Value: source[15:16]},
		{Type: floatType, Value: source[16:18]},
        {Type: lexer.UNKNOWN, Value: source[18:19]},
		{Type: floatType, Value: source[19:21]},
        {Type: lexer.UNKNOWN, Value: source[21:22]},
		{Type: floatExpType, Value: source[22:25]},
        {Type: lexer.UNKNOWN, Value: source[25:26]},
		{Type: floatExpType, Value: source[26:29]},
        {Type: lexer.UNKNOWN, Value: source[29:30]},
		{Type: floatNegExpType, Value: source[30:34]},
        {Type: lexer.UNKNOWN, Value: source[34:35]},
		{Type: floatNegExpType, Value: source[35:39]},
        {Type: lexer.UNKNOWN, Value: source[39:40]},
    }

	lexer.CheckTokens(t, l, expected, source)
}

func TestNoLeadingUnderscore(t *testing.T) {
	source := "_"

	l, _, _, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestNoLeadingDot(t *testing.T) {
	source := "."

	l, _, _, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestNoLeadingE(t *testing.T) {
	source := "e"

	l, _, _, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}
