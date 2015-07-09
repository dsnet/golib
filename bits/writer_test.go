// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"
import "bytes"
import "testing"
import "github.com/stretchr/testify/assert"

type buffer struct {
	bytes.Buffer
	fail bool
}

func (b *buffer) WriteByte(val byte) (err error) {
	if b.fail {
		return io.ErrShortWrite
	}
	return b.Buffer.WriteByte(val)
}

func TestWriter(t *testing.T) {
	b := new(buffer)
	bw := NewWriter(b)

	// Write first byte
	assert.Equal(t, true, bw.WriteAligned())
	assert.Equal(t, 0, bw.BytesWritten())
	assert.Equal(t, 0, bw.BitsWritten())

	assert.Nil(t, bw.WriteBit(true))
	assert.Equal(t, false, bw.WriteAligned())
	assert.Equal(t, 0, bw.BytesWritten())
	assert.Equal(t, 1, bw.BitsWritten())

	assert.Nil(t, bw.WriteBit(false))
	assert.Equal(t, false, bw.WriteAligned())
	assert.Equal(t, 0, bw.BytesWritten())
	assert.Equal(t, 2, bw.BitsWritten())

	assert.Nil(t, bw.WriteBit(true))
	assert.Equal(t, false, bw.WriteAligned())
	assert.Equal(t, 0, bw.BytesWritten())
	assert.Equal(t, 3, bw.BitsWritten())

	assert.Equal(t, ErrAlign, bw.WriteByte(0xff))
	assert.Equal(t, []byte(nil), nb(b.Bytes()))

	assert.Nil(t, bw.WriteBit(true))
	assert.Nil(t, bw.WriteBit(true))
	assert.Nil(t, bw.WriteBit(false))
	assert.Nil(t, bw.WriteBit(false))
	assert.Nil(t, bw.WriteBit(true))
	assert.Equal(t, true, bw.WriteAligned())
	assert.Equal(t, 1, bw.BytesWritten())
	assert.Equal(t, 8, bw.BitsWritten())
	assert.Equal(t, []byte{0x9d}, nb(b.Bytes()))

	// Write second byte
	assert.Nil(t, bw.WriteByte(0xa7))
	assert.Equal(t, 2, bw.BytesWritten())
	assert.Equal(t, 16, bw.BitsWritten())
	assert.Equal(t, []byte{0x9d, 0xa7}, nb(b.Bytes()))

	// Write third byte
	b.fail = true
	assert.Equal(t, io.ErrShortWrite, bw.WriteByte(0xff))
	assert.Equal(t, 2, bw.BytesWritten())
	assert.Equal(t, 16, bw.BitsWritten())
	assert.Equal(t, []byte{0x9d, 0xa7}, nb(b.Bytes()))

	assert.Nil(t, bw.WriteBit(false))
	assert.Nil(t, bw.WriteBit(false))
	assert.Nil(t, bw.WriteBit(true))
	assert.Nil(t, bw.WriteBit(false))
	assert.Nil(t, bw.WriteBit(true))
	assert.Nil(t, bw.WriteBit(true))
	assert.Nil(t, bw.WriteBit(true))

	assert.Equal(t, io.ErrShortWrite, bw.WriteBit(false))
	assert.Equal(t, 2, bw.BytesWritten())
	assert.Equal(t, 23, bw.BitsWritten())
	assert.Equal(t, []byte{0x9d, 0xa7}, nb(b.Bytes()))

	b.fail = false
	assert.Nil(t, bw.WriteBit(false))
	assert.Equal(t, 3, bw.BytesWritten())
	assert.Equal(t, 24, bw.BitsWritten())
	assert.Equal(t, []byte{0x9d, 0xa7, 0x74}, nb(b.Bytes()))

	// Reset
	bw.Reset(nil)
	assert.Equal(t, true, bw.WriteAligned())
	assert.Equal(t, 0, bw.BytesWritten())
	assert.Equal(t, 0, bw.BitsWritten())
}

func BenchmarkWriter(b *testing.B) {
	cnt := 1 << 20 // 1 MiB
	buf := bytes.NewBuffer(nil)
	bw := NewWriter(buf)

	b.SetBytes(int64(cnt))
	b.ResetTimer()

	for ni := 0; ni < b.N; ni++ {
		for bi := 0; bi < cnt; bi++ {
			bw.WriteBit(true)
			bw.WriteBit(false)
			bw.WriteBit(true)
			bw.WriteBit(false)
			bw.WriteBit(true)
			bw.WriteBit(false)
			bw.WriteBit(true)
			bw.WriteBit(false)
		}
	}
}
