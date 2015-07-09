// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "testing"
import "github.com/stretchr/testify/assert"

// Helper test function that converts any empty byte slice to the nil slice so
// that equality checks work out fine.
func nb(buf []byte) []byte {
	if len(buf) == 0 {
		return nil
	}
	return buf
}

func TestInterfaces(t *testing.T) {
	var ok bool

	_, ok = interface{}(NewReader(nil)).(BitsReader)
	assert.True(t, ok)
	_, ok = interface{}(NewWriter(nil)).(BitsWriter)
	assert.True(t, ok)
	_, ok = interface{}(NewBuffer(nil)).(BitsReader)
	assert.True(t, ok)
	_, ok = interface{}(NewBuffer(nil)).(BitsWriter)
	assert.True(t, ok)
}
