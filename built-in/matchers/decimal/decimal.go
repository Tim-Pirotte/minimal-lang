package decimal

import "minimal/minimal-lang/built-in/lexer"

type DecimalMatcher struct {
	whole       lexer.TokenType
    fraction    lexer.TokenType
    matchedType lexer.TokenType
}

func NewDecimalMatcher(whole, fraction lexer.TokenType) *DecimalMatcher {
	return &DecimalMatcher{whole, fraction, whole}
}

func (d *DecimalMatcher) New(_ *lexer.Lexer) lexer.Matcher {
	return d
}

func (d *DecimalMatcher) Match(l *lexer.Lexer) uint32 {
	pos := uint32(0)
    previousWasUnderScore := false

    c, ok := l.Get(pos)

	for ; ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
		if c == '_' && previousWasUnderScore {
            // TODO error
        }

        previousWasUnderScore = c == '_'
        pos++
	}

    if previousWasUnderScore {
        // TODO error
    }

    if !ok || c != '.' {
        d.matchedType = d.whole

        return pos
    }

    posBefore := pos
    previousWasUnderScore = false

    for ; ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
		if c == '_' && previousWasUnderScore {
            // TODO error
        }

        previousWasUnderScore = c == '_'
        pos++
	}

    if previousWasUnderScore {
        // TODO error
    }

    if pos == posBefore {
        // TODO error
    }

    d.matchedType = d.fraction

	return pos
}

func (d *DecimalMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: d.matchedType, Value: l.GetNextN(length)})
}
