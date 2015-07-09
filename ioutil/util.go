// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"
import "math"
import "bitbucket.org/rawr/golib/errs"

// Determine the size of a ReaderAt using a binary search. Given that file
// offsets are no larger than int64, there is an upper limit of 64 iterations
// before the EOF is found.
func ReaderAtSize(rd io.ReaderAt) (pos int64, err error) {
	errs.Recover(&err)

	// Function to check if the given position is at EOF
	buf := make([]byte, 2)
	checkEOF := func(pos int64) int {
		if pos > 0 {
			cnt, err := rd.ReadAt(buf[:2], pos-1)
			errs.Panic(errs.Ignore(err, io.EOF))
			return 1 - cnt // RetVal[Cnt] = {0: +1, 1: 0, 2: -1}
		} else { // Special case where position is zero
			cnt, err := rd.ReadAt(buf[:1], pos-0)
			errs.Panic(errs.Ignore(err, io.EOF))
			return 0 - cnt // RetVal[Cnt] = {0: 0, 1: -1}
		}
	}

	// Obtain the size via binary search O(log n) => 64 iterations
	posMin, posMax := int64(0), int64(math.MaxInt64)
	for posMax >= posMin {
		pos = (posMax + posMin) / 2
		switch checkEOF(pos) {
		case -1: // Below EOF
			posMin = pos + 1
		case 0: // At EOF
			return pos, nil
		case +1: // Above EOF
			posMax = pos - 1
		}
	}
	panic(errs.New("EOF is in a transient state"))
}
