package bpm

var Algorithms struct {
	AC  string
	BMH string
	KMP string
}

func init() {
	Algorithms.AC = "Aho-Corasick"
	Algorithms.BMH = "Boyer–Moore–Horspool"
	Algorithms.KMP = "Knuth–Morris–Pratt"
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
	FindAllOccurrences() (SearchResults, error)
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
}

type MultiplePatternSearchEngine interface {
	SearchEngine
	PreprocessPatterns(patterns [][]byte)
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
