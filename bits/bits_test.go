// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"
import "testing"
import "github.com/stretchr/testify/assert"

// Extension of Buffer, that allows fine control of failures.
type faultyBuffer struct {
	Buffer
	nr, nw int  // Bits to read/write before failure.
	fr, fw bool // Can read/write fail?
}

func (fb *faultyBuffer) ReadBit() (val bool, err error) {
	if !fb.fr || fb.nr >= 1 {
		val, err = fb.Buffer.ReadBit()
	} else {
		err = io.ErrNoProgress
	}
	fb.nr -= int(Btoi(err == nil))
	return
}

func (fb *faultyBuffer) ReadBits(num int) (val uint, cnt int, err error) {
	if !fb.fr || fb.nr >= num {
		val, cnt, err = fb.Buffer.ReadBits(num)
	} else {
		num = min(max(fb.nr, 0), num)
		val, cnt, err = fb.Buffer.ReadBits(num)
		err = io.ErrNoProgress
	}
	fb.nr -= cnt
	return
}

func (fb *faultyBuffer) ReadByte() (val byte, err error) {
	if !fb.fr || fb.nr >= 8 {
		val, err = fb.Buffer.ReadByte()
	} else {
		err = io.ErrNoProgress
	}
	fb.nr -= int(8 * Btoi(err == nil))
	return
}

func (fb *faultyBuffer) Read(buf []byte) (cnt int, err error) {
	if !fb.fr || fb.nr >= 8*len(buf) {
		cnt, err = fb.Buffer.Read(buf)
	} else {
		buf = buf[:min(max(fb.nr, 0)/8, len(buf))]
		cnt, err = fb.Buffer.Read(buf)
		err = io.ErrNoProgress
	}
	fb.nr -= 8 * cnt
	return
}

func (fb *faultyBuffer) WriteBit(val bool) (err error) {
	if !fb.fw || fb.nw >= 1 {
		err = fb.Buffer.WriteBit(val)
	} else {
		err = io.ErrShortWrite
	}
	fb.nw -= int(Btoi(err == nil))
	return
}

func (fb *faultyBuffer) WriteBits(val uint, num int) (cnt int, err error) {
	if !fb.fw || fb.nw >= num {
		cnt, err = fb.Buffer.WriteBits(val, num)
	} else {
		num = min(max(fb.nw, 0), num)
		cnt, err = fb.Buffer.WriteBits(val, num)
		err = io.ErrShortWrite
	}
	fb.nw -= cnt
	return
}

func (fb *faultyBuffer) WriteByte(val byte) (err error) {
	if !fb.fw || fb.nw >= 8 {
		err = fb.Buffer.WriteByte(val)
	} else {
		err = io.ErrShortWrite
	}
	fb.nw -= int(8 * Btoi(err == nil))
	return
}

func (fb *faultyBuffer) Write(buf []byte) (cnt int, err error) {
	if !fb.fw || fb.nw >= 8*len(buf) {
		cnt, err = fb.Buffer.Write(buf)
	} else {
		buf = buf[:min(max(fb.nw, 0)/8, len(buf))]
		cnt, err = fb.Buffer.Write(buf)
		err = io.ErrShortWrite
	}
	fb.nw -= 8 * cnt
	return
}

// Helper test function that converts any empty byte slice to the nil slice so
// that equality checks work out fine.
func nb(buf []byte) []byte {
	if len(buf) == 0 {
		return nil
	}
	return buf
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func TestInterfaces(t *testing.T) {
	assert.Implements(t, (*BitsReader)(nil), NewReader(nil))
	assert.Implements(t, (*BitsWriter)(nil), NewWriter(nil))
	assert.Implements(t, (*BitsReader)(nil), NewBuffer(nil))
	assert.Implements(t, (*BitsWriter)(nil), NewBuffer(nil))
}
