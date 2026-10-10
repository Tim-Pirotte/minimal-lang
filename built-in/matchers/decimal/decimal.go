package decimalmatcher

import "minimal/minimal-lang/built-in/lexer"

// TODO don't match on leading '_'s or '.'s or 'e's

type DecimalMatcher struct {
	whole                    lexer.TokenType
    wholeExponent            lexer.TokenType
    wholeNegativeExponent    lexer.TokenType
    fraction                 lexer.TokenType
    fractionExponent         lexer.TokenType
    fractionNegativeExponent lexer.TokenType
    matchedType              lexer.TokenType
}

func New(
    whole, 
    wholeExponent, 
    wholeNegativeExponent, 
    fraction, 
    fractionExponent, 
    fractionNegativeExponent lexer.TokenType,
) *DecimalMatcher {
	return &DecimalMatcher{
        whole, 
        wholeExponent, 
        wholeNegativeExponent, 
        fraction, 
        fractionExponent, 
        fractionNegativeExponent,
        whole,
    }
}

func (d *DecimalMatcher) New(_ *lexer.Lexer) lexer.Matcher {
	return &DecimalMatcher{
        d.whole, 
        d.wholeExponent, 
        d.wholeNegativeExponent, 
        d.fraction, 
        d.fractionExponent, 
        d.fractionNegativeExponent,
        d.whole,
    }
}

func (d *DecimalMatcher) Match(l *lexer.Lexer) uint32 {
	pos := skipDigits(l, 0)

    if c, ok := l.Get(pos); !ok || c != '.' {
        if ok && c == 'e' {
            pos++

            if c, ok = l.Get(pos); c == '-' {
                pos++
                d.matchedType = d.wholeNegativeExponent
            } else {
                d.matchedType = d.wholeExponent
            }

            pos = skipDigits(l, pos)

            return pos
        }

        d.matchedType = d.whole

        return pos
    }

    pos = skipDigits(l, pos + 1)

    if c, ok := l.Get(pos); !ok || c != 'e' {
        d.matchedType = d.fraction

        return pos
    }

    pos++

    if c, ok := l.Get(pos); ok && c == '-' {
        pos++
        d.matchedType = d.fractionNegativeExponent
    } else {
        d.matchedType = d.fractionExponent
    }

    pos = skipDigits(l, pos)

	return pos
}

func (d *DecimalMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: d.matchedType, Value: l.GetNextN(length)})
}

func skipDigits(l *lexer.Lexer, pos uint32) uint32 {
    for c, ok := l.Get(pos); ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
        pos++
	}

    return pos
}
