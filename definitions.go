package bmp

import (
	"fmt"
)

var Algorithms struct {
	BMH string
	KMP string
	AC  string
}

func init() {
	Algorithms.BMH = "Boyer–Moore–Horspool"
	Algorithms.KMP = "Knuth–Morris–Pratt"
	Algorithms.AC = "Aho-Corasick"
}

type SearchResults map[string][]int64

func newSearchResults() SearchResults {
	srs := make(map[string][]int64)
	return srs
}

func (sr SearchResults) putOne(pattern []byte, position int64) SearchResults {
	sr[string(pattern)] = append(sr[string(pattern)], position)
	return sr
}

func (sr SearchResults) putMany(pattern []byte, positions []int64) SearchResults {
	sr[string(pattern)] = append(sr[string(pattern)], positions...)
	return sr
}

// Get returns positions of the given byte pattern.
func (sr SearchResults) Get(pattern []byte) []int64 {
	return sr[string(pattern)]
}

type SearchEngine interface {
	Name() string
	SetFile(filePath string) error
}

type SinglePatternSearchEngine interface {
	SearchEngine
	FindAllOccurrences(pattern []byte) (SearchResults, error)
}

type MultiplePatternSearchEngine interface {
	SearchEngine
	FindAllOccurrences(patterns [][]byte) (SearchResults, error)
}

func NewSearchEngine(alg string) (SearchEngine, error) {
	switch alg {
	case Algorithms.BMH:
		return &_BMHSearchEngine{name: alg}, nil
	case Algorithms.KMP:
		return &_KMPSearchEngine{name: alg}, nil
	case Algorithms.AC:
		return &_ACSearchEngine{name: alg}, nil
	default:
		return nil, fmt.Errorf("unknown algorithm")
	}
}
