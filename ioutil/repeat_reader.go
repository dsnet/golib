// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"
import "bytes"

// The RepeatReader creates a reader stream that emits the following:
//	head body{n} tail
type RepeatReader struct {
	curr, body, tail []byte
	n                int64
}

// Create a new RepeatReader using copies of the provided slices.
func NewRepeatReader(head, body, tail []byte, n int64) *RepeatReader {
	headCpy := append([]byte(nil), head...)
	bodyCpy := append([]byte(nil), body...)
	tailCpy := append([]byte(nil), tail...)
	return newRepeatReader(headCpy, bodyCpy, tailCpy, n)
}

// Create a new RepeatReader using copies of the provided strings.
func NewRepeatReaderString(head, body, tail string, n int64) *RepeatReader {
	headCpy := []byte(head)
	bodyCpy := []byte(body)
	tailCpy := []byte(tail)
	return newRepeatReader(headCpy, bodyCpy, tailCpy, n)
}

func newRepeatReader(head, body, tail []byte, n int64) *RepeatReader {
	const minLength = 4096
	if n == 0 || len(body) == 0 || len(body) >= minLength {
		return &RepeatReader{head, body, tail, n}
	}

	// Repeating a short body can lead to poor performance. Thus, replicate the
	// body to exceed the minLength.
	var ceil = func(n, m int) int { return (n + m - 1) / m }
	nx := int64(ceil(minLength, len(body)))
	if nx > n {
		nx = n
	}
	headBody := append(head, bytes.Repeat(body, int(nx))...)
	repBody := headBody[len(head):]
	headBody = headBody[:len(head)+int(n%nx)*len(body)]
	return &RepeatReader{headBody, repBody, tail, n / nx}
}

// Read from the RepeatReader.
func (rr *RepeatReader) Read(buf []byte) (cnt int, err error) {
	for len(buf) > 0 {
		cpyCnt := copy(buf, rr.curr)
		rr.curr = rr.curr[cpyCnt:]
		buf = buf[cpyCnt:]
		cnt += cpyCnt

		if len(rr.curr) == 0 {
			switch {
			case rr.n > 0:
				rr.curr = rr.body
			case rr.n == 0:
				rr.curr = rr.tail
			case rr.n < 0:
				return cnt, io.EOF
			}
			rr.n--
		}
	}
	return cnt, nil
}
