// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"
import "bytes"
import "testing"
import "github.com/stretchr/testify/assert"

func TestBtoi(t *testing.T) {
	assert.Equal(t, Btoi(false), uint(0))
	assert.Equal(t, Btoi(true), uint(1))
}

func TestItob(t *testing.T) {
	assert.Equal(t, Itob(0), false)
	assert.Equal(t, Itob(1), true)
	assert.Equal(t, Itob(2), true)
	assert.Equal(t, Itob(MaxUint), true)
}

func TestGet(t *testing.T) {
	b := []byte{0x7b, 0x3a}

	assert.Equal(t, Get(b, 0), true)
	assert.Equal(t, Get(b, 1), true)
	assert.Equal(t, Get(b, 2), false)
	assert.Equal(t, Get(b, 3), true)
	assert.Equal(t, Get(b, 4), true)
	assert.Equal(t, Get(b, 5), true)
	assert.Equal(t, Get(b, 6), true)
	assert.Equal(t, Get(b, 7), false)

	assert.Equal(t, Get(b, 8), false)
	assert.Equal(t, Get(b, 9), true)
	assert.Equal(t, Get(b, 10), false)
	assert.Equal(t, Get(b, 11), true)
	assert.Equal(t, Get(b, 12), true)
	assert.Equal(t, Get(b, 13), true)
	assert.Equal(t, Get(b, 14), false)
	assert.Equal(t, Get(b, 15), false)
}

func TestGetN(t *testing.T) {
	assert.Equal(t, GetN([]byte(nil), 0, 0), uint(0x00))
	assert.Equal(t, GetN([]byte{0xaf}, 4, 0), uint(0x0f))
	assert.Equal(t, GetN([]byte{0xaf, 0x3a}, 4, 4), uint(0x0a))
	assert.Equal(t, GetN([]byte{0xaf, 0xb8}, 6, 2), uint(0x2b))
	assert.Equal(t, GetN([]byte{0xba, 0x64}, 13, 1), uint(0x125d))
	assert.Equal(t, GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 8, 0), uint(0x04))
	assert.Equal(t, GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 8, 8), uint(0x5c))
	assert.Equal(t, GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 8, 16), uint(0xeb))
	assert.Equal(t, GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 8, 24), uint(0x2d))
	assert.Equal(t, GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 5, 17), uint(0x15))
	assert.Equal(t, GetN([]byte{0x04, 0x5c, 0xeb, 0x2d}, 26, 3), uint(0x1bd6b80))
}

func TestSet(t *testing.T) {
	b := []byte{0x00, 0xff}

	Set(b, true, 0)
	Set(b, true, 1)
	Set(b, false, 2)
	Set(b, true, 3)
	assert.Equal(t, b, []byte{0x0b, 0xff})
	Set(b, true, 4)
	Set(b, true, 5)
	Set(b, true, 6)
	Set(b, false, 7)
	assert.Equal(t, b, []byte{0x7b, 0xff})

	Set(b, false, 8)
	Set(b, true, 9)
	Set(b, false, 10)
	Set(b, true, 11)
	assert.Equal(t, b, []byte{0x7b, 0xfa})
	Set(b, true, 12)
	Set(b, true, 13)
	Set(b, false, 14)
	Set(b, false, 15)
	assert.Equal(t, b, []byte{0x7b, 0x3a})
}

