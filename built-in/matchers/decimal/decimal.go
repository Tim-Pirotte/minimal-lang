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
    previousWasUnderScore := false

    c, ok := l.Get(pos)
    posBefore := pos

	for ; ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
        if c == '_' {
            if pos == posBefore {
                panic("integer starts with _")
            } else if previousWasUnderScore {
                // TODO error
                // Errors should be accumulated for the consume
                panic("already _ in whole part")
            }
        }

        previousWasUnderScore = c == '_'
        pos++
	}

    if previousWasUnderScore {
        // TODO error
        panic("whole ended on _")
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

            posBefore = pos
            previousWasUnderScore = false

            for ; ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
                if c == '_' {
                    if pos == posBefore {
                        panic("whole exponent part starts with _")
                    } else if previousWasUnderScore {
                        panic("already _ in whole exponent part")
                    }
                }

                previousWasUnderScore = c == '_'
                pos++
            }

            if previousWasUnderScore {
                panic("whole exponent ended on _")
            }

            if pos == posBefore {
                panic("nothing after e or - (whole)")
            }

            return pos
        }

        d.matchedType = d.whole

        return pos
    }

    pos++
    posBefore = pos
    previousWasUnderScore = false

    for c, ok = l.Get(pos); ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
		if c == '_' {
            if pos == posBefore {
                panic("fractional parts starts with _")
            } else if previousWasUnderScore {
                panic("already _ in fractional part")
            }
        }

        previousWasUnderScore = c == '_'
        pos++
	}

    if previousWasUnderScore {
        panic("fractional ended on _")
    }

    if pos == posBefore {
        panic("nothing after .")
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

    posBefore = pos
    previousWasUnderScore = false

    for ; ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
		if c == '_' {
            if pos == posBefore {
                panic("exponent part starts with _")
            } else if previousWasUnderScore {
                panic("already _ in exponent part")
            }
        }

        previousWasUnderScore = c == '_'
        pos++
	}

    if previousWasUnderScore {
        panic("exponent ended on _")
    }

    if pos == posBefore {
        panic("nothing after e or -")
    }

	return pos
}

func (d *DecimalMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: d.matchedType, Value: l.GetNextN(length)})
}
