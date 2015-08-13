// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"
import "io/ioutil"
import "strings"
import "testing"
import "github.com/stretchr/testify/assert"

func TestRepeatReader(t *testing.T) {
	var rr io.Reader
	var data []byte

	rr = NewRepeatReaderString("", "", "", 0)
	data, _ = ioutil.ReadAll(rr)
	assert.Equal(t, "", string(data))

	rr = NewRepeatReaderString("a", "", "c", 0)
	data, _ = ioutil.ReadAll(rr)
	assert.Equal(t, "ac", string(data))

	rr = NewRepeatReaderString("", "", "", 5)
	data, _ = ioutil.ReadAll(rr)
	assert.Equal(t, "", string(data))

	rr = NewRepeatReaderString("a", "", "c", 5)
	data, _ = ioutil.ReadAll(rr)
	assert.Equal(t, "ac", string(data))

	rr = NewRepeatReaderString("a", "b", "c", 0)
	data, _ = ioutil.ReadAll(rr)
	assert.Equal(t, "ac", string(data))

	rr = NewRepeatReaderString("a", "b", "c", 5)
	data, _ = ioutil.ReadAll(rr)
	assert.Equal(t, "abbbbbc", string(data))

	rr = NewRepeatReaderString("aa", "bb", "cc", 1)
	data, _ = ioutil.ReadAll(rr)
	assert.Equal(t, "aabbcc", string(data))

	rr = NewRepeatReaderString("aa", "bbb", "cc", 54321)
	data, _ = ioutil.ReadAll(rr)
	assert.Equal(t, "aa"+strings.Repeat("bbb", 54321)+"cc", string(data))
}

func BenchmarkRepeatReader1(b *testing.B) {
	// This benchmarks the worse-case scenario of repeating only a single byte.
	rr := NewRepeatReaderString("", " ", "", int64(b.N))
	b.SetBytes(1)
	b.ResetTimer()
	io.Copy(ioutil.Discard, rr)
}
