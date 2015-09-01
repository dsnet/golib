// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"

type Writer struct {
	wr   io.ByteWriter
	val  byte
	mask byte
	cnt  int64
}

// New bit Writer.
//
// In order to satisfy the error conditions specified in the bits.BitWriter
// interface, the io.ByteWriter must return an error if and only if a byte was
// not written.
func NewWriter(wr io.ByteWriter) *Writer {
	return &Writer{wr: wr, mask: 0x01}
}

// Write a bit. The first bit written is the LSB of the current byte.
// If internal offset is byte-aligned, then a byte is flushed.
func (bw *Writer) WriteBit(bit bool) (err error) {
	if bit {
		bw.val |= bw.mask
	}
	bw.mask <<= 1
	if bw.mask == 0x00 {
		bw.mask = 0x01
		if err = bw.wr.WriteByte(bw.val); err != nil {
			bw.val, bw.mask = bw.val&0x7f, 0x80 // Revert write
			return
		}
		bw.val = 0x00
		bw.cnt++
	}
	return
}

// Write multiple bits.
func (bw *Writer) WriteBits(val uint, num int) (cnt int, err error) {
	// The logic here is based on the generic bits.WriteBits function. It is
	// manually unfolded here for efficiency. Furthermore, the call to
	// bw.WriteBit is inlined.
	var bit bool
	for cnt = 0; cnt < num; cnt++ {
		bit = Itob(val & (1 << uint(cnt)))

		// Nearly identical to the bw.WriteBit function.
		if bit {
			bw.val |= bw.mask
		}
		bw.mask <<= 1
		if bw.mask == 0x00 {
			bw.mask = 0x01
			if err = bw.wr.WriteByte(bw.val); err != nil {
				bw.val, bw.mask = bw.val&0x7f, 0x80 // Revert write
				return
			}
			bw.val = 0x00
			bw.cnt++
		}
	}
	return
}

// Write a single byte.
// The internal offset must be byte-aligned.
func (bw *Writer) WriteByte(val byte) (err error) {
	if bw.mask != 0x01 {
		return ErrAlign
	}
	if err = bw.wr.WriteByte(val); err != nil {
		return
	}
	bw.cnt++
	return
}

// Is the stream currently at a byte boundary?
func (bw *Writer) WriteAligned() bool { return bw.mask == 0x01 }

// Number of bytes written to underlying ByteWriter.
func (bw *Writer) BytesWritten() int64 { return bw.cnt }

// Number of bits written in total.
func (bw *Writer) BitsWritten() int64 {
	bits := 8 * bw.cnt
	for mask := bw.mask; mask > 1; mask >>= 1 {
		bits++
	}
	return bits
}

// Reset the writer.
func (bw *Writer) Reset(wr io.ByteWriter) {
	bw.wr, bw.val, bw.mask, bw.cnt = wr, 0x00, 0x01, 0
}
