// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "bytes"

// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
// The zero value for Buffer is an empty buffer ready to use.
type Buffer bytes.Buffer

// Create a new buffer from a byte slice.
func NewBuffer(buf []byte) *Buffer {
	return (*Buffer)(bytes.NewBuffer(buf))
}

// Create a new buffer from a string.
func NewBufferString(str string) *Buffer {
	return (*Buffer)(bytes.NewBufferString(str))
}
