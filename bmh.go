package bpm

import (
	"fmt"
)

type _BMHSearchEngine struct {
	*_BaseEngine
	pattern []byte
}

func (bmh *_BMHSearchEngine) PreprocessPattern(pattern []byte) {
	bmh.pattern = pattern
}

func (bmh *_BMHSearchEngine) FindAllOccurrences() (srs SearchResults, err error) {
	dl, pl := bmh.bfr.FileSize(), int64(len(bmh.pattern))
	if pl > dl {
		err = fmt.Errorf("pattern must not be longer than data")
		return
	}
	var (
		indices  []int64
		badChars []int64 = make([]int64, 256)
		index    int64   = 0
		i        int64   = 0
		// Difference in data and pattern length.
		ld int64 = dl - pl
		// Last pattern's byte position.
		lpbp int64 = pl - 1
	)
	for i = 0; i < 256; i++ {
		badChars[i] = pl
	}
	for i = 0; i < lpbp; i++ {
		badChars[bmh.pattern[i]] = lpbp - i
	}
	// Start search.
	for index <= ld {
	innerLoop:
		for i = lpbp; bmh.bfr.ReadByteAt(index+i) == bmh.pattern[i]; i-- {
			if i == 0 {
				indices = append(indices, index)
				break innerLoop
			}
		}
		index += badChars[bmh.bfr.ReadByteAt(index+lpbp)]
	}
	srs = newSearchResults().putMany(bmh.pattern, indices)
	return
}
