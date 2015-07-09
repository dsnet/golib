// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"

type Reader struct {
	rd   io.ByteReader
	val  byte
	mask byte
	cnt  int64
}

// New bit Reader.
//
// In order to satisfy the error conditions specified in the bits.BitReader
// interface, the io.ByteReader must return an error if and only if a byte was
// not read.
func NewReader(rd io.ByteReader) *Reader {
	return &Reader{rd: rd, mask: 0x00}
}

// Read a single bit. The first bit returned is the LSB of the current byte.
// A byte is read from underlying reader on an as-needed basis.
func (br *Reader) ReadBit() (bit bool, err error) {
	if br.mask == 0x00 {
		if br.val, err = br.rd.ReadByte(); err != nil {
			return
		}
		br.mask = 0x01
		br.cnt++
	}
	bit = br.val&br.mask > 0
	br.mask <<= 1
	return
}

// Read multiple bits.
func (br *Reader) ReadBits(num int) (val uint, cnt int, err error) {
	// The logic here is based on the generic bits.ReadBits function. It is
	// manually unfolded here for efficiency. Furthermore, the call to
	// br.ReadBit is inlined.
	var bit bool
	for cnt = 0; cnt < num; cnt++ {
		// Nearly identical to the br.ReadBit function
		if br.mask == 0x00 {
			if br.val, err = br.rd.ReadByte(); err != nil {
				if err == io.EOF && cnt > 0 {
					err = io.ErrUnexpectedEOF
				}
				return
			}
			br.mask = 0x01
			br.cnt++
		}
		bit = br.val&br.mask > 0
		br.mask <<= 1

		val |= (Btoi(bit) << uint(cnt))
	}
	return
}

// Read a single byte.
// The internal offset must be aligned to a byte.
func (br *Reader) ReadByte() (val byte, err error) {
	if br.mask != 0x00 {
		return 0, ErrAlign
	}
	if val, err = br.rd.ReadByte(); err != nil {
		return
	}
	br.cnt++
	return
}

// Is the stream currently at a byte boundary?
func (br *Reader) ReadAligned() bool { return br.mask == 0x00 }

// Number of bytes read from underlying ByteReader.
func (br *Reader) BytesRead() int64 { return br.cnt }

// Number of bits read.
func (br *Reader) BitsRead() int64 {
	bits := 8 * br.cnt
	for mask := br.mask; mask > 0; mask <<= 1 {
		bits--
	}
	return bits
}

// Reset the reader.
func (br *Reader) Reset(rd io.ByteReader) {
	br.rd, br.val, br.mask, br.cnt = rd, 0x00, 0x00, 0
}
