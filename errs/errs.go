// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package errs implements helpers functions to deal with errors.
package errs

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

// Recovers from any panics and stores errors to the given pointer. If the
// source of the panic was not an error, then the panic continues.
func Recover(err *error) {
	if ex := recover(); ex != nil {
		if _err, ok := ex.(error); ok {
			(*err) = _err
		} else {
			panic(ex)
		}
	}
}

// Recovers from any panics and ignores them.
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
