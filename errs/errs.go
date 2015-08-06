// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package errs implements helpers functions to deal with errors.
package errs

import "runtime"
import "errors"

// Create a new error.
func New(str string) error {
	return errors.New(str)
}

// Panic if the error is not nil.
func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

// Recovers from any panics and stores any errors to the given pointer. If the
// source of the panic was a Runtime error or not an error at all, then the
// panic continues.
func Recover(err *error) {
	switch ex := recover().(type) {
	case nil:
		// Do nothing
	case runtime.Error:
		panic(ex)
	case error:
		*err = ex
	default:
		panic(ex)
	}
}

// TODO(jtsai): Remove the declaration of Recover above and replace it with the
// one below once Go's escape analysis improves. Otherwise, this function incurs
// some allocation penalty for every invocation.
//
// See: https://golang.org/issues/12006
/*
// Recovers from any panics and stores any errors to all of the input pointers.
// If the source of the panic was a Runtime error or not an error at all, then
// the panic continues.
func Recover(errs ...*error) {
	switch ex := recover().(type) {
	case nil:
		// Do nothing
	case runtime.Error:
		panic(ex)
	case error:
		for _, err := range errs {
			*err = ex
		}
	default:
		panic(ex)
	}
}
*/

// Recovers from any panics and ignores them.
// This is dangerous and should be used sparingly.
func NilRecover() {
	_ = recover()
}

// Convert errors from one type to another.
func Convert(err error, errNew error, errChks ...error) error {
	for _, errChk := range errChks {
		if err == errChk {
			return errNew
		}
	}
	return err
}

// Check if the provided error matches any of the others.
func Match(err error, errChks ...error) bool {
	return err != nil && Convert(err, nil, errChks...) == nil
}

// Ignore certain types of errors.
func Ignore(err error, errChks ...error) error {
	return Convert(err, nil, errChks...)
}

// Return the first non-nil error.
func First(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// Assert that an condition is true.
func Assert(cond bool, err error) {
	if !cond {
		panic(err)
	}
}
