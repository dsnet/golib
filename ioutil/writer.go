// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"
import "os"
import "errors"
import "unicode/utf8"

type Writer struct {
	buf []byte
	idx int64
}

func NewWriter(data []byte) *Writer {
	return &Writer{data, 0}
}

func (w *Writer) Len() int {
	if w.idx > int64(len(w.buf)) {
		return len(w.buf)
	} else {
		return int(w.idx)
	}
}

func (w *Writer) Write(data []byte) (cnt int, err error) {
	cnt, err = w.WriteAt(data, w.idx)
	w.idx += int64(cnt)
	return
}

func (w *Writer) WriteAt(data []byte, off int64) (cnt int, err error) {
	if off < 0 {
		return 0, errors.New("ioutil.Writer.WriteAt: invalid argument")
	}
	if off > int64(len(w.buf)) {
		off = int64(len(w.buf))
	}
	cnt = copy(w.buf[off:], data)
	if cnt < len(data) {
		err = io.ErrShortWrite
	}
	return
}

func (w *Writer) WriteByte(b byte) error {
	if w.idx >= int64(len(w.buf)) {
		return io.ErrShortWrite
	}
	w.buf[w.idx] = b
	w.idx++
	return nil
}

func (w *Writer) WriteRune(r rune) (cnt int, err error) {
	cnt = utf8.RuneLen(r)
	if cnt == -1 {
		r = utf8.RuneError
		cnt = utf8.RuneLen(r)
	}
	if availCnt := int64(len(w.buf)) - w.idx; availCnt < int64(cnt) {
		return 0, io.ErrShortWrite
	}
	cnt = utf8.EncodeRune(w.buf[w.idx:], r)
	w.idx += int64(cnt)
	return cnt, nil
}

func (w *Writer) Seek(offset int64, whence int) (pos int64, err error) {
	switch whence {
	case os.SEEK_SET:
		pos = offset
	case os.SEEK_CUR:
		pos = offset + w.idx
	case os.SEEK_END:
		pos = offset + int64(len(w.buf))
	default:
		return 0, errors.New("ioutil.Writer.Seek: invalid whence")
	}
	if pos < 0 {
		return 0, errors.New("ioutil.Writer.Seek: invalid offset")
	}
	w.idx = pos
	return pos, nil
}

func (w *Writer) ReadFrom(rd io.Reader) (cnt int64, err error) {
	cnt, err = io.Copy(w, rd)
	if err == io.ErrShortWrite {
		err = nil
	}
	return
}
