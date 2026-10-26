package textsourcefinder

import "minimal/minimal-lang/built-in/substring"

type TextSourceFinder struct {
	sources []string
}

func New() *TextSourceFinder {
	return &TextSourceFinder{[]string{}}
}

func (t *TextSourceFinder) AddText(text string) {
	t.sources = append(t.sources, text)
}

func (t *TextSourceFinder) GetSourceOfText(text string) (string, bool) {
	for _, source := range t.sources {
		if substring.IsSubString(source, text) {
			return source, true
		}
	}

	return "", false
}
