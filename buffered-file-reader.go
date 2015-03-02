package bmp

import (
	"os"
)

type bufferedFileReader struct {
	baseFile            *os.File
	fileSize            int64
	buffer              []byte
	bufferSize          int64
	bufferPivotPosition int64
}

func NewBufferedFileReader(fp string) (bfr *bufferedFileReader, err error) {
	return NewBufferedFileReaderWithSize(fp, 64*1024*1024)
}

func NewBufferedFileReaderWithSize(fp string, bs int64) (bfr *bufferedFileReader, err error) {
	var (
		bf  *os.File
		fi  os.FileInfo
		buf []byte
	)
	if bf, err = os.Open(fp); err != nil {
		return
	}
	if fi, err = bf.Stat(); err != nil {
		return
	}
	buf = make([]byte, bs)
	if _, err = bf.Read(buf); err != nil {
		return
	}
	bfr = &bufferedFileReader{bf, fi.Size(), buf, bs, 0}
	return
}

func (bfr *bufferedFileReader) FileSize() int64 {
	return bfr.fileSize
}

func (bfr *bufferedFileReader) ReadByteAt(offset int64) byte {
	if offset >= bfr.bufferPivotPosition && offset < (bfr.bufferPivotPosition+bfr.bufferSize) {
		return bfr.buffer[offset-bfr.bufferPivotPosition]
	}
	if n, _ := bfr.baseFile.ReadAt(bfr.buffer, offset); n == 0 {
		panic("(*bufferedFileReader).ReadByteAt: unexpected zero bytes read")
	}
	bfr.bufferPivotPosition = offset
	return bfr.buffer[0]
}

func (bfr *bufferedFileReader) Close() error {
	return bfr.baseFile.Close()
}
