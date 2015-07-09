// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"

// A Buffer is a variable-sized buffer that satisfies both bits.BitsReader and
// bits.BitsWriter. The zero value for Buffer is an empty buffer ready to use.
type Buffer struct {
	buf []byte
	off int

	rdMask byte
	wrMask byte
}

// New bit buffer.
func NewBuffer(data []byte) *Buffer {
	return &Buffer{buf: data}
}

// Read a bit in LSB order.
func (b *Buffer) ReadBit() (bit bool, err error) {
	if b.off == len(b.buf) && b.rdMask == b.wrMask {
		err = io.EOF
		return
	}
	if b.rdMask == 0x00 {
		b.rdMask = 0x01
		b.off++
	}
	bit = b.buf[b.off-1]&b.rdMask > 0
	b.rdMask <<= 1
	return
}

// Read multiple bits.
func (b *Buffer) ReadBits(num int) (val uint, cnt int, err error) {
	// The logic here is based on the generic bits.ReadBits function. It is
	// manually unfolded here for efficiency. Furthermore, the call to
	// b.ReadBit is inlined.
	var bit bool
	for cnt = 0; cnt < num; cnt++ {
		// Nearly identical to the b.ReadBit function
		if b.off == len(b.buf) && b.rdMask == b.wrMask {
			err = io.EOF
			if err == io.EOF && cnt > 0 {
				err = io.ErrUnexpectedEOF
			}
			return
		}
		if b.rdMask == 0x00 {
			b.rdMask = 0x01
			b.off++
		}
		bit = b.buf[b.off-1]&b.rdMask > 0
		b.rdMask <<= 1

		val |= (Btoi(bit) << uint(cnt))
	}
	return
}

// Read a single byte.
// The internal read offset must be byte-aligned.
func (b *Buffer) ReadByte() (val byte, err error) {
	if b.rdMask != 0x00 {
		return val, ErrAlign
	}
	if b.off == len(b.buf) {
		if b.wrMask != 0x00 {
			return val, ErrAlign
		}
		return val, io.EOF
	}
	val = b.buf[b.off]
	b.off++
	return
}

// Read multiple bytes.
//
// The internal read offset must be byte-aligned.
// If the EOF is hit, but the write offset is not byte-aligned, then
// bits.ErrAlign will be returned.
func (b *Buffer) Read(buf []byte) (cnt int, err error) {
	if b.rdMask != 0x00 {
		return cnt, ErrAlign
	}
	buf2, err2 := b.buf, io.EOF
	if b.wrMask != 0x00 {
		buf2, err2 = buf2[:len(buf2)-1], ErrAlign
	}
	cnt = copy(buf, buf2)
	b.off += cnt
	if cnt != len(buf) {
		err = err2
	}
	return
}

// Is the read offset aligned to a byte boundary?
func (b *Buffer) ReadAligned() bool { return b.rdMask == 0x00 }

// Total number of bits read.
func (b *Buffer) BitsRead() int64 {
	cnt := b.off * 8
	for mask := b.rdMask; mask > 0; mask <<= 1 {
		cnt--
	}
	return int64(cnt)
}

// Write a bit in LSB order.
func (b *Buffer) WriteBit(bit bool) (err error) {
	if b.wrMask == 0x00 {
		b.wrMask = 0x01
		b.buf = append(b.buf, 0x00)
	}
	if bit {
		b.buf[len(b.buf)-1] |= b.wrMask
	}
	b.wrMask <<= 1
	return
}

// Write multiple bits.
func (b *Buffer) WriteBits(val uint, num int) (cnt int, err error) {
	// The logic here is based on the generic bits.WriteBits function. It is
	// manually unfolded here for efficiency. Furthermore, the call to
	// b.WriteBit is inlined.
	var bit bool
	for cnt = 0; cnt < num; cnt++ {
		bit = Itob(val & (1 << uint(cnt)))

		// Nearly identical to the b.WriteBit function
		if b.wrMask == 0x00 {
			b.wrMask = 0x01
			b.buf = append(b.buf, 0x00)
		}
		if bit {
			b.buf[len(b.buf)-1] |= b.wrMask
		}
		b.wrMask <<= 1
	}
	return
}

// Write a single byte.
// The internal write offset must be byte-aligned.
func (b *Buffer) WriteByte(val byte) (err error) {
	if b.wrMask != 0x00 {
		return ErrAlign
	}
	b.buf = append(b.buf, val)
	return
}

// Write multiple bytes.
// The internal write offset must be byte-aligned.
func (b *Buffer) Write(buf []byte) (cnt int, err error) {
	if b.wrMask != 0x00 {
		return cnt, ErrAlign
	}
	b.buf = append(b.buf, buf...)
	cnt = len(buf)
	return
}

// Is the write offset aligned to a byte boundary?
func (b *Buffer) WriteAligned() bool { return b.wrMask == 0x00 }

// Total number of bits written.
func (b *Buffer) BitsWritten() int64 {
	cnt := len(b.buf) * 8
	for mask := b.wrMask; mask > 0; mask <<= 1 {
		cnt--
	}
	return int64(cnt)
}

// Return a slice of the contents of the unread portion of the buffer.
//
// If buffer is not read aligned, the first byte will include extra bits.
// If buffer is not write aligned, the last byte will include extra bits.
func (b *Buffer) Bytes() []byte {
	if b.rdMask != 0x00 { // Not read aligned
		return b.buf[b.off-1:]
	}
	return b.buf[b.off:]
}

// Reset the buffer so that it has no content.
func (b *Buffer) Reset() {
	b.buf, b.off = b.buf[:0], 0
	b.rdMask, b.wrMask = 0x00, 0x00
}

// Reset the buffer so that it uses the input data.
func (b *Buffer) ResetData(data []byte) {
	b.Reset()
	b.buf = data
}
