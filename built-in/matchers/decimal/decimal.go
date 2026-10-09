package decimalmatcher

import "minimal/minimal-lang/built-in/lexer"

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
	return d
}

func (d *DecimalMatcher) Match(l *lexer.Lexer) uint32 {
	pos := uint32(0)

    c, ok := l.Get(pos)

	for ; ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
        pos++
	}

    if !ok || c != '.' {
        if ok && c == 'e' {
            pos++

            if c, ok = l.Get(pos); c == '-' {
                pos++
                c, ok = l.Get(pos)

                d.matchedType = d.wholeNegativeExponent
            } else {
                d.matchedType = d.wholeExponent
            }

            for ; ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
                pos++
            }

            return pos
        }

        d.matchedType = d.whole

        return pos
    }

    pos++

    for c, ok = l.Get(pos); ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
        pos++
	}

    if !ok || c != 'e' {
        d.matchedType = d.fraction

        return pos
    }

    pos++

    if c, ok = l.Get(pos); c == '-' {
        pos++
        c, ok = l.Get(pos)
    
        d.matchedType = d.fractionNegativeExponent
    } else {
        d.matchedType = d.fractionExponent
    }

    for ; ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
        pos++
	}

	return pos
}

func (d *DecimalMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: d.matchedType, Value: l.GetNextN(length)})
}
