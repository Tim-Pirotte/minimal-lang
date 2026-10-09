package numbers

// Maybe we should do all text processing during lexing
// so that the parses doesnt need to load source anymore
// That would mean we handle bases in the lexer
// We should just match the first char on a digit 0-9
// Then we get the next char and if the first char was 0 the next should be a base
// If not we just create a decimal int with value 0
// Otherwise for every following char we call a boolean function on it via an interface
// After int parsing we check if the next char is a . and if so we repeat the process to produce a float
// The base should not be respecified for a float.
// After parsing the second int we check for a floating point scientific notation prefix
// If that exists we check for - and then we parse another int
// Al of these should produce tokens. Per integer base a different token but one for float
// and one for scientific notation
// _ should be skipped during int parsing automatically but we shouldnt allow it at the end or multiple times
// After error production we should skip the invalid _'s.

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

	for c, ok := l.Get(pos); ok && '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c > asciiMax || c == '_' || c == '.'; c, ok = l.Get(pos) {
		pos++
	}

	return pos
}

func (i *NumberMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: i.tokenType, Value: l.GetNextN(length)})
}
