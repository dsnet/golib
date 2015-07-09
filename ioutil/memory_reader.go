// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"
import "os"
import "errors"

// The MemoryReader contains a ring buffer that acts like a sliding window. It
// records the last N bytes of data read from the underlying reader. It provides
// a seeker-like interface that allows clients to seek backwards within the
// sliding window. The MemoryReader is useful for wrapping streamed readers,
// where a need arises to either read slightly older data or to "undo" a read.
//
// For most use cases, the UndoReader is probably sufficient and much simpler
// to use. The major use case of MemoryReader over UndoReader is when the data
// read is not available, and thus cannot be fed into the UndoRead method.
// Since MemoryReader records all data that is read, this is not an issue.
type MemoryReader struct {
	rd    io.Reader
	buf   []byte
	rdPtr int64
	wrPtr int64
}

// Create a new MemoryReader using buf as the sliding window.
func NewMemoryReader(rd io.Reader, buf []byte) *MemoryReader {
	return &MemoryReader{rd, buf, 0, 0}
}

// If the read pointer has been placed within the ring buffer, read from it
// first. Then, finish off the input buffer with a read to the underlying
// reader. All data read from the underlying reader will be copied into the
// internal ring buffer.
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

	// Top off the input buffer with actual Read call
	var rdCnt int
	rdCnt, err = r.rd.Read(data)
	data = data[:rdCnt]
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

	return cnt, err
}

// The MemoryReader provides a seek-like interface, but there are strict
// restrictions on valid operations. First, seeking relative to the end is not
// supported since the position of the EOF is not known. Secondly, seeks can
// only be made that place the read pointer somewhere inside the internal
// ring buffer. Seeking too far backwards or forwards will cause an error. The
// absolute position returned by seek is relative to when the MemoryReader
// was created.
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

// Get the offsets that are seekable. Three values are returned: the lowest
// position that can be seeked to, the current position, and the highest
// position that can be seeked to. The current position is guaranteed to be
// in between the low and high offsets inclusively. The range of offsets that a
// relative seek can make can be computed as the range [posLo-pos, posHi-pos].
func (r *MemoryReader) SeekOffsets() (posLo int64, pos int64, posHi int64) {
	posLo = r.wrPtr - int64(len(r.buf))
	if posLo < 0 {
		posLo = 0
	}
	return posLo, r.rdPtr, r.wrPtr
}
