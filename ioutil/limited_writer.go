// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

type LimitedWriter struct {
	W io.Writer
	N int64
}

func NewLimitedWriter(wr io.Writer, cnt int64) io.Writer {
	return &LimitedWriter{wr, cnt}
}

func (l *LimitedWriter) Write(data []byte) (cnt int, err error) {
	inCnt := len(data)
	if int64(inCnt) > l.N {
		inCnt = l.N
	}
	if l.N < 0 {
		inCnt = 0
	}
	cnt, err = l.W.Write(data[:inCnt])
	if err == nil && cnt < len(data) {
		err = io.ErrShortWrite
	}
	l.N -= cnt
	return
}
