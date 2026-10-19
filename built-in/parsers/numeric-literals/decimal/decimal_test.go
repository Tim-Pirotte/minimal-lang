package decimalparser

import (
	"minimal/minimal-lang/built-in/ast"
	"minimal/minimal-lang/built-in/lexer"
	decimalmatcher "minimal/minimal-lang/built-in/matchers/decimal"
	"minimal/minimal-lang/built-in/parsers/prattparser"
	"testing"
)

func TestWhole(t *testing.T) {
	lexerSchema := lexer.NewScheme()
    astSchema := ast.NewSchema()

	intType := lexerSchema.NewTokenType(
        lexer.TokenTypeMetadata{NounPhrase: "a decimal integer literal", DebugName: "DecimalInt"},
    )

    wholeFloatExpType := lexerSchema.NewTokenType(
        lexer.TokenTypeMetadata{
            NounPhrase: "a decimal whole float literal with exponent", 
            DebugName: "DecimalWholeFloatExponent",
        },
    )

    wholeFloatNegExpType := lexerSchema.NewTokenType(
        lexer.TokenTypeMetadata{
            NounPhrase: "a decimal whole float literal with negative exponent", 
            DebugName: "DecimalWholeFloatNegativeExponent",
        },
    )

    floatType := lexerSchema.NewTokenType(
        lexer.TokenTypeMetadata{NounPhrase: "a decimal float literal", DebugName: "DecimalFloat"},
    )

    floatExpType := lexerSchema.NewTokenType(
        lexer.TokenTypeMetadata{
            NounPhrase: "a decimal float literal with exponent", 
            DebugName: "DecimalFloatExponent",
        },
    )

    floatNegativeExpType := lexerSchema.NewTokenType(
        lexer.TokenTypeMetadata{
            NounPhrase: "a decimal float literal with negative exponent", 
            DebugName: "DecimalFloatNegativeExponent",
        },
    )

	lexerSchema.AddMatcher(
        decimalmatcher.New(
            intType, 
            wholeFloatExpType, 
            wholeFloatNegExpType, 
            floatType, 
            floatExpType, 
            floatNegativeExpType,
        ),
    )

    intNodeType := astSchema.NewNodeType(
        &ast.StructNodeTypeMetadata{
            DisplayName: "a decimal integer literal", 
            DebugName: "DecimalInt", 
            ChildCount: 0,
        },
    )

    parser := NewWholeParser(intType, intNodeType)

    pp := prattparser.New(lexerSchema)
    pp.AddPrefix(parser)

    l := lexerSchema.Lex("14567894567898977984465446565545488778")
    result := pp.Parse(l, 0)

    t.Errorf("%+v %+v", result, parser.numbers)
}
