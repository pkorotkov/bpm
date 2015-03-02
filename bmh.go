package bmp

import (
	"fmt"
)

type _BMHSearchEngine struct {
	name string
	bfr  *bufferedFileReader
}

func (bmh *_BMHSearchEngine) Name() string {
	return bmh.name
}

func (bmh *_BMHSearchEngine) SetFile(fp string) (err error) {
	var bfr *bufferedFileReader
	if bfr, err = NewBufferedFileReader(fp); err != nil {
		return
	}
	bmh.bfr = bfr
	return
}

func (bmh *_BMHSearchEngine) FindAllOccurrences(pattern []byte) (srs SearchResults, err error) {
	dl, pl := bmh.bfr.FileSize(), int64(len(pattern))
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
		badChars[pattern[i]] = lpbp - i
	}
	// Start search.
	for index <= ld {
	innerLoop:
		for i = lpbp; bmh.bfr.ReadByteAt(index+i) == pattern[i]; i-- {
			if i == 0 {
				indices = append(indices, index)
				break innerLoop
			}
		}
		index += badChars[bmh.bfr.ReadByteAt(index+lpbp)]
	}
	srs = newSearchResults().putMany(pattern, indices)
	return
}
