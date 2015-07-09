// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil_test

import "io"
import "os"
import "testing"
import "github.com/stretchr/testify/assert"
import . "bitbucket.org/rawr/golib/ioutil"

func TestMemoryReaderOperations(t *testing.T) {
	var rd *MemoryReader
	var cnt int
	var pos, posLo, posHi int64
	var err error
	data := make([]byte, 64)

	// Empty reader and empty ring buffer
	rd = NewMemoryReader(NewReader(nil), nil)

	cnt, err = rd.Read(data)
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, io.EOF)

	pos, err = rd.Seek(0, os.SEEK_CUR)
	assert.Equal(t, pos, 0)
	assert.Equal(t, err, nil)

	pos, err = rd.Seek(1, os.SEEK_CUR)
	assert.Equal(t, pos, 0)
	assert.NotEqual(t, err, nil)

	pos, err = rd.Seek(-1, os.SEEK_CUR)
	assert.Equal(t, pos, 0)
	assert.NotEqual(t, err, nil)

	pos, err = rd.Seek(0, os.SEEK_END)
	assert.Equal(t, pos, 0)
	assert.NotEqual(t, err, nil)

	posLo, pos, posHi = rd.SeekOffsets()
	assert.Equal(t, posLo, 0)
	assert.Equal(t, pos, 0)
	assert.Equal(t, posHi, 0)

	// Actual reader, but empty ring buffer
	rd = NewMemoryReader(NewReader([]byte("Hello, world!")), nil)

	cnt, err = rd.Read(data)
	assert.Equal(t, cnt, 13)
	assert.Equal(t, string(data[:cnt]), "Hello, world!")
	assert.True(t, err == nil || err == io.EOF)

	pos, err = rd.Seek(13, os.SEEK_SET)
	assert.Equal(t, pos, 13)
	assert.Equal(t, err, nil)

	pos, err = rd.Seek(0, os.SEEK_CUR)
	assert.Equal(t, pos, 13)
	assert.Equal(t, err, nil)

	pos, err = rd.Seek(1, os.SEEK_CUR)
	assert.Equal(t, pos, 0)
	assert.NotEqual(t, err, nil)

	pos, err = rd.Seek(-1, os.SEEK_CUR)
	assert.Equal(t, pos, 0)
	assert.NotEqual(t, err, nil)

	posLo, pos, posHi = rd.SeekOffsets()
	assert.Equal(t, posLo, 13)
	assert.Equal(t, pos, 13)
	assert.Equal(t, posHi, 13)

	// Actual reader with ring buffer
	rd = NewMemoryReader(NewReader([]byte("Hello, world!")), make([]byte, 3))

	pos, err = rd.Seek(0, os.SEEK_CUR)
	assert.Equal(t, pos, 0)
	assert.Equal(t, err, nil)

	pos, err = rd.Seek(1, os.SEEK_CUR)
	assert.Equal(t, pos, 0)
	assert.NotEqual(t, err, nil)

	pos, err = rd.Seek(-1, os.SEEK_CUR)
	assert.Equal(t, pos, 0)
	assert.NotEqual(t, err, nil)

	posLo, pos, posHi = rd.SeekOffsets()
	assert.Equal(t, posLo, 0)
	assert.Equal(t, pos, 0)
	assert.Equal(t, posHi, 0)

	cnt, err = rd.Read(data[:5])
	assert.Equal(t, cnt, 5)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(data[:cnt]), "Hello")

	pos, err = rd.Seek(0, os.SEEK_CUR)
	assert.Equal(t, pos, 5)
	assert.Equal(t, err, nil)

	pos, err = rd.Seek(1, os.SEEK_CUR)
	assert.Equal(t, pos, 0)
	assert.NotEqual(t, err, nil)

	pos, err = rd.Seek(-1, os.SEEK_CUR)
	assert.Equal(t, pos, 4)
	assert.Equal(t, err, nil)

	posLo, pos, posHi = rd.SeekOffsets()
	assert.Equal(t, posLo, 2)
	assert.Equal(t, pos, 4)
	assert.Equal(t, posHi, 5)

	cnt, err = rd.Read(data[:3])
	assert.Equal(t, cnt, 3)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(data[:cnt]), "o, ")

	posLo, pos, posHi = rd.SeekOffsets()
	assert.Equal(t, posLo, 4)
	assert.Equal(t, pos, 7)
	assert.Equal(t, posHi, 7)

	pos, err = rd.Seek(-3, os.SEEK_CUR)
	assert.Equal(t, pos, 4)
	assert.Equal(t, err, nil)

	cnt, err = rd.Read(data)
	assert.Equal(t, cnt, 9)
	assert.Equal(t, string(data[:cnt]), "o, world!")
	assert.True(t, err == nil || err == io.EOF)
}
