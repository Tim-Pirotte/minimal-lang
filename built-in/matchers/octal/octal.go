package octalmatcher

import "minimal/minimal-lang/built-in/lexer"

type OctalMatcher struct {
	whole       lexer.TokenType
	fraction    lexer.TokenType
	matchedType lexer.TokenType
}

func New(whole, fraction lexer.TokenType) *OctalMatcher {
    return &OctalMatcher{whole, fraction, whole}
}

func (o *OctalMatcher) New(_ *lexer.Lexer) lexer.Matcher {
    return &OctalMatcher{o.whole, o.fraction, o.whole}
}

func (o *OctalMatcher) Match(l *lexer.Lexer) uint32 {
    if c, _ := l.Get(0); c != '0' {
        return 0
    }

    if c, ok := l.Get(1); !ok || c != 'o' {
        return 0
    }

    pos := skipDigits(l, 2)

    if c, ok := l.Get(pos); !ok || c != '.' {
        o.matchedType = o.whole

        return pos
    }

    o.matchedType = o.fraction

    return skipDigits(l, pos + 1)
}

func (o *OctalMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: o.matchedType, Value: l.GetNextN(length)})
}

func skipDigits(l *lexer.Lexer, pos uint32) uint32 {
    for c, ok := l.Get(pos); ok && ('0' <= c && c <= '7' || c == '_'); c, ok = l.Get(pos) {
        pos++
	}

    return pos
}
