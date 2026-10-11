package hexmatcher

import "minimal/minimal-lang/built-in/lexer"

type HexMatcher struct {
	whole                    lexer.TokenType
	wholeExponent            lexer.TokenType
	wholeNegativeExponent    lexer.TokenType
	fraction                 lexer.TokenType
	fractionExponent         lexer.TokenType
	fractionNegativeExponent lexer.TokenType
	matchedType              lexer.TokenType
}

func (h *HexMatcher) New(_ *lexer.Lexer) lexer.Matcher {
	return &HexMatcher{
        h.whole, 
        h.wholeExponent, 
        h.wholeNegativeExponent, 
        h.fraction, 
        h.fractionExponent, 
        h.fractionNegativeExponent,
        h.whole,
    }
}

func (h *HexMatcher) Match(l *lexer.Lexer) uint32 {
	if c, _ := l.Get(0); c != '0' {
        return 0
    }

    if c, ok := l.Get(1); !ok || c != 'x' {
        return 0
    }
    
    pos := skipHexDigits(l, 2)

    if c, ok := l.Get(pos); !ok || c != '.' {
        if !ok || c != 'p' {
            h.matchedType = h.whole

            return pos
        }

        pos++

        if c, ok = l.Get(pos); c == '-' {
            pos++
            h.matchedType = h.wholeNegativeExponent
        } else {
            h.matchedType = h.wholeExponent
        }

        pos = skipDecimalDigits(l, pos)

        return pos
    }

    pos = skipHexDigits(l, pos + 1)

    if c, ok := l.Get(pos); !ok || c != 'p' {
        h.matchedType = h.fraction

        return pos
    }

    pos++

    if c, ok := l.Get(pos); ok && c == '-' {
        pos++
        h.matchedType = h.fractionNegativeExponent
    } else {
        h.matchedType = h.fractionExponent
    }

    pos = skipDecimalDigits(l, pos)

	return pos
}

func (h *HexMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: h.matchedType, Value: l.GetNextN(length)})
}

func skipHexDigits(l *lexer.Lexer, pos uint32) uint32 {
    for c, ok := l.Get(pos); ok && ('0' <= c && c <= '9' || 'A' <= c && c <= 'F' || c == '_'); c, ok = l.Get(pos) {
        pos++
	}

    return pos
}

func skipDecimalDigits(l *lexer.Lexer, pos uint32) uint32 {
    for c, ok := l.Get(pos); ok && ('0' <= c && c <= '9' || c == '_'); c, ok = l.Get(pos) {
        pos++
	}

    return pos
}
