// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

func NewTeeByteReader(r io.ByteReader, w io.ByteWriter) io.ByteReader {
	return &TeeByteReader{r, w}
}

type TeeByteReader struct {
	R io.ByteReader
	W io.ByteWriter
}

func (t *TeeByteReader) ReadByte() (val byte, err error) {
	val, err = t.R.ReadByte()
	if err == nil {
		err = t.W.WriteByte(val)
	}
	return
}
