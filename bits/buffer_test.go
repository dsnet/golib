// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"
import "testing"
import "github.com/stretchr/testify/assert"

// Helper test function that converts any empty byte slice to the nil slice so
// that equality checks work out fine.
func nb(buf []byte) []byte {
	if len(buf) == 0 {
		return nil
	}
	return buf
}

func TestBuffer(t *testing.T) {
	var bit bool
	var val uint
	var cnt int
	var err error

	bb := NewBuffer(nil)
	assert.Equal(t, true, bb.ReadAligned())
	assert.Equal(t, true, bb.WriteAligned())
	assert.Equal(t, 0, bb.BitsRead())
	assert.Equal(t, 0, bb.BitsWritten())
	assert.Equal(t, []byte(nil), nb(bb.Bytes()))

	_, err = bb.ReadBit()
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, true, bb.ReadAligned())
	assert.Equal(t, true, bb.WriteAligned())
	assert.Equal(t, 0, bb.BitsRead())
	assert.Equal(t, 0, bb.BitsWritten())
	assert.Equal(t, []byte(nil), nb(bb.Bytes()))

	err = bb.WriteBit(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, bb.ReadAligned())
	assert.Equal(t, false, bb.WriteAligned())
	assert.Equal(t, 0, bb.BitsRead())
	assert.Equal(t, 1, bb.BitsWritten())
	assert.Equal(t, []byte{0x01}, nb(bb.Bytes()))

	bit, err = bb.ReadBit()
	assert.Equal(t, true, bit)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, bb.ReadAligned())
	assert.Equal(t, false, bb.WriteAligned())
	assert.Equal(t, 1, bb.BitsRead())
	assert.Equal(t, 1, bb.BitsWritten())
	assert.Equal(t, []byte{0x01}, nb(bb.Bytes()))

	_, err = bb.ReadBit()
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, false, bb.ReadAligned())
	assert.Equal(t, false, bb.WriteAligned())
	assert.Equal(t, 1, bb.BitsRead())
	assert.Equal(t, 1, bb.BitsWritten())
	assert.Equal(t, []byte{0x01}, nb(bb.Bytes()))

	err = bb.WriteBit(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, bb.ReadAligned())
	assert.Equal(t, false, bb.WriteAligned())
	assert.Equal(t, 1, bb.BitsRead())
	assert.Equal(t, 2, bb.BitsWritten())
	assert.Equal(t, []byte{0x03}, nb(bb.Bytes()))

	err = bb.WriteBit(false)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, bb.ReadAligned())
	assert.Equal(t, false, bb.WriteAligned())
	assert.Equal(t, 1, bb.BitsRead())
	assert.Equal(t, 3, bb.BitsWritten())
	assert.Equal(t, []byte{0x03}, nb(bb.Bytes()))

	bit, err = bb.ReadBit()
	assert.Equal(t, true, bit)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, bb.ReadAligned())
	assert.Equal(t, false, bb.WriteAligned())
	assert.Equal(t, 2, bb.BitsRead())
	assert.Equal(t, 3, bb.BitsWritten())
	assert.Equal(t, []byte{0x03}, nb(bb.Bytes()))

	err = bb.WriteBit(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, bb.ReadAligned())
	assert.Equal(t, false, bb.WriteAligned())
	assert.Equal(t, 2, bb.BitsRead())
	assert.Equal(t, 4, bb.BitsWritten())
	assert.Equal(t, []byte{0x0b}, nb(bb.Bytes()))

	err = bb.WriteBit(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, bb.ReadAligned())
	assert.Equal(t, false, bb.WriteAligned())
	assert.Equal(t, 2, bb.BitsRead())
	assert.Equal(t, 5, bb.BitsWritten())
	assert.Equal(t, []byte{0x1b}, nb(bb.Bytes()))

	err = bb.WriteBit(false)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, bb.ReadAligned())
	assert.Equal(t, false, bb.WriteAligned())
	assert.Equal(t, 2, bb.BitsRead())
	assert.Equal(t, 6, bb.BitsWritten())
	assert.Equal(t, []byte{0x1b}, nb(bb.Bytes()))

	err = bb.WriteBit(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, bb.ReadAligned())
	assert.Equal(t, false, bb.WriteAligned())
	assert.Equal(t, 2, bb.BitsRead())
	assert.Equal(t, 7, bb.BitsWritten())
	assert.Equal(t, []byte{0x5b}, nb(bb.Bytes()))

	err = bb.WriteBit(false)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, bb.ReadAligned())
	assert.Equal(t, true, bb.WriteAligned())
	assert.Equal(t, 2, bb.BitsRead())
	assert.Equal(t, 8, bb.BitsWritten())
	assert.Equal(t, []byte{0x5b}, nb(bb.Bytes()))

	val, cnt, err = bb.ReadBits(6)
	assert.Equal(t, byte(0x16), val)
	assert.Equal(t, 6, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, bb.ReadAligned())
	assert.Equal(t, true, bb.WriteAligned())
	assert.Equal(t, 8, bb.BitsRead())
	assert.Equal(t, 8, bb.BitsWritten())
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
