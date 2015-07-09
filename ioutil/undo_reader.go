// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

type UndoReader struct {
	rd  io.Reader
	buf []byte // Reverse stack, whichs grows starting from the end
	off int    // The stack pointer, grows from the end
}

func NewUndoReader(rd io.Reader, buf []byte) *UndoReader {
	return &UndoReader{rd, buf, len(buf)}
}

func (r *UndoReader) Read(data []byte) (cnt int, err error) {
	// Pop data off from internal undo stack
	if len(r.buf) > r.off {
		cnt = copy(data, r.buf[r.off:])
		data = data[cnt:]
		r.off += cnt
	}

	// Top off the input buffer with actual Read() call
	if len(data) > 0 {
		var rdCnt int
		rdCnt, err = r.rd.Read(data)
		cnt += rdCnt
	}
	return cnt, err
}

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
