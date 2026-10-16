package decimalparser

import (
	"fmt"
	"math/big"
	"minimal/minimal-lang/built-in/ast"
	"minimal/minimal-lang/built-in/lexer"
	"minimal/minimal-lang/built-in/parsers/prattparser"
)

type number struct {
	big.Int
}

type WholeParser struct {
	tokenType lexer.TokenType
	nodeType  ast.NodeType
	numbers   []number
}

type WholeExponentParser struct {
    tokenType lexer.TokenType
	nodeType  ast.NodeType
	numbers   []wholeExponent
}

type WholeNegativeExponentParser struct {
	tokenType lexer.TokenType
	nodeType  ast.NodeType
    numbers   []wholeExponent
}

type FractionParser struct {
    tokenType lexer.TokenType
	nodeType  ast.NodeType
	numbers   []fraction
}

type FractionExponentParser struct {
	tokenType lexer.TokenType
	nodeType  ast.NodeType
    numbers   []fractionExponent
}

type FractionNegativeExponentParser struct {
	tokenType lexer.TokenType
	nodeType  ast.NodeType
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

func NewWholeParser(tokenType lexer.TokenType, nodeType ast.NodeType) *WholeParser {
	// These are initialized lazily
    return &WholeParser{tokenType, nodeType, nil}
}

func NewWholeExponentParser(tokenType lexer.TokenType, nodeType ast.NodeType) *WholeExponentParser {
	return &WholeExponentParser{tokenType, nodeType, nil}
}

func NewWholeNegativeExponentParser(
	tokenType lexer.TokenType, 
	nodeType ast.NodeType,
) *WholeNegativeExponentParser {
	return &WholeNegativeExponentParser{tokenType, nodeType, nil}
}

func NewFractionParser(tokenType lexer.TokenType, nodeType ast.NodeType) *FractionParser {
	return &FractionParser{tokenType, nodeType, nil}
}

func NewFractionExponentParser(
	tokenType lexer.TokenType, 
	nodeType ast.NodeType,
) *FractionExponentParser {
	return &FractionExponentParser{tokenType, nodeType, nil}
}

func NewFractionNegativeExponentParser(
	tokenType lexer.TokenType, 
	nodeType ast.NodeType,
) *FractionNegativeExponentParser {
	return &FractionNegativeExponentParser{tokenType, nodeType, nil}
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

	w.numbers = append(w.numbers, number)

	return []ast.Node{{Type: w.nodeType, Reference: uint32(len(w.numbers)) - 1}}
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

	bigBase := big.NewInt(int64(base))
	bigDigit := big.NewInt(int64(digit))

	n.Mul(&n.Int, bigBase)
	n.Add(&n.Int, bigDigit)
}

func (n number) String() string {
	return n.Text(10)
}
