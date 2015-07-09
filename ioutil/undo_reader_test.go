// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil_test

import "io"
import "testing"
import "github.com/stretchr/testify/assert"
import . "bitbucket.org/rawr/golib/ioutil"

func TestUndoReaderOperations(t *testing.T) {
	var rd *UndoReader
	var cnt int
	var err error
	data := make([]byte, 64)

	// Test with empty reader and stack
	rd = NewUndoReader(NewReader(nil), nil)

	cnt, err = rd.Read(nil)
	assert.Equal(t, cnt, 0)

	cnt, err = rd.Read(data)
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, io.EOF)

	rd.UndoRead(nil)

	cnt, err = rd.Read(data)
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, io.EOF)

	rd.UndoRead([]byte("angry penguins"))

	cnt, err = rd.Read(data[:5])
	assert.Equal(t, cnt, 5)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(data[:cnt]), "angry")

	rd.UndoRead([]byte("happy"))

	cnt, err = rd.Read(data)
	assert.Equal(t, cnt, 14)
	assert.Equal(t, string(data[:cnt]), "happy penguins")
	assert.True(t, err == nil || err == io.EOF)

	cnt, err = rd.Read(data)
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, io.EOF)

	// Test with actual reader and pre-allocated stack
	r := NewReader([]byte("doggies hate igloos"))
	rd = NewUndoReader(r, make([]byte, 17))

	cnt, err = rd.Read(data[:12])
	assert.Equal(t, cnt, 12)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(data[:cnt]), "doggies hate")

	rd.UndoRead([]byte("penguins love"))
	rd.UndoRead([]byte("happy "))

	cnt, err = rd.Read(data)
	assert.Equal(t, cnt, 26)
	assert.Equal(t, string(data[:cnt]), "happy penguins love igloos")
	assert.True(t, err == nil || err == io.EOF)

	cnt, err = rd.Read(data)
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, io.EOF)
}
