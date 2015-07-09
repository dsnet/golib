// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "testing"

// Helper test function that converts any empty byte slice to the nil slice so
// that equality checks work out fine.
func nb(buf []byte) []byte {
	if len(buf) == 0 {
		return nil
	}
	return buf
}

func TestInterfaces(_ *testing.T) {
	// These should compile just fine.
	var _ BitsReader = NewReader(nil)
	var _ BitsWriter = NewWriter(nil)
	var _ BitsReader = NewBuffer(nil)
	var _ BitsWriter = NewBuffer(nil)
}
