package numbers

import (
	"minimal/minimal-lang/built-in/lexer"
	"testing"
)

func getLexer() (*lexer.LexerScheme, lexer.TokenType) {
	l := lexer.NewScheme()
	numberType := l.NewTokenType(
		lexer.TokenTypeMetadata{NounPhrase: "a number", DebugName: "Number"},
	)

	identifierMatcher := NewNumberMatcher(numberType)
	l.AddMatcher(identifierMatcher)

	return l, numberType
}

func TestNumber(t *testing.T) {
	source := "0123456789"

	l, numberType := getLexer()

	expected := []lexer.Token{
		{Type: numberType, Value: source},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMixed(t *testing.T) {
	source := "/012345 :6789"

	l, numberType := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source[:1]},
		{Type: numberType, Value: source[1:7]},
		{Type: lexer.UNKNOWN, Value: source[7:8]},
		{Type: lexer.UNKNOWN, Value: source[8:9]},
		{Type: numberType, Value: source[9:13]},
	}

	// / and : are next to 0 and 9 respectively in the ASCII table
	lexer.CheckTokens(t, l, expected, source)
}

func TestNonDigits(t *testing.T) {
	source := "0_&a 9!"

	l, numberType := getLexer()

	expected := []lexer.Token{
		{Type: numberType, Value: source[:4]},
		{Type: lexer.UNKNOWN, Value: source[4:5]},
		{Type: numberType, Value: source[5:]},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestEndOfLine(t *testing.T) {
	source := "0\n0\r0"

	l, numberType := getLexer()

	expected := []lexer.Token{
		{Type: numberType, Value: source[:1]},
		{Type: lexer.UNKNOWN, Value: source[1:2]},
		{Type: numberType, Value: source[2:3]},
		{Type: lexer.UNKNOWN, Value: source[3:4]},
		{Type: numberType, Value: source[4:5]},
	}

	lexer.CheckTokens(t, l, expected, source)
}
