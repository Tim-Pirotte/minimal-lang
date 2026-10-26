package modules

import (
	"math"
	"minimal/minimal-lang/built-in/substring"
	"unique"
)

type RetrievalStatus int

const (
	NotFound RetrievalStatus = iota
	FirstRetrieved
	AlreadyRetrieved
	LibraryDoesNotExist
)

type Library interface {
	Retrieve([]unique.Handle[string]) (Module, bool)
	GetName() string
}

type Libraries struct {
	libraries map[unique.Handle[string]]Library
	modules map[[]unique.Handle[string]]Module
}

type Module struct {
	Name    []unique.Handle[string]
	Content string
}

func New() *Libraries {
	return &Libraries{map[unique.Handle[string]]Library{}, []Module{}}
}

func (s *Libraries) Retrieve(libraryName unique.Handle[string], path []unique.Handle[string]) (Module, RetrievalStatus) {
	
}

func (s *Libraries) Add(module Module) bool {
	if len(module.Content) > math.MaxUint32 {
		return false
	}

	s.modules = append(s.modules, module)

	return true
}
