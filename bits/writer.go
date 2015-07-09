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
func NewWriter(wr io.ByteWriter) *Writer {
	return &Writer{wr: wr, mask: 0x01}
}

// Write a bit. The first bit written is the LSB of the current byte.
// If internal offset is byte-aligned, then a byte is flushed.
func (bw *Writer) WriteBit(val bool) (err error) {
	if val {
		bw.val |= bw.mask
	}
	bw.mask <<= 1
	if bw.mask == 0x00 {
		bw.mask = 0x01
		if err = bw.WriteByte(bw.val); err != nil {
			bw.val, bw.mask = bw.val&0x7f, 0x80 // Revert write
			return
		}
		bw.val = 0x00
	}
	return
}

// Write a single byte.
// The internal offset must be byte-aligned.
func (bw *Writer) WriteByte(val byte) (err error) {
	if !bw.ByteAligned() {
		return ErrAlign
	}
	if err = bw.wr.WriteByte(val); err != nil {
		return
	}
	bw.cnt++
	return
}

// Is the stream currently at a byte boundary?
func (bw *Writer) ByteAligned() bool { return bw.mask == 0x01 }

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
