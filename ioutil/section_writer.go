// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"
import "os"
import "errors"

// SectionWriter implements Write, Seek, and WriteAt on a section of an
// underlying WriterAt.
type SectionWriter struct {
	wr    io.WriterAt
	off   int64
	offLo int64
	offHi int64
}

// NewSectionWriter returns a SectionWriter that writes to wr starting at offset
// off and stops with ErrShortWrite after cnt bytes.
func NewSectionWriter(wr io.WriterAt, off int64, cnt int64) *SectionWriter {
	return &SectionWriter{wr, off, off, off + cnt}
}

func (s *SectionWriter) Write(data []byte) (cnt int, err error) {
	cnt, err = s.WriteAt(data, s.off-s.offLo)
	s.off += int64(cnt)
	return
}

func (s *SectionWriter) WriteAt(data []byte, off int64) (cnt int, err error) {
	off += s.offLo
	if off < s.offLo {
		return 0, errors.New("ioutil.SectionWriter.WriteAt: invalid argument")
	}
	if off > s.offHi {
		off = s.offHi
	}
	cnt = int(s.offHi - off)
	cnt, err = s.wr.WriteAt(data[:cnt], off)
	if err == nil && cnt < len(data) {
		err = io.ErrShortWrite
	}
	return
}

func (s *SectionWriter) Seek(offset int64, whence int) (pos int64, err error) {
	switch whence {
	case os.SEEK_SET:
		pos = offset + s.offLo
	case os.SEEK_CUR:
		pos = offset + s.off
	case os.SEEK_END:
		pos = offset + s.offHi
	default:
		return 0, errors.New("ioutil.SectionWriter.Seek: invalid whence")
	}
	if pos < s.offLo {
		return 0, errors.New("ioutil.SectionWriter.Seek: invalid offset")
	}
	s.off = pos
	return pos - s.offLo, nil
}

func (s *SectionWriter) Size() int64 {
	return s.offHi - s.offLo
}
