// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "bytes"

type Reader struct {
	*bytes.Reader
}

func NewReader(data []byte) *Reader {
	return &Reader{bytes.NewReader(data)}
}
