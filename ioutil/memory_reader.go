// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"
import "os"
import "errors"

type MemoryReader struct {
	rd    io.Reader
	buf   []byte
	rdPtr int64
	wrPtr int64
}

func NewMemoryReader(rd io.Reader, buf []byte) *MemoryReader {
	return &MemoryReader{rd, buf, 0, 0}
}

func (r *MemoryReader) Read(data []byte) (cnt int, err error) {
	// If there is data to read from ring buffer, use it
	for len(data) > 0 && r.wrPtr > r.rdPtr {
		bufLen := int64(len(r.buf))
		offLo, offHi := r.rdPtr%bufLen, r.wrPtr%bufLen
		if offHi <= offLo {
			offHi = bufLen
		}
		cpyCnt := copy(data, r.buf[offLo:offHi])
		data = data[cpyCnt:]
		cnt += cpyCnt
		r.rdPtr += int64(cpyCnt)
	}

	// Top off the input buffer with actual Read() call
	if len(data) > 0 {
		var rdCnt int
		rdCnt, err = r.rd.Read(data)
		cnt += rdCnt

		// Write data to ring buffer
		if len(data) > len(r.buf) {
			skipCnt := len(data) - len(r.buf)
			data = data[skipCnt:]
			r.rdPtr += int64(skipCnt)
			r.wrPtr += int64(skipCnt)
		}
		for len(data) > 0 {
			off := r.wrPtr % int64(len(r.buf))
			cpyCnt := copy(r.buf[off:], data)
			data = data[cpyCnt:]
			r.rdPtr += int64(cpyCnt)
			r.wrPtr += int64(cpyCnt)
		}
	}
	return cnt, err
}

func (r *MemoryReader) Seek(offset int64, whence int) (pos int64, err error) {
	switch whence {
	case os.SEEK_SET:
		pos = offset
	case os.SEEK_CUR:
		pos = offset + r.rdPtr
	default:
		return 0, errors.New("ioutil.MemoryReader.Seek: invalid whence")
	}
	if posLo, _, posHi := r.SeekOffsets(); pos < posLo || pos > posHi {
		return 0, errors.New("ioutil.MemoryReader.Seek: buffer out of bounds")
	}
	r.rdPtr = pos
	return
}

func (r *MemoryReader) SeekOffsets() (int64, int64, int64) {
	ptrLo := r.wrPtr - int64(len(r.buf))
	if ptrLo < 0 {
		ptrLo = 0
	}
	return ptrLo, r.rdPtr, r.wrPtr
}
