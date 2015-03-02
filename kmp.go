package bmp

import (
	"fmt"
)

type _KMPSearchEngine struct {
	bfr *bufferedFileReader
}

func (kmp *_KMPSearchEngine) SetFile(fp string) (err error) {
	var bfr *bufferedFileReader
	if bfr, err = NewBufferedFileReader(fp); err != nil {
		return
	}
	kmp.bfr = bfr
	return
}

func computePrefix(pattern []byte) []int64 {
	pl := len(pattern)
	if pl == 1 {
		return []int64{-1}
	}
	var (
		pfx   []int64 = make([]int64, pl)
		count int64   = 0
		pos   int     = 2
	)
	pfx[0], pfx[1] = -1, 0
	for pos < pl {
		if pattern[pos-1] == pattern[count] {
			count++
			pfx[pos] = count
			pos++
		} else {
			if count > 0 {
				count = pfx[count]
			} else {
				pfx[pos] = 0
				pos++
			}
		}
	}
	return pfx
}

func (kmp *_KMPSearchEngine) FindAllOccurrences(pattern []byte) (srs SearchResults, err error) {
	dl, pl := kmp.bfr.FileSize(), int64(len(pattern))
	if pl > dl {
		err = fmt.Errorf("pattern not must be longer than data")
		return
	}
	var (
		i, m    int64
		indices []int64
		plm1    int64 = pl - 1
	)
	prefix := computePrefix(pattern)
	for m+i < dl {
		if pattern[i] == kmp.bfr.ReadByteAt(m+i) {
			if i == plm1 {
				indices = append(indices, m)
				m += i - prefix[i]
				if prefix[i] > -1 {
					i = prefix[i]
				} else {
					i = 0
				}
			} else {
				i++
			}
		} else {
			m += i - prefix[i]
			if prefix[i] > -1 {
				i = prefix[i]
			} else {
				i = 0
			}
		}
	}
	srs = newSearchResults().putMany(pattern, indices)
	return
}
