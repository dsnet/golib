// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"
import "bytes"
import "testing"
import "github.com/stretchr/testify/assert"

func TestBtoi(t *testing.T) {
	assert.Equal(t, uint(0), Btoi(false))
	assert.Equal(t, uint(1), Btoi(true))
}

func TestItob(t *testing.T) {
	assert.Equal(t, false, Itob(0))
	assert.Equal(t, true, Itob(1))
	assert.Equal(t, true, Itob(2))
	assert.Equal(t, true, Itob(MaxUint))
}

func TestGet(t *testing.T) {
	b := []byte{0x7b, 0x3a}

	assert.Equal(t, true, Get(b, 0))
	assert.Equal(t, true, Get(b, 1))
	assert.Equal(t, false, Get(b, 2))
	assert.Equal(t, true, Get(b, 3))
	assert.Equal(t, true, Get(b, 4))
	assert.Equal(t, true, Get(b, 5))
	assert.Equal(t, true, Get(b, 6))
	assert.Equal(t, false, Get(b, 7))

	assert.Equal(t, false, Get(b, 8))
	assert.Equal(t, true, Get(b, 9))
	assert.Equal(t, false, Get(b, 10))
	assert.Equal(t, true, Get(b, 11))
	assert.Equal(t, true, Get(b, 12))
	assert.Equal(t, true, Get(b, 13))
	assert.Equal(t, false, Get(b, 14))
	assert.Equal(t, false, Get(b, 15))
}

func TestGetN(t *testing.T) {
	assert.Equal(t, uint(0x00), GetN([]byte(nil), 0, 0))
	assert.Equal(t, uint(0x00), GetN([]byte{}, 0, 0))
	assert.Equal(t, uint(0x0f), GetN([]byte{0xaf}, 4, 0))
	assert.Equal(t, uint(0x0a), GetN([]byte{0xaf, 0x3a}, 4, 4))
	assert.Equal(t, uint(0x2b), GetN([]byte{0xaf, 0xb8}, 6, 2))
	assert.Equal(t, uint(0x125d), GetN([]byte{0xba, 0x64}, 13, 1))
	assert.Equal(t, uint(0x04), GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 8, 0))
	assert.Equal(t, uint(0x5c), GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 8, 8))
	assert.Equal(t, uint(0xeb), GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 8, 16))
	assert.Equal(t, uint(0x2d), GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 8, 24))
	assert.Equal(t, uint(0x15), GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 5, 17))
	assert.Equal(t, uint(0x1bd6b80), GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 26, 3))
}

func TestSet(t *testing.T) {
	b := []byte{0x00, 0xff}

	Set(b, true, 0)
	Set(b, true, 1)
	Set(b, false, 2)
	Set(b, true, 3)
	assert.Equal(t, []byte{0x0b, 0xff}, b)
	Set(b, true, 4)
	Set(b, true, 5)
	Set(b, true, 6)
	Set(b, false, 7)
	assert.Equal(t, []byte{0x7b, 0xff}, b)

	Set(b, false, 8)
	Set(b, true, 9)
	Set(b, false, 10)
	Set(b, true, 11)
	assert.Equal(t, []byte{0x7b, 0xfa}, b)
	Set(b, true, 12)
	Set(b, true, 13)
	Set(b, false, 14)
	Set(b, false, 15)
	assert.Equal(t, []byte{0x7b, 0x3a}, b)
}

func TestSetN(t *testing.T) {
	var b []byte

	b = []byte(nil)
	SetN(b, 0, 0, 0) // Should not crash

	b = []byte{0xaa}
	SetN(b, 0x0f, 4, 0)
	assert.Equal(t, []byte{0xaf}, b)

	b = []byte{0x55, 0x55}
	SetN(b, 0x0a, 4, 4)
	assert.Equal(t, []byte{0xa5, 0x55}, b)

	b = []byte{0x55, 0x55}
	SetN(b, 0x2b, 6, 2)
	assert.Equal(t, []byte{0xad, 0x55}, b)

	b = []byte{0x55, 0x55}
	SetN(b, 0x125d, 13, 1)
	assert.Equal(t, []byte{0xbb, 0x64}, b)

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0x04, 8, 0)
	assert.Equal(t, []byte{0x04, 0x55, 0x55, 0x55}, b)

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0x5c, 8, 8)
	assert.Equal(t, []byte{0x55, 0x5c, 0x55, 0x55}, b)

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0xeb, 8, 16)
	assert.Equal(t, []byte{0x55, 0x55, 0xeb, 0x55}, b)

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0x2d, 8, 24)
	assert.Equal(t, []byte{0x55, 0x55, 0x55, 0x2d}, b)

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0x15, 5, 17)
	assert.Equal(t, []byte{0x55, 0x55, 0x6b, 0x55}, b)

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0x1bd6b80, 26, 3)
	assert.Equal(t, []byte{0x05, 0x5c, 0xeb, 0x4d}, b)

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0xf1bd6b80, 26, 3)
	assert.Equal(t, []byte{0x05, 0x5c, 0xeb, 0x4d}, b)
}

