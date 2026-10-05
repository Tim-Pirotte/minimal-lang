package numbers

// TODO format:
// First char is digit?
// If so keep lexing undtil the next space
// We will only check the rules during parsing so other matchers
// Can also lex sequences that start with a digit but would be invalid numeric literals

// Base digits are case sensitive (uppercase), base specifiers too (lowercase)
// Do we allow utf-8 digits?

// TODO this wont work, it will consume everything else. We need some sort of valid range

import "minimal/minimal-lang/built-in/lexer"

const asciiMax = 127

type NumberMatcher struct {
	tokenType lexer.TokenType
}

func NewNumberMatcher(tt lexer.TokenType) *NumberMatcher {
	return &NumberMatcher{tt}
}

func (n *NumberMatcher) New(_ *lexer.Lexer) lexer.Matcher {
	return n
}

func (n *NumberMatcher) Match(l *lexer.Lexer) uint32 {
	if c, _ := l.Get(0); c < '0' || c > '9' {
		return 0
	}

	pos := uint32(1)

	for c, ok := l.Get(pos); ok && '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c > asciiMax || c == '_'; c, ok = l.Get(pos) {
		pos++
	}

	return pos
}

func (i *NumberMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: i.tokenType, Value: l.GetNextN(length)})
}
