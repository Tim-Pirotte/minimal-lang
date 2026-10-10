package binarymatcher

import "minimal/minimal-lang/built-in/lexer"

type BinaryMatcher struct {
	whole       lexer.TokenType
    fraction    lexer.TokenType
    matchedType lexer.TokenType
}

func New(whole, fraction lexer.TokenType) *BinaryMatcher {
    return &BinaryMatcher{whole, fraction, whole}
}

func (b *BinaryMatcher) New(_ *lexer.Lexer) lexer.Matcher {
    return &BinaryMatcher{b.whole, b.fraction, b.whole}
}

func (b *BinaryMatcher) Match(l *lexer.Lexer) uint32 {
    if c, _ := l.Get(0); c != '0' {
        return 0
    }

    if c, ok := l.Get(1); !ok || c != 'b' {
        return 0
    }

    pos := skipDigits(l, 2)

    if c, ok := l.Get(pos); !ok || c != '.' {
        b.matchedType = b.whole

        return pos
    }

    b.matchedType = b.fraction

    return skipDigits(l, pos + 1)
}

func (b *BinaryMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: b.matchedType, Value: l.GetNextN(length)})
}

func skipDigits(l *lexer.Lexer, pos uint32) uint32 {
    for c, ok := l.Get(pos); ok && (c == '0' || c == '1' || c == '_'); c, ok = l.Get(pos) {
        pos++
	}

    return pos
}
