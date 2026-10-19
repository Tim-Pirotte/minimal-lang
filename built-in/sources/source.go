package sources

import (
	"math"
	"minimal/minimal-lang/built-in/substring"
)

type Sources struct {
	Sources []Source
}

type Source struct {
	Name    string
	Content string
}

func New() *Sources {
	return &Sources{[]Source{}}
}

func (s *Sources) Add(source Source) bool {
	if len(source.Content) > math.MaxUint32 {
		return false
	}

	s.Sources = append(s.Sources, source)

	return true
}

func (s *Sources) GetSourceOfText(text string) (Source, bool) {
	for _, source := range s.Sources {
		if substring.IsSubString(source.Content, text) {
			return source, true
		}
	}
	
	return Source{}, false
}
