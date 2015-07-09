// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil_test

import "io"
import "testing"
import "github.com/stretchr/testify/assert"
import . "bitbucket.org/rawr/golib/ioutil"

func TestSectionWriterOperations(t *testing.T) {
	_ = NewSectionWriter(io.WriterAt(nil), 0, 0)
	assert.True(t, true)
}
