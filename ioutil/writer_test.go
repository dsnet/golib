// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "testing"

func TestWriter(t *testing.T) {
	// TODO(jtsai): Finish this!
	/*
		var rd *Reader
		var wr *Writer
		var pos int64
		var cnt64 int64
		var cnt int
		var err error
		data := make([]byte, 64)

		// Empty writer
		wr = NewWriter(nil)

		assert.Equal(t, wr.Len(), 0)

		cnt, err = wr.Write(data)
		assert.Equal(t, cnt, 0)
		assert.Equal(t, err, io.ErrShortWrite)

		cnt, err = wr.WriteAt(data, -1)
		assert.Equal(t, cnt, 0)
		assert.NotEqual(t, err, nil)

		cnt, err = wr.WriteAt(data, 1)
		assert.Equal(t, cnt, 0)
		assert.Equal(t, err, io.ErrShortWrite)

		err = wr.WriteByte(127)
		assert.Equal(t, err, io.ErrShortWrite)

		cnt, err = wr.WriteRune(12345)
		assert.Equal(t, cnt, 0)
		assert.Equal(t, err, io.ErrShortWrite)

		rd = NewReader([]byte("Hello, world!"))
		cnt64, err = wr.ReadFrom(rd)
		assert.Equal(t, cnt64, 0)
		assert.Equal(t, err, nil)

		pos, err = wr.Seek(0, os.SEEK_SET)
		assert.Equal(t, pos, 0)
		assert.Equal(t, err, nil)

		pos, err = wr.Seek(0, os.SEEK_CUR)
		assert.Equal(t, pos, 0)
		assert.Equal(t, err, nil)

		pos, err = wr.Seek(0, os.SEEK_END)
		assert.Equal(t, pos, 0)
		assert.Equal(t, err, nil)

		pos, err = wr.Seek(-1, os.SEEK_SET)
		assert.Equal(t, pos, 0)
		assert.NotEqual(t, err, nil)

		pos, err = wr.Seek(-1, os.SEEK_CUR)
		assert.Equal(t, pos, 0)
		assert.NotEqual(t, err, nil)

		pos, err = wr.Seek(-1, os.SEEK_END)
		assert.Equal(t, pos, 0)
		assert.NotEqual(t, err, nil)

		pos, err = wr.Seek(+1, os.SEEK_SET)
		assert.Equal(t, pos, 1)
		assert.Equal(t, err, nil)

		pos, err = wr.Seek(+1, os.SEEK_CUR)
		assert.Equal(t, pos, 2)
		assert.Equal(t, err, nil)

		pos, err = wr.Seek(+1, os.SEEK_END)
		assert.Equal(t, pos, 1)
		assert.Equal(t, err, nil)

		assert.Equal(t, wr.Len(), 0)

		cnt, err = wr.Write(data)
		assert.Equal(t, cnt, 0)
		assert.Equal(t, err, io.ErrShortWrite)

		// Actual writer
		wr = NewWriter(make([]byte, 32))

		assert.Equal(t, wr.Len(), 0)
	*/
}