func TestSetN(t *testing.T) {
	var b []byte

	b = []byte(nil)
	SetN(b, 0, 0, 0) // Should not crash

	b = []byte{0xaa}
	SetN(b, 0x0f, 4, 0)
	assert.Equal(t, b, []byte{0xaf})

	b = []byte{0x55, 0x55}
	SetN(b, 0x0a, 4, 4)
	assert.Equal(t, b, []byte{0xa5, 0x55})

	b = []byte{0x55, 0x55}
	SetN(b, 0x2b, 6, 2)
	assert.Equal(t, b, []byte{0xad, 0x55})

	b = []byte{0x55, 0x55}
	SetN(b, 0x125d, 13, 1)
	assert.Equal(t, b, []byte{0xbb, 0x64})

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0x04, 8, 0)
	assert.Equal(t, b, []byte{0x04, 0x55, 0x55, 0x55})

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0x5c, 8, 8)
	assert.Equal(t, b, []byte{0x55, 0x5c, 0x55, 0x55})

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0xeb, 8, 16)
	assert.Equal(t, b, []byte{0x55, 0x55, 0xeb, 0x55})

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0x2d, 8, 24)
	assert.Equal(t, b, []byte{0x55, 0x55, 0x55, 0x2d})

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0x15, 5, 17)
	assert.Equal(t, b, []byte{0x55, 0x55, 0x6b, 0x55})

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0x1bd6b80, 26, 3)
	assert.Equal(t, b, []byte{0x05, 0x5c, 0xeb, 0x4d})

	b = []byte{0x55, 0x55, 0x55, 0x55}
	SetN(b, 0xf1bd6b80, 26, 3)
	assert.Equal(t, b, []byte{0x05, 0x5c, 0xeb, 0x4d})
}

func TestCount(t *testing.T) {
	var zeros, ones int

	zeros, ones = Count(nil)
	assert.Equal(t, zeros, 0)
	assert.Equal(t, ones, 0)

	zeros, ones = Count([]byte{0xaa})
	assert.Equal(t, zeros, 4)
	assert.Equal(t, ones, 4)

	zeros, ones = Count([]byte{0x7b})
	assert.Equal(t, zeros, 2)
	assert.Equal(t, ones, 6)

	zeros, ones = Count([]byte{0xf3, 0xd1})
	assert.Equal(t, zeros, 6)
	assert.Equal(t, ones, 10)

	zeros, ones = Count([]byte{0xff, 0xff, 0xff})
	assert.Equal(t, zeros, 0)
	assert.Equal(t, ones, 24)
}

func TestInvert(t *testing.T) {
	var b []byte

	b = []byte(nil)
	Invert(b)
	assert.Equal(t, b, []byte(nil))

	b = []byte{0xaa}
	Invert(b)
	assert.Equal(t, b, []byte{0x55})

	b = []byte{0x7b}
	Invert(b)
	assert.Equal(t, b, []byte{0x84})

	b = []byte{0xf3, 0xd1}
	Invert(b)
	assert.Equal(t, b, []byte{0xc, 0x2e})

	b = []byte{0xff, 0xff, 0xff}
	Invert(b)
	assert.Equal(t, b, []byte{0x00, 0x00, 0x00})
}

func TestReverseUint(t *testing.T) {
	assert.Equal(t, ReverseUint(0), uint(0))
	assert.Equal(t, ReverseUint(MaxUint), MaxUint)
	assert.Equal(t, ReverseUint(MaxUint>>1), MaxUint&(^uint(1)))
	assert.Equal(t, ReverseUint(MaxUint>>2), MaxUint&(^uint(3)))
	assert.Equal(t, ReverseUint(MaxUint>>3), MaxUint&(^uint(7)))
	assert.Equal(t, ReverseUint(0xed), uint(0xb7<<uint(NumUintBits-8)))
	assert.Equal(t, ReverseUint(0xabcde), uint(0x7b3d5<<uint(NumUintBits-20)))
}

func TestReverseUintN(t *testing.T) {
	assert.Equal(t, ReverseUintN(MaxUint, 0), uint(0))
	assert.Equal(t, ReverseUintN(MaxUint, NumUintBits), MaxUint)
	assert.Equal(t, ReverseUintN(MaxUint, NumUintBits-1), MaxUint>>1)
	assert.Equal(t, ReverseUintN(MaxUint, NumUintBits-2), MaxUint>>2)
	assert.Equal(t, ReverseUintN(MaxUint>>1, NumUintBits), MaxUint&(^uint(1)))
	assert.Equal(t, ReverseUintN(MaxUint>>2, NumUintBits), MaxUint&(^uint(3)))
	assert.Equal(t, ReverseUintN(MaxUint>>3, NumUintBits), MaxUint&(^uint(7)))
	assert.Equal(t, ReverseUintN(0xed, 8), uint(0xb7))
	assert.Equal(t, ReverseUintN(0xed, 10), uint(0xb7)<<2)
	assert.Equal(t, ReverseUintN(0xabcde, 20), uint(0x7b3d5))
	assert.Equal(t, ReverseUintN(0xfabcde, 20), uint(0x7b3d5))
	assert.Equal(t, ReverseUintN(0xabcde, 23), uint(0x7b3d5)<<3)
}

