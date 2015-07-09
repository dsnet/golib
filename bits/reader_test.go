// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"
import "bytes"
import "testing"
import "github.com/stretchr/testify/assert"

func TestReader(t *testing.T) {
	var bit bool
	var val byte
	var err error

	b := bytes.NewBuffer(nil)
	br := NewReader(b)

	// Read zeroth byte
	_, err = br.ReadBit()
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, br.ByteAligned(), true)
	assert.Equal(t, br.BytesRead(), 0)
	assert.Equal(t, br.BitsRead(), 0)

	// Read first byte
	b.WriteByte(0x9d)

	bit, err = br.ReadBit()
	assert.Equal(t, bit, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, br.ByteAligned(), false)
	assert.Equal(t, br.BytesRead(), 1)
	assert.Equal(t, br.BitsRead(), 1)

	bit, err = br.ReadBit()
	assert.Equal(t, bit, false)
	assert.Equal(t, err, nil)
	assert.Equal(t, br.ByteAligned(), false)
	assert.Equal(t, br.BytesRead(), 1)
	assert.Equal(t, br.BitsRead(), 2)

	bit, err = br.ReadBit()
	assert.Equal(t, bit, true)
	assert.Equal(t, err, nil)

	bit, err = br.ReadBit()
	assert.Equal(t, bit, true)
	assert.Equal(t, err, nil)

	bit, err = br.ReadBit()
	assert.Equal(t, bit, true)
	assert.Equal(t, err, nil)

	bit, err = br.ReadBit()
	assert.Equal(t, bit, false)
	assert.Equal(t, err, nil)

	_, err = br.ReadByte()
	assert.Equal(t, err, ErrAlign)
	assert.Equal(t, br.ByteAligned(), false)
	assert.Equal(t, br.BytesRead(), 1)
	assert.Equal(t, br.BitsRead(), 6)

	bit, err = br.ReadBit()
	assert.Equal(t, bit, false)
	assert.Equal(t, err, nil)

	bit, err = br.ReadBit()
	assert.Equal(t, bit, true)
	assert.Equal(t, err, nil)

	assert.Equal(t, br.ByteAligned(), true)
	assert.Equal(t, br.BytesRead(), 1)
	assert.Equal(t, br.BitsRead(), 8)

	_, err = br.ReadByte()
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, br.ByteAligned(), true)
	assert.Equal(t, br.BytesRead(), 1)
	assert.Equal(t, br.BitsRead(), 8)

	// Read second byte
	b.WriteByte(0xa7)

	val, err = br.ReadByte()
	assert.Equal(t, val, byte(0xa7))
	assert.Equal(t, err, nil)
	assert.Equal(t, br.ByteAligned(), true)
	assert.Equal(t, br.BytesRead(), 2)
	assert.Equal(t, br.BitsRead(), 16)

	_, err = br.ReadBit()
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, br.ByteAligned(), true)
	assert.Equal(t, br.BytesRead(), 2)
	assert.Equal(t, br.BitsRead(), 16)

	// Reset
	br.Reset(nil)
	assert.Equal(t, br.ByteAligned(), true)
	assert.Equal(t, br.BytesRead(), 0)
	assert.Equal(t, br.BitsRead(), 0)
}
