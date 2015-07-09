// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "bytes"

// A Reader implements the io.Reader, io.ReaderAt, io.WriterTo, io.Seeker,
// io.ByteScanner, and io.RuneScanner interfaces by reading from a byte slice.
type Reader struct {
	*bytes.Reader
}

// A wrapper around the Reader provided by package io. It is provided in this
// library so that one does not need to import io also. It is the functional
// complement of Writer.
func NewReader(data []byte) Reader {
	return Reader{bytes.NewReader(data)}
}
