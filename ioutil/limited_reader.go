// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

type LimitedReader io.LimitedReader

func NewLimitedReader(rd io.Reader, cnt int64) io.Reader {
	return io.LimitReader(rd, cnt)
}
