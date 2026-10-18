package decimalparser

import (
	"math/big"
	"minimal/minimal-lang/built-in/ast"
	"minimal/minimal-lang/built-in/lexer"
	"minimal/minimal-lang/built-in/messenger"
	"minimal/minimal-lang/built-in/parsers/prattparser"
)

type WholeParser struct {
	messenger *messenger.Messenger
	tokenType lexer.TokenType
	nodeType  ast.NodeType
	numbers   []*big.Int
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
	numbers   []*big.Rat
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

type wholeExponent struct {
	whole    *big.Int
	exponent *big.Int
}

type fractionExponent struct {
	fraction *big.Rat
	exponent *big.Int
}

func NewWholeParser(
	tokenType lexer.TokenType, 
	nodeType ast.NodeType, 
	messenger *messenger.Messenger,
) *WholeParser {
    return &WholeParser{messenger, tokenType, nodeType, nil}
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

func (w *WholeParser) ParsePrefix(
	_ *prattparser.PrattParser, 
	l *lexer.Lexer, 
	minBindingPower uint32,
) []ast.Node {
	text := l.Peek(0).Value
	number := big.NewInt(0)

	number.Mul(number, big.NewInt(10))
	number.Add(number, big.NewInt(int64(text[0] - 48)))

	previousWasUnderScore := false
	
	for i, c := range []byte(text[1:]) {
		if c != '_' {
			number.Mul(number, big.NewInt(10))
			number.Add(number, big.NewInt(int64(c - 48)))
		} else if previousWasUnderScore {
			w.messenger.Send(
				messenger.Message{
					Message: "Two successive underscore in number literal", 
					Severity: messenger.Warning,
					Context: []messenger.Span{{Content: text[i:i + 2]}},
				},
			)
		} else {
			previousWasUnderScore = true
		}
	}

	if previousWasUnderScore {
		w.messenger.Send(
			messenger.Message{
				Message: "Number literal ends with underscore", 
				Severity: messenger.Warning,
				Context: []messenger.Span{{Content: text[len(text) - 1:]}},
			},
		)
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
