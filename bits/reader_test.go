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
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, true, br.ReadAligned())
	assert.Equal(t, 0, br.BytesRead())
	assert.Equal(t, 0, br.BitsRead())

	// Read first byte
	b.WriteByte(0x9d)

	bit, err = br.ReadBit()
	assert.Equal(t, true, bit)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, br.ReadAligned())
	assert.Equal(t, 1, br.BytesRead())
	assert.Equal(t, 1, br.BitsRead())

	bit, err = br.ReadBit()
	assert.Equal(t, false, bit)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, br.ReadAligned())
	assert.Equal(t, 1, br.BytesRead())
	assert.Equal(t, 2, br.BitsRead())

	bit, err = br.ReadBit()
	assert.Equal(t, true, bit)
	assert.Equal(t, nil, err)

	bit, err = br.ReadBit()
	assert.Equal(t, true, bit)
	assert.Equal(t, nil, err)

	bit, err = br.ReadBit()
	assert.Equal(t, true, bit)
	assert.Equal(t, nil, err)

	bit, err = br.ReadBit()
	assert.Equal(t, false, bit)
	assert.Equal(t, nil, err)

	_, err = br.ReadByte()
	assert.Equal(t, ErrAlign, err)
	assert.Equal(t, false, br.ReadAligned())
	assert.Equal(t, 1, br.BytesRead())
	assert.Equal(t, 6, br.BitsRead())

	bit, err = br.ReadBit()
	assert.Equal(t, false, bit)
	assert.Equal(t, nil, err)

	bit, err = br.ReadBit()
	assert.Equal(t, true, bit)
	assert.Equal(t, nil, err)

	assert.Equal(t, true, br.ReadAligned())
	assert.Equal(t, 1, br.BytesRead())
	assert.Equal(t, 8, br.BitsRead())

	_, err = br.ReadByte()
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, true, br.ReadAligned())
	assert.Equal(t, 1, br.BytesRead())
	assert.Equal(t, 8, br.BitsRead())

	// Read second byte
	b.WriteByte(0xa7)

	val, err = br.ReadByte()
	assert.Equal(t, byte(0xa7), val)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, br.ReadAligned())
	assert.Equal(t, 2, br.BytesRead())
	assert.Equal(t, 16, br.BitsRead())

	_, err = br.ReadBit()
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, true, br.ReadAligned())
	assert.Equal(t, 2, br.BytesRead())
	assert.Equal(t, 16, br.BitsRead())

	// Reset
	br.Reset(nil)
	assert.Equal(t, true, br.ReadAligned())
	assert.Equal(t, 0, br.BytesRead())
	assert.Equal(t, 0, br.BitsRead())
}