func TestInvert(t *testing.T) {
	var b []byte

	b = []byte(nil)
	Invert(b)
	assert.Equal(t, []byte(nil), b)

	b = []byte{0xaa}
	Invert(b)
	assert.Equal(t, []byte{0x55}, b)

	b = []byte{0x7b}
	Invert(b)
	assert.Equal(t, []byte{0x84}, b)

	b = []byte{0xf3, 0xd1}
	Invert(b)
	assert.Equal(t, []byte{0xc, 0x2e}, b)

	b = []byte{0xff, 0xff, 0xff}
	Invert(b)
	assert.Equal(t, []byte{0x00, 0x00, 0x00}, b)
}

func TestCount(t *testing.T) {
	assert.Equal(t, 0, Count(nil))
	assert.Equal(t, 4, Count([]byte{0xaa}))
	assert.Equal(t, 6, Count([]byte{0x7b}))
	assert.Equal(t, 10, Count([]byte{0xf3, 0xd1}))
	assert.Equal(t, 24, Count([]byte{0xff, 0xff, 0xff}))
}

func TestCountByte(t *testing.T) {
	assert.Equal(t, 0, CountByte(0x00))
	assert.Equal(t, 3, CountByte(0x13))
	assert.Equal(t, 4, CountByte(0xf0))
	assert.Equal(t, 4, CountByte(0x0f))
	assert.Equal(t, 4, CountByte(0xaa))
	assert.Equal(t, 6, CountByte(0x7e))
	assert.Equal(t, 8, CountByte(0xff))
}

func TestCountUint(t *testing.T) {
	assert.Equal(t, 0, CountUint(0x0))
	assert.Equal(t, 2, CountUint(0x3))
	assert.Equal(t, 3, CountUint(0xe))
	assert.Equal(t, 5, CountUint(0x3d))
	assert.Equal(t, 8, CountUint(0xff))
	assert.Equal(t, 6, CountUint(0xa8d))
	assert.Equal(t, 10, CountUint(0x7fe))
	assert.Equal(t, 10, CountUint(0xa8df))
	assert.Equal(t, 16, CountUint(0xffff))
	assert.Equal(t, NumUintBits-2, CountUint(MaxUint>>2))
	assert.Equal(t, NumUintBits-1, CountUint(MaxUint>>1))
	assert.Equal(t, NumUintBits, CountUint(MaxUint))
}

func TestReverseByte(t *testing.T) {
	assert.Equal(t, byte(0x00), ReverseByte(0x00))
	assert.Equal(t, byte(0x37), ReverseByte(0xec))
	assert.Equal(t, byte(0x9c), ReverseByte(0x39))
	assert.Equal(t, byte(0x88), ReverseByte(0x11))
	assert.Equal(t, byte(0x8e), ReverseByte(0x71))
	assert.Equal(t, byte(0xaa), ReverseByte(0x55))
	assert.Equal(t, byte(0x3a), ReverseByte(0x5c))
	assert.Equal(t, byte(0x2f), ReverseByte(0xf4))
	assert.Equal(t, byte(0xd3), ReverseByte(0xcb))
	assert.Equal(t, byte(0xff), ReverseByte(0xff))
}

func TestReverseUint(t *testing.T) {
	assert.Equal(t, uint(0), ReverseUint(0))
	assert.Equal(t, MaxUint, ReverseUint(MaxUint))
	assert.Equal(t, MaxUint&(^uint(1)), ReverseUint(MaxUint>>1))
	assert.Equal(t, MaxUint&(^uint(3)), ReverseUint(MaxUint>>2))
	assert.Equal(t, MaxUint&(^uint(7)), ReverseUint(MaxUint>>3))
	assert.Equal(t, uint(0xb7<<uint(NumUintBits-8)), ReverseUint(0xed))
	assert.Equal(t, uint(0x7b3d5<<uint(NumUintBits-20)), ReverseUint(0xabcde))
}

