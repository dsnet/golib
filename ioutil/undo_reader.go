// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

// The MemoryReader allows a client to unread data back writing data to an
// internal stack. Reads will pop data from the internal stack before reading
// from the underlying reader.
//
// The MemoryReader provides similar functionality as the UndoReader, but is
// more complex. It solves certain cases where UndoReader is not sufficient.
type UndoReader struct {
	rd  io.Reader
	buf []byte // Reverse stack, which grows starting from the end
	off int    // The stack pointer, grows from the end
}

// Create a new UndoReader using buf as the internal stack.
func NewUndoReader(rd io.Reader, buf []byte) *UndoReader {
	return &UndoReader{rd, buf, len(buf)}
}

// If there is data on the internal stack, read from it first. Then, finish off
// the input buffer with a read to the underlying reader.
func (r *UndoReader) Read(data []byte) (cnt int, err error) {
	// Pop data off from internal undo stack
	if len(data) > 0 && len(r.buf) > r.off {
		cnt = copy(data, r.buf[r.off:])
		data = data[cnt:]
		r.off += cnt
	}

	// Top off the input buffer with actual Read call
	var rdCnt int
	rdCnt, err = r.rd.Read(data)
	cnt += rdCnt

	return cnt, err
}

// Unread some data by copying it to the internal stack and moving the read
// pointer backwards.
//
// Since the data read from the underlying reader is not recorded. This method
// does nothing to validate that the data being unread really was the data that
// was read in the first place.
func (r *UndoReader) UndoRead(data []byte) {
	if len(data) > r.off {
		// Undo stack is too small, allocate a larger one
		rdyBuf := r.buf[r.off:]
		needCnt := len(data) + len(rdyBuf)
		padCnt := int(8 + 0.25*float32(needCnt)) // Extra room for growth
		r.buf = make([]byte, padCnt, padCnt+needCnt)
		r.buf = append(r.buf, data...)
		r.buf = append(r.buf, rdyBuf...)
		r.off = padCnt
	} else {
		// Push onto undo stack
		offNew := r.off - len(data)
		copy(r.buf[offNew:r.off], data)
		r.off = offNew
	}
}
