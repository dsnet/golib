// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"
import "testing"
import "github.com/stretchr/testify/assert"

func TestByteReader(t *testing.T) {
	var test = "hello, world!"
	var rd io.Reader
	var brd *ByteReader
	var buf = make([]byte, 1024)
	var val byte
	var cnt int
	var err error

	rd = NewReader([]byte(nil))
	brd = NewByteReader(rd)

	_, err = brd.ReadByte()
	assert.Equal(t, io.EOF, err)

	rd = NewReader([]byte("hello, world!"))
	brd = NewByteReader(rd)

	val, err = brd.ReadByte()
	assert.Equal(t, byte(test[0]), val)
	assert.Nil(t, err)

	cnt, err = brd.Read(buf[:10])
	assert.Equal(t, 10, cnt)
	assert.Equal(t, test[1:11], string(buf[:cnt]))
	assert.Nil(t, err)

	val, err = brd.ReadByte()
	assert.Equal(t, byte(test[11]), val)
	assert.Nil(t, err)

	val, err = brd.ReadByte()
	assert.Equal(t, byte(test[12]), val)
	assert.Nil(t, err)

	_, err = brd.ReadByte()
	assert.Equal(t, io.EOF, err)
}
