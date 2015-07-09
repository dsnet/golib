// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"

type SectionReader io.SectionReader

func NewSectionReader(rd io.ReaderAt, off int64, cnt int64) *SectionReader {
	return (*SectionReader)(io.NewSectionReader(rd, off, cnt))
}
