// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil_test

import "io"
import "bytes"
import "testing"
import "github.com/stretchr/testify/assert"
import . "bitbucket.org/rawr/golib/ioutil"

func TestLimitedWriterOperations(t *testing.T) {
	var w *bytes.Buffer
	var wr io.Writer
	var cnt int
	var err error
	data := make([]byte, 64)

	// Negative limit
	w = new(bytes.Buffer)
	wr = NewLimitedWriter(w, -1)

	cnt, err = wr.Write(data)
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, io.ErrShortWrite)

	// Zero limit
	w = new(bytes.Buffer)
	wr = NewLimitedWriter(w, 0)

	cnt, err = wr.Write(data)
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, io.ErrShortWrite)

	// Positive limit
	w = new(bytes.Buffer)
	wr = NewLimitedWriter(w, 5)

	cnt, err = wr.Write([]byte("hello, world!"))
	assert.Equal(t, cnt, 5)
	assert.Equal(t, err, io.ErrShortWrite)
	assert.Equal(t, string(w.Bytes()), "hello")

	// Multiple writes
	w = new(bytes.Buffer)
	wr = NewLimitedWriter(w, 12)

	cnt, err = wr.Write([]byte("hello"))
	assert.Equal(t, cnt, 5)
	assert.Equal(t, err, nil)
	assert.Equal(t, string(w.Bytes()), "hello")

	cnt, err = wr.Write([]byte(", world!"))
	assert.Equal(t, cnt, 7)
	assert.Equal(t, err, io.ErrShortWrite)
	assert.Equal(t, string(w.Bytes()), "hello, world")

	cnt, err = wr.Write([]byte("eof"))
	assert.Equal(t, cnt, 0)
	assert.Equal(t, err, io.ErrShortWrite)
	assert.Equal(t, string(w.Bytes()), "hello, world")
}
