package decimalparser

import (
	"fmt"
	"math/bits"
	"minimal/minimal-lang/built-in/ast"
	"minimal/minimal-lang/built-in/lexer"
	"minimal/minimal-lang/built-in/parsers/prattparser"
)

type number []uint

type WholeParser struct {
	tokenType lexer.TokenType
	numbers   []number
}

type WholeExponentParser struct {
    tokenType lexer.TokenType
	numbers   []wholeExponent
}

type WholeNegativeExponentParser struct {
	tokenType lexer.TokenType
    numbers   []wholeExponent
}

type FractionParser struct {
    tokenType lexer.TokenType
	numbers   []fraction
}

type FractionExponentParser struct {
	tokenType lexer.TokenType
    numbers   []fractionExponent
}

type FractionNegativeExponentParser struct {
	tokenType lexer.TokenType
    numbers   []fractionExponent
}

type fraction struct {
	decimal    number
	fractional number
}

type wholeExponent struct {
	whole    number
	exponent number
}

type fractionExponent struct {
	fraction fraction
	exponent number
}

func NewWholeParser(tokenType lexer.TokenType) *WholeParser {
	// These are initialized lazily
    return &WholeParser{tokenType, nil}
}

func NewWholeExponentParser(tokenType lexer.TokenType) *WholeExponentParser {
	return &WholeExponentParser{tokenType, nil}
}

func NewWholeNegativeExponentParser(tokenType lexer.TokenType) *WholeNegativeExponentParser {
	return &WholeNegativeExponentParser{tokenType, nil}
}

func NewFractionParser(tokenType lexer.TokenType) *FractionParser {
	return &FractionParser{tokenType, nil}
}

func NewFractionExponentParser(tokenType lexer.TokenType) *FractionExponentParser {
	return &FractionExponentParser{tokenType, nil}
}

func NewFractionNegativeExponentParser(tokenType lexer.TokenType) *FractionNegativeExponentParser {
	return &FractionNegativeExponentParser{tokenType, nil}
}

func (w *WholeParser) GetTokenType() lexer.TokenType {
    return w.tokenType
}

func (w *WholeParser) ParsePrefix(_ *prattparser.PrattParser, l *lexer.Lexer, minBindingPower uint32) []ast.Node {
	text := l.Peek(0).Value
	number := number{}

	number.addDigit(text[0] - 48, 10)

	previousWasUnderScore := false
	
	for _, c := range []byte(text[1:]) {
		if c != '_' {
			number.addDigit(c - 48, 10)
		} else if previousWasUnderScore {
			// TODO error
		} else {
			previousWasUnderScore = true
		}
	}

	if previousWasUnderScore {
		// TODO error
	}

	return []ast.Node{}
}

func (w *WholeExponentParser) GetTokenType() lexer.TokenType {
    return w.tokenType
}

func (w *WholeExponentParser) ParsePrefix(_ *prattparser.PrattParser, l *lexer.Lexer, minBindingPower uint32) []ast.Node {
	return []ast.Node{}
}

func (w *WholeNegativeExponentParser) GetTokenType() lexer.TokenType {
    return w.tokenType
}

func (w *WholeNegativeExponentParser) ParsePrefix(_ *prattparser.PrattParser, l *lexer.Lexer, minBindingPower uint32) []ast.Node {
	return []ast.Node{}
}

func (f *FractionParser) GetTokenType() lexer.TokenType {
    return f.tokenType
}

func (f *FractionParser) ParsePrefix(_ *prattparser.PrattParser, l *lexer.Lexer, minBindingPower uint32) []ast.Node {
	return []ast.Node{}
}

func (f *FractionExponentParser) GetTokenType() lexer.TokenType {
    return f.tokenType
}

func (f *FractionExponentParser) ParsePrefix(_ *prattparser.PrattParser, l *lexer.Lexer, minBindingPower uint32) []ast.Node {
	return []ast.Node{}
}

func (f *FractionNegativeExponentParser) GetTokenType() lexer.TokenType {
    return f.tokenType
}

func (f *FractionNegativeExponentParser) ParsePrefix(_ *prattparser.PrattParser, l *lexer.Lexer, minBindingPower uint32) []ast.Node {
	return []ast.Node{}
}

func (n *number) addDigit(digit, base uint8) {
	if digit >= base {
        panic(fmt.Sprintf("digit %d is not valid for base %d", digit, base))
    }

	carry := digit

	for j := range *n {
        // I think that hi can only reach MAX - 1 so we can ignore the carry below
        hi, lo := bits.Mul((*n)[j], uint(base))
        lo, carryOut := bits.Add(lo, uint(carry), 0)
        hi, _ = bits.Add(hi, 0, carryOut)

        (*n)[j] = lo
        carry = uint8(hi)
    }

    if carry > 0 {
		*n = append(*n, uint(carry))
	}
}
