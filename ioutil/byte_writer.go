// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

// ByteWriter allows an io.Writer to satisfy the io.ByteWriter interface.
//
// This does not perform any buffering, so performance may be very poor.
// If efficiency is needed, consider using bufio.Writer.
type ByteWriter struct {
	io.Writer
	buf [1]byte
}

// Create a new ByteWriter.
func NewByteWriter(wr io.Writer) *ByteWriter { return &ByteWriter{Writer: wr} }

// Write a single byte.
func (bw *ByteWriter) WriteByte(val byte) (err error) {
	bw.buf[0] = val
	if cnt, err := bw.Write(bw.buf[:]); cnt == 0 {
		if err == nil {
			return io.ErrShortWrite
		}
		return err
	}
	return nil
}