func TestReverseUintN(t *testing.T) {
	assert.Equal(t, uint(0), ReverseUintN(MaxUint, 0))
	assert.Equal(t, MaxUint, ReverseUintN(MaxUint, NumUintBits))
	assert.Equal(t, MaxUint>>1, ReverseUintN(MaxUint, NumUintBits-1))
	assert.Equal(t, MaxUint>>2, ReverseUintN(MaxUint, NumUintBits-2))
	assert.Equal(t, MaxUint&(^uint(1)), ReverseUintN(MaxUint>>1, NumUintBits))
	assert.Equal(t, MaxUint&(^uint(3)), ReverseUintN(MaxUint>>2, NumUintBits))
	assert.Equal(t, MaxUint&(^uint(7)), ReverseUintN(MaxUint>>3, NumUintBits))
	assert.Equal(t, uint(0xb7), ReverseUintN(0xed, 8))
	assert.Equal(t, uint(0xb7)<<2, ReverseUintN(0xed, 10))
	assert.Equal(t, uint(0x7b3d5), ReverseUintN(0xabcde, 20))
	assert.Equal(t, uint(0x7b3d5), ReverseUintN(0xfabcde, 20))
	assert.Equal(t, uint(0x7b3d5)<<3, ReverseUintN(0xabcde, 23))
}

func TestWriteSameBit(t *testing.T) {
	var cnt int
	var err error
	var bb *Buffer

	for _, x := range []struct {
		b bool
		n int
	}{
		{false, 0},
		{true, 2},
		{false, 16},
		{true, NumUintBits},
		{false, 4321},
	} {
		bb = NewBuffer(nil)
		cnt, err = WriteSameBit(bb, x.b, x.n)
		assert.Equal(t, x.n, cnt)
		assert.Equal(t, nil, err)
		assert.Equal(t, int64(x.n), bb.BitsWritten())
		if x.b {
			assert.Equal(t, x.n, Count(bb.Bytes()))
		} else {
			assert.Equal(t, 0, Count(bb.Bytes()))
		}
	}
}

func TestReadBits(t *testing.T) {
	var val uint
	var cnt int
	var err error

	b := bytes.NewBuffer(nil)
	br := NewReader(b)

	val, cnt, err = ReadBits(br, 0)
	assert.Equal(t, uint(0), val)
	assert.Equal(t, 0, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, b.Len())

	val, cnt, err = ReadBits(br, 1)
	assert.Equal(t, uint(0), val)
	assert.Equal(t, 0, cnt)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, b.Len())

	b.Write([]byte{0xc9})
	assert.Equal(t, 1, b.Len())

	val, cnt, err = ReadBits(br, 3)
	assert.Equal(t, uint(1), val)
	assert.Equal(t, 3, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, b.Len())

	val, cnt, err = ReadBits(br, 8)
	assert.Equal(t, uint(0x19), val)
	assert.Equal(t, 5, cnt)
	assert.Equal(t, io.ErrUnexpectedEOF, err)
	assert.Equal(t, 0, b.Len())

	val, cnt, err = ReadBits(br, 8)
	assert.Equal(t, uint(0), val)
	assert.Equal(t, 0, cnt)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, b.Len())

	b.Write([]byte{0xeb, 0xad, 0xe2})
	assert.Equal(t, 3, b.Len())

	val, cnt, err = ReadBits(br, 7)
	assert.Equal(t, uint(0x6b), val)
	assert.Equal(t, 7, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, b.Len())

	val, cnt, err = ReadBits(br, 9)
	assert.Equal(t, uint(0x15b), val)
	assert.Equal(t, 9, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, b.Len())

	val, cnt, err = ReadBits(br, 8)
	assert.Equal(t, uint(0xe2), val)
	assert.Equal(t, 8, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, b.Len())

	val, cnt, err = ReadBits(br, 3)
	assert.Equal(t, uint(0), val)
	assert.Equal(t, 0, cnt)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, b.Len())
}

func TestWriteBits(t *testing.T) {
	var cnt int
	var err error

	b := new(faultyBuffer)
	bw := NewWriter(b)

	cnt, err = WriteBits(bw, 0x16, 5)
	assert.Equal(t, 5, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(b.Bytes()))

	cnt, err = WriteBits(bw, 0x0b, 5)
	assert.Equal(t, 5, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(b.Bytes()))

	cnt, err = WriteBits(bw, 0x2d, 6)
	assert.Equal(t, 6, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(b.Bytes()))

	b.fw = true

	cnt, err = WriteBits(bw, 0x1a6d, 13)
	assert.Equal(t, 7, cnt)
	assert.Equal(t, io.ErrShortWrite, err)
	assert.Equal(t, 2, len(b.Bytes()))

	b.fw = false

	cnt, err = WriteBits(bw, 0x1a7b1, 17)
	assert.Equal(t, 17, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, 5, len(b.Bytes()))

	assert.Equal(t, []byte{0x76, 0xb5, 0xed, 0xd8, 0xd3}, b.Bytes())
}
