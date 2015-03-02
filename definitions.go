package bpm

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

type _BaseEngine struct {
	name string
	bfr  *bufferedFileReader
}

func (be *_BaseEngine) Name() string {
	return be.name
}

func (be *_BaseEngine) SetFile(fp string) (err error) {
	var bfr *bufferedFileReader
	if bfr, err = newBufferedFileReader(fp); err != nil {
		return
	}
	be.bfr = bfr
	return
}

type SinglePatternSearchEngine interface {
	SearchEngine
	PreprocessPattern(pattern []byte)
	FindAllOccurrences() (SearchResults, error)
}

type MultiplePatternSearchEngine interface {
	SearchEngine
	PreprocessPatterns(patterns [][]byte)
	FindAllOccurrences() (SearchResults, error)
}

func NewSearchEngine(alg string) SearchEngine {
	switch alg {
	case Algorithms.BMH:
		return &_BMHSearchEngine{_BaseEngine: &_BaseEngine{name: alg}}
	case Algorithms.KMP:
		return &_KMPSearchEngine{_BaseEngine: &_BaseEngine{name: alg}}
	case Algorithms.AC:
		return &_ACSearchEngine{_BaseEngine: &_BaseEngine{name: alg}}
	default:
		return nil
	}
}
