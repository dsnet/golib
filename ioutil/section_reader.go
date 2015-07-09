// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

// SectionReader implements Read, Seek, and ReadAt on a section of an
// underlying ReaderAt.
type SectionReader io.SectionReader

// A wrapper around the SectionReader provided by package io. It is provided in
// this library so that one does not need to import io also. It is the
// functional complement of SectionWriter.
func NewSectionReader(rd io.ReaderAt, off int64, cnt int64) *SectionReader {
	return (*SectionReader)(io.NewSectionReader(rd, off, cnt))
}
