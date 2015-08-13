// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"
import "bytes"

// Wrap some of the data structures provided by package io and bytes so that
// importers of ioutil do not need to import those packages as well.

// A wrap of bytes.Buffer.
type Buffer struct{ bytes.Buffer }

// A wrap of bytes.NewBuffer.
func NewBuffer(buf []byte) *Buffer {
	return &Buffer{*bytes.NewBuffer(buf)}
}

// A wrap of bytes.NewBufferString.
func NewBufferString(str string) *Buffer {
	return &Buffer{*bytes.NewBufferString(str)}
}

// A wrap of bytes.Reader.
type Reader struct{ bytes.Reader }

// A wrap of bytes.NewReader.
func NewReader(data []byte) *Reader {
	return &Reader{*bytes.NewReader(data)}
}

// A wrap of io.LimitedReader.
type LimitedReader struct{ io.LimitedReader }

// A wrap of io.LimitedReader.
func NewLimitedReader(rd io.Reader, cnt int64) *LimitedReader {
	return &LimitedReader{io.LimitedReader{R: rd, N: cnt}}
}

// A wrap of io.SectionReader.
type SectionReader struct{ io.SectionReader }

// A wrap of io.NewSectionReader.
func NewSectionReader(rd io.ReaderAt, off int64, cnt int64) *SectionReader {
	return &SectionReader{*io.NewSectionReader(rd, off, cnt)}
}

// A wrap of io.TeeReader.
func NewTeeReader(rd io.Reader, wr io.Writer) io.Reader {
	return io.TeeReader(rd, wr)
}

// A wrap of io.MultiReader.
func NewMultiReader(rd ...io.Reader) io.Reader {
	return io.MultiReader(rd...)
}

// A wrap of io.MultiWriter.
func NewMultiWriter(wr ...io.Writer) io.Writer {
	return io.MultiWriter(wr...)
}
