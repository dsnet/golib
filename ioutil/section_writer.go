// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

type SectionWriter struct {
}

func NewSectionWriter(rd io.ReaderAt, off int64, cnt int64) *SectionWriter {
	return nil
}

func (s *SectionWriter) Write(data []byte) (cnt int, err error) {
	return 0, nil
}

func (s *SectionWriter) WriteAt(data []byte, off int64) (cnt int, err error) {
	return 0, nil
}

func (s *SectionWriter) Seek(offset int64, whence int) (pos int64, err error) {
	return 0, nil
}

func (s *SectionWriter) Size() int64 {
	return 0
}
