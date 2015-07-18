// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package errs

import "io"
import "testing"
import "github.com/stretchr/testify/assert"

func TestPanic(t *testing.T) {
	assert.NotPanics(t, func() { Panic(nil) })
	assert.NotPanics(t, func() { Panic(error(nil)) })
	assert.Panics(t, func() { Panic(io.EOF) })
	assert.Panics(t, func() { Panic(New("error")) })
}

func TestRecover(t *testing.T) {
	var err error

	// Nil type, should recover
	assert.NotPanics(t, func() {
		defer Recover(&err)
		panic(nil)
	})
	assert.Equal(t, nil, err)

	// Error type, should recover
	assert.NotPanics(t, func() {
		defer Recover(&err)
		panic(io.EOF)
	})
	assert.Equal(t, io.EOF, err)

	// Non error type, should panic
	assert.Panics(t, func() {
		defer Recover(&err)
		panic(5)
	})

	// Runtime error, should panic
	assert.Panics(t, func() {
		defer Recover(&err)
		var s string = "abc"
		_ = s[100]
	})

	// Runtime error, should panic
	assert.Panics(t, func() {
		defer Recover(&err)
		var s interface{} = "abc"
		_ = s.(int)
	})

	// Runtime error, should panic
	assert.Panics(t, func() {
		defer Recover(&err)
		cnt := -1
		_ = make([]byte, cnt)
	})

	// Non-error, should panic
	assert.Panics(t, func() {
		defer Recover(&err)
		panic((*int)(nil))
	})

	// Non-error, should panic
	assert.Panics(t, func() {
		defer Recover(&err)
		panic(([]byte)(nil))
	})
}

func TestNilRecover(t *testing.T) {
	assert.NotPanics(t, func() {
		defer NilRecover()
		panic(nil)
	})

	assert.NotPanics(t, func() {
		defer NilRecover()
		panic(io.EOF)
	})

	assert.NotPanics(t, func() {
		defer NilRecover()
		panic(5)
	})

	assert.NotPanics(t, func() {
		defer NilRecover()
		var s string = "abc"
		_ = s[100]
	})

	assert.NotPanics(t, func() {
		defer NilRecover()
		var s interface{} = "abc"
		_ = s.(int)
	})

	assert.NotPanics(t, func() {
		defer NilRecover()
		cnt := -1
		_ = make([]byte, cnt)
	})

	assert.NotPanics(t, func() {
		defer NilRecover()
		panic((*int)(nil))
	})

	assert.NotPanics(t, func() {
		defer NilRecover()
		panic(([]byte)(nil))
	})
}

func TestConvert(t *testing.T) {
	assert.Equal(t, io.EOF, Convert(io.EOF, io.ErrUnexpectedEOF))
	assert.Equal(t, io.ErrUnexpectedEOF, Convert(io.EOF, io.ErrUnexpectedEOF, io.EOF))
	assert.Equal(t, io.ErrUnexpectedEOF, Convert(io.EOF, io.ErrUnexpectedEOF, io.ErrShortWrite, io.EOF))
}

func TestMatch(t *testing.T) {
	assert.False(t, Match(nil))
	assert.False(t, Match(io.EOF))
	assert.False(t, Match(io.EOF, io.ErrUnexpectedEOF))
	assert.True(t, Match(io.EOF, io.ErrUnexpectedEOF, io.EOF))
}

func TestIgnore(t *testing.T) {
	assert.Equal(t, nil, Ignore(nil))
	assert.Equal(t, io.EOF, Ignore(io.EOF))
	assert.Equal(t, io.EOF, Ignore(io.EOF, io.ErrUnexpectedEOF))
	assert.Equal(t, nil, Ignore(io.EOF, io.ErrUnexpectedEOF, io.EOF))
}

func TestFirst(t *testing.T) {
	assert.Equal(t, nil, First())
	assert.Equal(t, io.EOF, First(io.EOF))
	assert.Equal(t, io.ErrUnexpectedEOF, First(io.ErrUnexpectedEOF))
	assert.Equal(t, io.EOF, First(io.EOF, io.ErrUnexpectedEOF))
}

func TestAssert(t *testing.T) {
	var err error

	err = nil
	func() {
		defer Recover(&err)
		Assert(true, io.EOF)
	}()
	assert.Equal(t, nil, err)

	err = nil
	func() {
		defer Recover(&err)
		Assert(false, io.ErrUnexpectedEOF)
	}()
	assert.Equal(t, io.ErrUnexpectedEOF, err)
}
