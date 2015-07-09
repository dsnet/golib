// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"
import "testing"
import "github.com/stretchr/testify/assert"

func TestBuffer(t *testing.T) {
	type X struct {
		ra, wa bool
		rn, wn int64
		buf    []byte
	}
	state := func(bb *Buffer) X {
		return X{
			bb.ReadAligned(), bb.WriteAligned(),
			bb.BitsRead(), bb.BitsWritten(),
			nb(bb.Bytes()),
		}
	}

	var bit bool
	var dat byte
	var val uint
	var buf [1024]byte
	var cnt int
	var err error

	bb := NewBuffer(nil)
	assert.Equal(t, X{true, true, 0, 0, nil}, state(bb))

	_, err = bb.ReadBit()
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, X{true, true, 0, 0, nil}, state(bb))

	err = bb.WriteBit(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, false, 0, 1, []byte{0x01}}, state(bb))

	bit, err = bb.ReadBit()
	assert.Equal(t, true, bit)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, false, 1, 1, []byte{0x01}}, state(bb))

	_, err = bb.ReadBit()
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, X{false, false, 1, 1, []byte{0x01}}, state(bb))

	err = bb.WriteBit(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, false, 1, 2, []byte{0x03}}, state(bb))

	err = bb.WriteBit(false)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, false, 1, 3, []byte{0x03}}, state(bb))

	bit, err = bb.ReadBit()
	assert.Equal(t, true, bit)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, false, 2, 3, []byte{0x03}}, state(bb))

	cnt, err = bb.WriteBits(0x0b, 5)
	assert.Equal(t, 5, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, true, 2, 8, []byte{0x5b}}, state(bb))

	val, cnt, err = bb.ReadBits(6)
	assert.Equal(t, uint(0x16), val)
	assert.Equal(t, 6, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, true, 8, 8, nil}, state(bb))

	err = bb.WriteByte(0xd3)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, true, 8, 16, []byte{0xd3}}, state(bb))

	cnt, err = bb.Write([]byte{0xc3, 0x1c, 0x3b})
	assert.Equal(t, 3, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, true, 8, 40, []byte{0xd3, 0xc3, 0x1c, 0x3b}}, state(bb))

	cnt, err = bb.WriteBits(0x7bfe3, 19)
	assert.Equal(t, 19, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, false, 8, 59, []byte{0xd3, 0xc3, 0x1c, 0x3b, 0xe3, 0xbf, 0x07}}, state(bb))

	err = bb.WriteByte(0xff)
	assert.Equal(t, ErrAlign, err)
	assert.Equal(t, X{true, false, 8, 59, []byte{0xd3, 0xc3, 0x1c, 0x3b, 0xe3, 0xbf, 0x07}}, state(bb))

	cnt, err = bb.Write([]byte{0xff, 0x00})
	assert.Equal(t, 0, cnt)
	assert.Equal(t, ErrAlign, err)
	assert.Equal(t, X{true, false, 8, 59, []byte{0xd3, 0xc3, 0x1c, 0x3b, 0xe3, 0xbf, 0x07}}, state(bb))

	cnt, err = bb.Read(buf[:2])
	assert.Equal(t, 2, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte{0xd3, 0xc3}, nb(buf[:cnt]))
	assert.Equal(t, X{true, false, 24, 59, []byte{0x1c, 0x3b, 0xe3, 0xbf, 0x07}}, state(bb))

	dat, err = bb.ReadByte()
	assert.Equal(t, byte(0x1c), dat)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, false, 32, 59, []byte{0x3b, 0xe3, 0xbf, 0x07}}, state(bb))

	val, cnt, err = bb.ReadBits(13)
	assert.Equal(t, uint(0x33b), val)
	assert.Equal(t, 13, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, false, 45, 59, []byte{0xe3, 0xbf, 0x07}}, state(bb))

	cnt, err = bb.Read(buf[:])
	assert.Equal(t, 0, cnt)
	assert.Equal(t, ErrAlign, err)
	assert.Equal(t, X{false, false, 45, 59, []byte{0xe3, 0xbf, 0x07}}, state(bb))

	_, err = bb.ReadByte()
	assert.Equal(t, ErrAlign, err)
	assert.Equal(t, X{false, false, 45, 59, []byte{0xe3, 0xbf, 0x07}}, state(bb))

	val, cnt, err = bb.ReadBits(21)
	assert.Equal(t, uint(0x3dff), val)
	assert.Equal(t, 14, cnt)
	assert.Equal(t, io.ErrUnexpectedEOF, err)
	assert.Equal(t, X{false, false, 59, 59, []byte{0x07}}, state(bb))

	_, cnt, err = bb.ReadBits(3)
	assert.Equal(t, 0, cnt)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, X{false, false, 59, 59, []byte{0x07}}, state(bb))

	cnt, err = bb.WriteBits(0xabcde, 20)
	assert.Equal(t, 20, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, false, 59, 79, []byte{0xf7, 0xe6, 0x55}}, state(bb))

	val, cnt, err = bb.ReadBits(5)
	assert.Equal(t, uint(0x1e), val)
	assert.Equal(t, 5, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, false, 64, 79, []byte{0xe6, 0x55}}, state(bb))

	cnt, err = bb.Read(buf[:])
	assert.Equal(t, 1, cnt)
	assert.Equal(t, ErrAlign, err)
	assert.Equal(t, []byte{0xe6}, nb(buf[:cnt]))
	assert.Equal(t, X{true, false, 72, 79, []byte{0x55}}, state(bb))

	_, err = bb.ReadByte()
	assert.Equal(t, ErrAlign, err)
	assert.Equal(t, X{true, false, 72, 79, []byte{0x55}}, state(bb))

	err = bb.WriteBit(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, true, 72, 80, []byte{0xd5}}, state(bb))

	cnt, err = bb.Read(buf[:])
	assert.Equal(t, 1, cnt)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, []byte{0xd5}, nb(buf[:cnt]))
	assert.Equal(t, X{true, true, 80, 80, nil}, state(bb))

	_, err = bb.ReadByte()
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, X{true, true, 80, 80, nil}, state(bb))

	// Reset
	bb.Reset()
	assert.Equal(t, X{true, true, 0, 0, nil}, state(bb))

	cnt, err = bb.Read(buf[:])
	assert.Equal(t, 0, cnt)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, []byte(nil), nb(buf[:cnt]))
	assert.Equal(t, X{true, true, 0, 0, nil}, state(bb))

	// Reset with data
	bb.ResetData([]byte{0xab, 0xcd, 0xef})
	assert.Equal(t, X{true, true, 0, 24, []byte{0xab, 0xcd, 0xef}}, state(bb))

	cnt, err = bb.Read(buf[:])
	assert.Equal(t, 3, cnt)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, []byte{0xab, 0xcd, 0xef}, nb(buf[:cnt]))
	assert.Equal(t, X{true, true, 24, 24, nil}, state(bb))
}

func BenchmarkBufferWriter(b *testing.B) {
	cnt := 1 << 20 // 1 MiB
	bb := NewBuffer(nil)

	b.SetBytes(int64(cnt))
	b.ResetTimer()

	for ni := 0; ni < b.N; ni++ {
		for bi := 0; bi < cnt; bi++ {
			bb.WriteBit(true)
			bb.WriteBit(false)
			bb.WriteBit(true)
			bb.WriteBit(false)
			bb.WriteBit(true)
			bb.WriteBit(false)
			bb.WriteBit(true)
			bb.WriteBit(false)
		}
	}
}

func BenchmarkBufferReader(b *testing.B) {
	cnt := 1 << 20 // 1 MiB
	data := make([]byte, cnt)
	for i := range data {
		data[i] = 0x55
	}
	bb := NewBuffer(data)

	b.SetBytes(int64(cnt))
	b.ResetTimer()

	for ni := 0; ni < b.N; ni++ {
		for bi := 0; bi < cnt; bi++ {
			bb.ReadBit()
			bb.ReadBit()
			bb.ReadBit()
			bb.ReadBit()
			bb.ReadBit()
			bb.ReadBit()
			bb.ReadBit()
			bb.ReadBit()
		}
	}
}
