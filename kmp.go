package bpm

import (
	"fmt"
)

type _KMPSearchEngine struct {
	*_BaseEngine
	pattern []byte
	prefix  []int64
}

func (kmp *_KMPSearchEngine) PreprocessPattern(pattern []byte) {
	kmp.pattern = pattern
	kmp.computePrefix()
}

func (kmp *_KMPSearchEngine) computePrefix() {
	pl := len(kmp.pattern)
	if pl == 1 {
		kmp.prefix = []int64{-1}
		return
	}
	var (
		count int64 = 0
		pos   int   = 2
	)
	kmp.prefix = make([]int64, pl)
	kmp.prefix[0], kmp.prefix[1] = -1, 0
	for pos < pl {
		if kmp.pattern[pos-1] == kmp.pattern[count] {
			count++
			kmp.prefix[pos] = count
			pos++
		} else {
			if count > 0 {
				count = kmp.prefix[count]
			} else {
				kmp.prefix[pos] = 0
				pos++
			}
		}
	}
	return
}

func (kmp *_KMPSearchEngine) FindAllOccurrences() (srs SearchResults, err error) {
	dl, pl := kmp.bfr.FileSize(), int64(len(kmp.pattern))
	if pl > dl {
		err = fmt.Errorf("pattern not must be longer than data")
		return
	}
	var (
		i, m    int64
		indices []int64
		plm1    int64 = pl - 1
	)
	for m+i < dl {
		if kmp.pattern[i] == kmp.bfr.ReadByteAt(m+i) {
			if i == plm1 {
				indices = append(indices, m)
				m += i - kmp.prefix[i]
				if kmp.prefix[i] > -1 {
					i = kmp.prefix[i]
				} else {
					i = 0
				}
			} else {
				i++
			}
		} else {
			m += i - kmp.prefix[i]
			if kmp.prefix[i] > -1 {
				i = kmp.prefix[i]
			} else {
				i = 0
			}
		}
	}
	srs = newSearchResults().putMany(kmp.pattern, indices)
	return
}