func TestReadBits(t *testing.T) {
	var val uint
	var cnt int
	var err error

	b := bytes.NewBuffer(nil)
	br := NewReader(b)

	val, cnt, err = ReadBits(br, 0)
	assert.Equal(t, val, uint(0))
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, nil)
	assert.Equal(t, b.Len(), 0)

	val, cnt, err = ReadBits(br, 1)
	assert.Equal(t, val, uint(0))
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, b.Len(), 0)

	b.Write([]byte{0xc9})
	assert.Equal(t, b.Len(), 1)

	val, cnt, err = ReadBits(br, 3)
	assert.Equal(t, val, uint(1))
	assert.Equal(t, cnt, 3)
	assert.Equal(t, err, nil)
	assert.Equal(t, b.Len(), 0)

	val, cnt, err = ReadBits(br, 8)
	assert.Equal(t, val, uint(0x19))
	assert.Equal(t, cnt, 5)
	assert.Equal(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, b.Len(), 0)

	val, cnt, err = ReadBits(br, 8)
	assert.Equal(t, val, uint(0))
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, b.Len(), 0)

	b.Write([]byte{0xeb, 0xad, 0xe2})
	assert.Equal(t, b.Len(), 3)

	val, cnt, err = ReadBits(br, 7)
	assert.Equal(t, val, uint(0x6b))
	assert.Equal(t, cnt, 7)
	assert.Equal(t, err, nil)
	assert.Equal(t, b.Len(), 2)

	val, cnt, err = ReadBits(br, 9)
	assert.Equal(t, val, uint(0x15b))
	assert.Equal(t, cnt, 9)
	assert.Equal(t, err, nil)
	assert.Equal(t, b.Len(), 1)

	val, cnt, err = ReadBits(br, 8)
	assert.Equal(t, val, uint(0xe2))
	assert.Equal(t, cnt, 8)
	assert.Equal(t, err, nil)
	assert.Equal(t, b.Len(), 0)

	val, cnt, err = ReadBits(br, 3)
	assert.Equal(t, val, uint(0))
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, b.Len(), 0)
}

func TestWriteBits(t *testing.T) {
	var cnt int
	var err error

	b := new(buffer)
	bw := NewWriter(b)

	cnt, err = WriteBits(bw, 0x16, 5)
	assert.Equal(t, cnt, 5)
	assert.Equal(t, err, nil)
	assert.Equal(t, b.Len(), 0)

	cnt, err = WriteBits(bw, 0x0b, 5)
	assert.Equal(t, cnt, 5)
	assert.Equal(t, err, nil)
	assert.Equal(t, b.Len(), 1)

	cnt, err = WriteBits(bw, 0x2d, 6)
	assert.Equal(t, cnt, 6)
	assert.Equal(t, err, nil)
	assert.Equal(t, b.Len(), 2)

	b.fail = true

	cnt, err = WriteBits(bw, 0x1a6d, 13)
	assert.Equal(t, cnt, 7)
	assert.Equal(t, err, io.ErrShortWrite)
	assert.Equal(t, b.Len(), 2)

	b.fail = false

	cnt, err = WriteBits(bw, 0x1a7b1, 17)
	assert.Equal(t, cnt, 17)
	assert.Equal(t, err, nil)
	assert.Equal(t, b.Len(), 5)

	assert.Equal(t, b.Bytes(), []byte{0x76, 0xb5, 0xed, 0xd8, 0xd3})
}
