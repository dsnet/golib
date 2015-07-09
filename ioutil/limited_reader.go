// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

// A LimitedReader reads from R but limits the amount of data returned to just N
// bytes. Each call to Read updates N to reflect the new amount remaining.
type LimitedReader io.LimitedReader

// A wrapper around the LimitedReader provided by package io. It is provided in
// this library so that one does not need to import io also. It is the
// functional complement of LimitedWriter.
func NewLimitedReader(rd io.Reader, cnt int64) *LimitedReader {
	return (*LimitedReader)(&io.LimitedReader{R: rd, N: cnt})
}
