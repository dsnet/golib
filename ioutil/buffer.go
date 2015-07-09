// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package ioutil is a collection of small io related implementations.
package ioutil

import "bytes"

// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
// The zero value for Buffer is an empty buffer ready to use.
type Buffer struct {
	bytes.Buffer
}

// A wrapper around the Buffer provided by package bytes. It is provided in this
// library so that one does not need to import bytes also.
func NewBuffer(buf []byte) *Buffer {
	return &Buffer{*bytes.NewBuffer(buf)}
}

func NewBufferString(str string) *Buffer {
	return &Buffer{*bytes.NewBufferString(str)}
}
