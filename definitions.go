package bmp

import (
	"fmt"
)

var Algorithm struct {
	BMH string
	KMP string
	AC  string
}

func init() {
	Algorithm.BMH = "Boyer–Moore–Horspool"
	Algorithm.KMP = "Knuth–Morris–Pratt"
	Algorithm.AC = "Aho-Corasick"
}

type SearchResults map[string][]int64

func NewSearchResults() SearchResults {
	srs := make(map[string][]int64)
	return srs
}

func (sr SearchResults) putOne(r []byte, p int64) SearchResults {
	sr[string(r)] = append(sr[string(r)], p)
	return sr
}

func (sr SearchResults) putMany(r []byte, ps []int64) SearchResults {
	sr[string(r)] = append(sr[string(r)], ps...)
	return sr
}

func (sr SearchResults) Get(r []byte) []int64 {
	return sr[string(r)]
}

type SearchEngine interface {
	SetFile(fp string) error
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
	case Algorithm.BMH:
		return &_BMHSearchEngine{}, nil
	case Algorithm.KMP:
		return &_KMPSearchEngine{}, nil
	case Algorithm.AC:
		return &_ACSearchEngine{}, nil
	default:
		return nil, fmt.Errorf("unknown algorithm")
	}
}
