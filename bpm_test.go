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
		srs      SearchResults
		patterns [][]byte
	)
	for _, sp := range spatterns {
		patterns = append(patterns, []byte(sp))
	}
	// Test Aho-Corasick.
	ac := NewSearchEngine(Algorithms.AC).(MultiplePatternSearchEngine)
	if err = ac.SetFile(randomFilePath); err != nil {
		t.Fatal(err)
	}
	ac.PreprocessPatterns(patterns)
	if srs, err = ac.FindAllOccurrences(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		fps, sps := srs.Get(patterns[i]), samplePositions[spatterns[i]]
		for k, fp := range fps {
			if sp := sps[k]; fp != sp {
				t.Errorf("[AC] position comparison is broken: expected %d but encountered %d", sp, fp)
			}
		}
	}
	// Test Boyer–Moore–Horspool.
	bmh := NewSearchEngine(Algorithms.BMH).(SinglePatternSearchEngine)
	if err = bmh.SetFile(randomFilePath); err != nil {
		t.Fatal(err)
	}
	bmh.PreprocessPattern(patterns[2])
	if srs, err = bmh.FindAllOccurrences(); err != nil {
		t.Fatal(err)
	}
	fps, sps := srs.Get(patterns[2]), samplePositions[spatterns[2]]
	for k, fp := range fps {
		if sp := sps[k]; fp != sp {
			t.Errorf("[BMH] position comparison is broken: expected %d but encountered %d", sp, fp)
		}
	}
	// Test Knuth–Morris–Pratt.
	kmp := NewSearchEngine(Algorithms.KMP).(SinglePatternSearchEngine)
	if err = kmp.SetFile(randomFilePath); err != nil {
		t.Fatal(err)
	}
	kmp.PreprocessPattern(patterns[3])
	if srs, err = kmp.FindAllOccurrences(); err != nil {
		t.Fatal(err)
	}
	fps, sps = srs.Get(patterns[3]), samplePositions[spatterns[3]]
	for k, fp := range fps {
		if sp := sps[k]; fp != sp {
			t.Errorf("[KMP] position comparison is broken: expected %d but encountered %d", sp, fp)
		}
	}
}
