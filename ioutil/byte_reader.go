// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

// ByteReader allows an io.Reader to satisfy the io.ByteReader interface.
//
// This does not perform any buffering, so performance may be very poor.
// If efficiency is needed, consider using bufio.Reader.
type ByteReader struct {
	io.Reader
	buf [1]byte
}

// Create a new ByteReader.
func NewByteReader(rd io.Reader) *ByteReader { return &ByteReader{Reader: rd} }

// Read a single byte.
func (br *ByteReader) ReadByte() (val byte, err error) {
	if cnt, err := br.Read(br.buf[:]); cnt == 0 {
		if err == nil {
			return 0, io.ErrUnexpectedEOF
		}
		return 0, err
	}
	val = br.buf[0]
	return val, nil
}
