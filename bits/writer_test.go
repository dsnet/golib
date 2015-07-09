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
	assert.Equal(t, bw.ByteAligned(), true)
	assert.Equal(t, bw.BytesWritten(), 0)
	assert.Equal(t, bw.BitsWritten(), 0)

	assert.Nil(t, bw.WriteBit(true))
	assert.Equal(t, bw.ByteAligned(), false)
	assert.Equal(t, bw.BytesWritten(), 0)
	assert.Equal(t, bw.BitsWritten(), 1)

	assert.Nil(t, bw.WriteBit(false))
	assert.Equal(t, bw.ByteAligned(), false)
	assert.Equal(t, bw.BytesWritten(), 0)
	assert.Equal(t, bw.BitsWritten(), 2)

	assert.Nil(t, bw.WriteBit(true))
	assert.Equal(t, bw.ByteAligned(), false)
	assert.Equal(t, bw.BytesWritten(), 0)
	assert.Equal(t, bw.BitsWritten(), 3)

	assert.Equal(t, bw.WriteByte(0xff), ErrAlign)
	assert.Equal(t, b.Bytes(), []byte(nil))

	assert.Nil(t, bw.WriteBit(true))
	assert.Nil(t, bw.WriteBit(true))
	assert.Nil(t, bw.WriteBit(false))
	assert.Nil(t, bw.WriteBit(false))
	assert.Nil(t, bw.WriteBit(true))
	assert.Equal(t, bw.ByteAligned(), true)
	assert.Equal(t, bw.BytesWritten(), 1)
	assert.Equal(t, bw.BitsWritten(), 8)
	assert.Equal(t, b.Bytes(), []byte{0x9d})

	// Write second byte
	assert.Nil(t, bw.WriteByte(0xa7))
	assert.Equal(t, bw.BytesWritten(), 2)
	assert.Equal(t, bw.BitsWritten(), 16)
	assert.Equal(t, b.Bytes(), []byte{0x9d, 0xa7})

	// Write third byte
	b.fail = true
	assert.Equal(t, bw.WriteByte(0xff), io.ErrShortWrite)
	assert.Equal(t, bw.BytesWritten(), 2)
	assert.Equal(t, bw.BitsWritten(), 16)
	assert.Equal(t, b.Bytes(), []byte{0x9d, 0xa7})

	assert.Nil(t, bw.WriteBit(false))
	assert.Nil(t, bw.WriteBit(false))
	assert.Nil(t, bw.WriteBit(true))
	assert.Nil(t, bw.WriteBit(false))
	assert.Nil(t, bw.WriteBit(true))
	assert.Nil(t, bw.WriteBit(true))
	assert.Nil(t, bw.WriteBit(true))

	assert.Equal(t, bw.WriteBit(false), io.ErrShortWrite)
	assert.Equal(t, bw.BytesWritten(), 2)
	assert.Equal(t, bw.BitsWritten(), 23)
	assert.Equal(t, b.Bytes(), []byte{0x9d, 0xa7})

	b.fail = false
	assert.Nil(t, bw.WriteBit(false))
	assert.Equal(t, bw.BytesWritten(), 3)
	assert.Equal(t, bw.BitsWritten(), 24)
	assert.Equal(t, b.Bytes(), []byte{0x9d, 0xa7, 0x74})

	// Reset
	bw.Reset(nil)
	assert.Equal(t, bw.ByteAligned(), true)
	assert.Equal(t, bw.BytesWritten(), 0)
	assert.Equal(t, bw.BitsWritten(), 0)
}
