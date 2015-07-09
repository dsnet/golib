// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

type MultiByteWriter struct{ W []io.ByteWriter }

func NewMultiByteWriter(wrs ...io.ByteWriter) *MultiByteWriter {
	w := make([]io.ByteWriter, len(wrs))
	copy(w, wrs)
	return &MultiByteWriter{w}
}

func (t *MultiByteWriter) WriteByte(val byte) (err error) {
	for _, wr := range t.W {
		if err = wr.WriteByte(val); err != nil {
			return
		}
	}
	return nil
}
