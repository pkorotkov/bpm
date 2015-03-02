package bpm

import (
	"testing"

	. "github.com/pkorotkov/randombulk"
)

var (
	randomFilePath = "./test-data/data.bin"
	spatterns      = []string{
		"3Hs9Wnlo1Q5",
		"o0plMs-EgdC3#",
		"D0zWholeBr7q4W9",
		"Kw87_uiX2Y",
		"RWo-45vZl",
	}
)

func createRandomFile() (err error) {
	bulk := NewRandomBulk()
	incls := []Inclusion{
		NewInclusionFromString(spatterns[0], Frequencies.Sometimes),
		NewInclusionFromString(spatterns[1], Frequencies.Rarely),
		NewInclusionFromString(spatterns[2], Frequencies.Sometimes),
		NewInclusionFromString(spatterns[3], Frequencies.Frequently),
		NewInclusionFromString(spatterns[4], Frequencies.Rarely),
	}
	_, err = bulk.DumpToFile(randomFilePath, 50*1024, false, incls)
	return
}

func TestAlgorithms(t *testing.T) {
	var (
		err      error
		se       SearchEngine
		srs      SearchResults
		patterns [][]byte
	)
	for _, sp := range spatterns {
		patterns = append(patterns, []byte(sp))
	}
    // Aho-Corasick.
	se = NewSearchEngine(Algorithms.AC)
	if err = se.SetFile(randomFilePath); err != nil {
		t.Fatal(err)
	}
	ac := se.(MultiplePatternSearchEngine)
	ac.PreprocessPatterns(patterns)
	if srs, err = ac.FindAllOccurrences(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		fps := srs.Get(patterns[i])
		sps := samplePositions[spatterns[i]]
		for k, fp := range fps {
			if sp := sps[k]; fp != sp {
				t.Errorf("comparison is broken: expected %d but given with %d", sp, fp)
			}
		}
	}
}
