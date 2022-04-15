// Copyright 2022, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package try emulates aspects of the ill-fated "try" proposal using generics.
// See https://golang.org/issue/32437 for inspiration.
//
// Example usage:
//
//	func Fizz(...) (..., err error) {
//		defer try.Catch(&err, func() {
//			if err == io.EOF {
//				err = io.ErrUnexpectedEOF
//			}
//		})
//		... := try.A2(Buzz(...))
//		return ..., nil
//	}
//
// This package is not intended for production critical code as quick and easy
// error handling can occlude critical error handling logic.
// Rather, it is intended for short Go programs and unit tests where
// development speed is a greater priority than reliability.
//
//
// Code before try:
//
//	func (a *MixedArray) UnmarshalNext(uo json.UnmarshalOptions, d *json.Decoder) error {
//		switch t, err := d.ReadToken(); {
//		case err != nil:
//			return err
//		case t.Kind() != '[':
//			return fmt.Errorf("got %v, expecting array start", t.Kind())
//		}
//
//		if err := uo.UnmarshalNext(d, &a.Scalar); err != nil {
//			return err
//		}
//		if err := uo.UnmarshalNext(d, &a.Slice); err != nil {
//			return err
//		}
//		if err := uo.UnmarshalNext(d, &a.Map); err != nil {
//			return err
//		}
//
//		switch t, err := d.ReadToken(); {
//		case err != nil:
//			return err
//		case t.Kind() != ']':
//			return fmt.Errorf("got %v, expecting array start", t.Kind())
//		}
//		return nil
//	}
//
// Code after try:
//
//	func (a *MixedArray) UnmarshalNext(uo json.UnmarshalOptions, d *json.Decoder) (err error) {
//		defer try.Catch(&err)
//		if t := try.A1(d.ReadToken()); t.Kind() != '[' {
//			return fmt.Errorf("found %v, expecting array start", t.Kind())
//		}
//		try.A0(uo.UnmarshalNext(d, &a.Scalar))
//		try.A0(uo.UnmarshalNext(d, &a.Slice))
//		try.A0(uo.UnmarshalNext(d, &a.Map))
//		if t := try.A1(d.ReadToken()); t.Kind() != ']' {
//			return fmt.Errorf("found %v, expecting array start", t.Kind())
//		}
//		return nil
//	}
//
package try

// wrapError wraps an error to ensure that we only recover from errors
// panicked by this package.
type wrapError struct{ error }

// Unwrap primarily exists for testing purposes.
func (e wrapError) Unwrap() error { return e.error }

// Catch catches a previously panicked error and stores it into err.
// If it successfully catches an error, it calls any provided handlers.
func Catch(err *error, handlers ...func()) {
	switch ex := recover().(type) {
	case nil:
		return
	case wrapError:
		*err = ex.error
		for _, handler := range handlers {
			handler()
		}
	default:
		panic(ex)
	}
}

// A0 panics if err is non-nil.
func A0(err error) {
	if err != nil {
		panic(wrapError{err})
	}
	return
}

// A1 panics if err is non-nil,
// otherwise it returns v1 as is.
func A1[T1 any](v1 T1, err error) T1 {
	A0(err)
	return v1
}

// A2 panics if err is non-nil,
// otherwise it returns v1 and v2 as is.
func A2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	A0(err)
	return v1, v2
}

// A3 panics if err is non-nil,
// otherwise it returns v1, v2, and v3 as is.
func A3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	A0(err)
	return v1, v2, v3
}

// A4 panics if err is non-nil,
// otherwise it returns v1, v2, v3, and v4 as is.
func A4[T1, T2, T3, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) (T1, T2, T3, T4) {
	A0(err)
	return v1, v2, v3, v4
}
