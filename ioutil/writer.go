// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

type Writer struct {
}

func NewWriter(data []byte) *Writer {
	return &Writer{}
}

func (w *Writer) Len() int {
	return 0
}

func (w *Writer) Write(data []byte) (cnt int, err error) {
	return 0, nil
}

func (w *Writer) WriteAt(data []byte, off int64) (cnt int, err error) {
	return 0, nil
}

func (w *Writer) WriteByte(b byte) error {
	return nil
}

func (w *Writer) WriteRune(r rune) (cnt int, err error) {
	return 0, nil
}

func (w *Writer) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (w *Writer) ReadFrom(rd io.Reader) (cnt int64, err error) {
	return 0, nil
}
