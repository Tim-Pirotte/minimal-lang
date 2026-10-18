package hexmatcher

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
		lexer.TokenTypeMetadata{NounPhrase: "a hexadecimal integer literal", DebugName: "HexInt"},
	)

	wholeFloatExpType := l.NewTokenType(
		lexer.TokenTypeMetadata{
			NounPhrase: "a hexadecimal whole float literal with exponent",
			DebugName:  "HexWholeFloatExponent",
		},
	)

	wholeFloatNegExpType := l.NewTokenType(
		lexer.TokenTypeMetadata{
			NounPhrase: "a hexadecimal whole float literal with negative exponent",
			DebugName:  "HexWholeFloatNegativeExponent",
		},
	)

	floatType := l.NewTokenType(
		lexer.TokenTypeMetadata{NounPhrase: "a hexadecimal float literal", DebugName: "HexFloat"},
	)

	floatExpType := l.NewTokenType(
		lexer.TokenTypeMetadata{
			NounPhrase: "a hexadecimal float literal with exponent",
			DebugName:  "HexFloatExponent",
		},
	)

	floatNegativeExpType := l.NewTokenType(
		lexer.TokenTypeMetadata{
			NounPhrase: "a hexadecimal float literal with negative exponent",
			DebugName:  "HexFloatNegativeExponent",
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
	source := "0x0123456789ABCDEF"

	l, intType, _, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: intType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestValidFloat(t *testing.T) {
	source := "0x0123456789ABCDEF.0123456789ABCDEF"

	l, _, _, _, floatType, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestWholeFloatExponent(t *testing.T) {
	source := "0x0123456789ABCDEFp0123456789"

	l, _, wholeFloatExpType, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestWholeFloatNegativeExponent(t *testing.T) {
	source := "0x0123456789ABCDEFp-0123456789"

	l, _, _, wholeFloatNegExpType, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatExponent(t *testing.T) {
	source := "0x0123456789ABCDEF.0123456789ABCDEFp0123456789"

	l, _, _, _, _, floatExpType, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatNegativeExponent(t *testing.T) {
	source := "0x0123456789ABCDEF.0123456789ABCDEFp-0123456789"

	l, _, _, _, _, _, floatNegExpType := getLexer()

	expected := []lexer.Token{
		{Type: floatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestIntUnderscores(t *testing.T) {
	source := "0x__0_1_2_3_4_5_6_7_8_9__ABCDEF__"

	l, intType, _, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: intType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatUnderscores(t *testing.T) {
	source := "0x_0_1_2_3_4_5_6_7_8_9_ABCDEF__.0_1_2_3_4_5_6_7_8_9_ABCDEF__"

	l, _, _, _, floatType, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestWholeFloatExponentUnderscores(t *testing.T) {
	source := "0x_0_123456789ABCDEF__p_0_1_2_3_4__5_6_7_8_9__"

	l, _, wholeFloatExpType, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestWholeFloatNegativeExponentUnderscores(t *testing.T) {
	source := "0x0123456789ABCDEF__p-__0__123456789__"

	l, _, _, wholeFloatNegExpType, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatExponentUnderscores(t *testing.T) {
	source := "0x012345678__9_ABCDEF_.__0123456789_ABCDEF_p__0_123456789__"

	l, _, _, _, _, floatExpType, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFloatNegativeExponentUnderscores(t *testing.T) {
	source := "0x0123456789ABCDEF.0123456789_ABCDEF_p-__0_123456789__"

	l, _, _, _, _, _, floatNegExpType := getLexer()

	expected := []lexer.Token{
		{Type: floatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingFractionalPart(t *testing.T) {
	source := "0x0123456789ABCDEF."

	l, _, _, _, floatType, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingWholeFloatExponent(t *testing.T) {
	source := "0x0123456789ABCDEFp"

	l, _, wholeFloatExpType, _, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingWholeFloatNegativeExponent(t *testing.T) {
	source := "0x0123456789p-"

	l, _, _, wholeFloatNegExpType, _, _, _ := getLexer()

	expected := []lexer.Token{
		{Type: wholeFloatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingFloatExponent(t *testing.T) {
	source := "0x0123456789ABCDEF.0123456789ABCDEFp"

	l, _, _, _, _, floatExpType, _ := getLexer()

	expected := []lexer.Token{
		{Type: floatExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMissingFloatNegativeExponent(t *testing.T) {
	source := "0x0123456789ABCDEF.0123456789ABCDEFp-"

	l, _, _, _, _, _, floatNegExpType := getLexer()

	expected := []lexer.Token{
		{Type: floatNegExpType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestDigitBounds(t *testing.T) {
	source := "0x/" + "0x:" + "0x@" + "0xG" + "0xp/" + "0xp:" + "0xp@" + "0xpG" + 
        "0xp-/" + "0xp-:" + "0xp-@" + "0xp-G" + "0x./" + "0x.:" + "0x.@" + "0x.G" + 
        "0x.p/" + "0x.p:" + "0x.p@" + "0x.pG" + "0x.p-/" + "0x.p-:" + "0x.p-@" + "0x.p-G"

	l, intType, intExpType, intNegExpType, floatType, floatExpType, floatNegExpType := getLexer()

	expected := []lexer.Token{
		{Type: intType, Value: source[:2]},
		{Type: lexer.UNKNOWN, Value: source[2:3]},
        {Type: intType, Value: source[3:5]},
		{Type: lexer.UNKNOWN, Value: source[5:6]},
        {Type: intType, Value: source[6:8]},
		{Type: lexer.UNKNOWN, Value: source[8:9]},
        {Type: intType, Value: source[9:11]},
		{Type: lexer.UNKNOWN, Value: source[11:12]},

        {Type: intExpType, Value: source[12:15]},
		{Type: lexer.UNKNOWN, Value: source[15:16]},
        {Type: intExpType, Value: source[16:19]},
		{Type: lexer.UNKNOWN, Value: source[19:20]},
        {Type: intExpType, Value: source[20:23]},
		{Type: lexer.UNKNOWN, Value: source[23:24]},
        {Type: intExpType, Value: source[24:27]},
		{Type: lexer.UNKNOWN, Value: source[27:28]},

        {Type: intNegExpType, Value: source[28:32]},
		{Type: lexer.UNKNOWN, Value: source[32:33]},
        {Type: intNegExpType, Value: source[33:37]},
		{Type: lexer.UNKNOWN, Value: source[37:38]},
        {Type: intNegExpType, Value: source[38:42]},
		{Type: lexer.UNKNOWN, Value: source[42:43]},
        {Type: intNegExpType, Value: source[43:47]},
		{Type: lexer.UNKNOWN, Value: source[47:48]},

        {Type: floatType, Value: source[48:51]},
		{Type: lexer.UNKNOWN, Value: source[51:52]},
        {Type: floatType, Value: source[52:55]},
		{Type: lexer.UNKNOWN, Value: source[55:56]},
        {Type: floatType, Value: source[56:59]},
		{Type: lexer.UNKNOWN, Value: source[59:60]},
        {Type: floatType, Value: source[60:63]},
		{Type: lexer.UNKNOWN, Value: source[63:64]},

        {Type: floatExpType, Value: source[64:68]},
		{Type: lexer.UNKNOWN, Value: source[68:69]},
        {Type: floatExpType, Value: source[69:73]},
		{Type: lexer.UNKNOWN, Value: source[73:74]},
        {Type: floatExpType, Value: source[74:78]},
		{Type: lexer.UNKNOWN, Value: source[78:79]},
        {Type: floatExpType, Value: source[79:83]},
		{Type: lexer.UNKNOWN, Value: source[83:84]},

        {Type: floatNegExpType, Value: source[84:89]},
		{Type: lexer.UNKNOWN, Value: source[89:90]},
        {Type: floatNegExpType, Value: source[90:95]},
		{Type: lexer.UNKNOWN, Value: source[95:96]},
        {Type: floatNegExpType, Value: source[96:101]},
		{Type: lexer.UNKNOWN, Value: source[101:102]},
        {Type: floatNegExpType, Value: source[102:107]},
		{Type: lexer.UNKNOWN, Value: source[107:108]},
	}

	lexer.CheckTokens(t, l, expected, source)
}
