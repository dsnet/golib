// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"
import "bytes"
import "runtime"
import "testing"
import "github.com/stretchr/testify/assert"

func TestReader(t *testing.T) {
	type X struct {
		wa      bool
		wbn, wn int64
	}
	state := func(br *Reader) X {
		return X{
			br.ReadAligned(),
			br.BytesRead(), br.BitsRead(),
		}
	}

	var bit bool
	var dat byte
	var val uint
	var cnt int
	var err error

	b := bytes.NewBuffer(nil)
	br := NewReader(b)
	assert.Equal(t, X{true, 0, 0}, state(br))

	_, err = br.ReadBit()
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, X{true, 0, 0}, state(br))

	b.WriteByte(0x9d)

	bit, err = br.ReadBit()
	assert.Equal(t, true, bit)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, 1, 1}, state(br))

	bit, err = br.ReadBit()
	assert.Equal(t, false, bit)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, 1, 2}, state(br))

	val, cnt, err = br.ReadBits(4)
	assert.Equal(t, uint(0x07), val)
	assert.Equal(t, 4, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, 1, 6}, state(br))

	_, err = br.ReadByte()
	assert.Equal(t, ErrAlign, err)
	assert.Equal(t, X{false, 1, 6}, state(br))

	bit, err = br.ReadBit()
	assert.Equal(t, false, bit)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, 1, 7}, state(br))

	bit, err = br.ReadBit()
	assert.Equal(t, true, bit)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, 1, 8}, state(br))

	_, err = br.ReadByte()
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, X{true, 1, 8}, state(br))

	b.WriteByte(0xa7)

	dat, err = br.ReadByte()
	assert.Equal(t, byte(0xa7), dat)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, 2, 16}, state(br))

	_, err = br.ReadBit()
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, X{true, 2, 16}, state(br))

	b.WriteByte(0x8e)
	b.WriteByte(0xa3)

	bit, err = br.ReadBit()
	assert.Equal(t, false, bit)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, 3, 17}, state(br))

	val, cnt, err = br.ReadBits(21)
	assert.Equal(t, uint(0x51c7), val)
	assert.Equal(t, 15, cnt)
	assert.Equal(t, io.ErrUnexpectedEOF, err)
	assert.Equal(t, X{true, 4, 32}, state(br))

	_, cnt, err = br.ReadBits(3)
	assert.Equal(t, 0, cnt)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, X{true, 4, 32}, state(br))

	// Reset
	br.Reset(nil)
	assert.Equal(t, X{true, 0, 0}, state(br))
}

func BenchmarkReader(b *testing.B) {
	cnt := 1 << 20 // 1 MiB
	data := make([]byte, cnt)
	for i := range data {
		data[i] = 0x55
	}
	buf := NewBuffer(data)
	br := NewReader(buf)
	brr := BitReader(br)

	runtime.GC()
	b.ReportAllocs()
	b.SetBytes(int64(cnt))
	b.ResetTimer()

	for ni := 0; ni < b.N; ni++ {
		buf.ResetBuffer(data)
		br.Reset(buf)
		for bi := 0; bi < cnt; bi++ {
			brr.ReadBit()
			brr.ReadBit()
			brr.ReadBit()
			brr.ReadBit()
			brr.ReadBit()
			brr.ReadBit()
			brr.ReadBit()
			brr.ReadBit()
		}
	}
}
