// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsoncs

import (
	"bytes"
)

// Valid reports whether the JSON input is in canonical form.
// Invalid JSON input is reported as false.
func Valid(b []byte) bool {
	b, ok := validValue(b)
	return ok && len(b) == 0
}

// validValue reports whether the next JSON value is in its canonical form.
// It consume the leading value, and returns the remaining bytes.
func validValue(b []byte) ([]byte, bool) {
	switch {
	case len(b) > 0 && b[0] == '{':
		return validObject(b)
	case len(b) > 0 && b[0] == '[':
		return validArray(b)
	case len(b) > 0 && b[0] == '"':
		return validString(b)
	case len(b) > 0 && (b[0] == '-' || ('0' <= b[0] && b[0] <= '9')):
		return validNumber(b)
	case bytes.HasPrefix(b, nullLiteral):
		return b[len(nullLiteral):], true
	case bytes.HasPrefix(b, trueLiteral):
		return b[len(trueLiteral):], true
	case bytes.HasPrefix(b, falseLiteral):
		return b[len(falseLiteral):], true
	default:
		return b, false
	}
}

// validObject reports whether the next JSON object is in its canonical form
// per RFC 8785, section 3.2.3 regarding object name ordering.
// It consume the leading value, and returns the remaining bytes.
func validObject(b []byte) ([]byte, bool) {
	if len(b) == 0 || b[0] != '{' {
		return b, false
	}
	b = b[1:]

	var init, ok bool
	var prevKey string
	for {
		if len(b) > 0 && b[0] == '}' {
			return b[1:], true
		}

		if init {
			if len(b) == 0 || b[0] != ',' {
				return b, false
			}
			b = b[1:]
		}

		currKey, _, _ := decodeString(b)
		b, ok = validString(b)
		if !ok {
			return b, ok
		}
		if init && !lessUTF16(prevKey, currKey) {
			return b, ok
		}
		prevKey = currKey

		if len(b) == 0 || b[0] != ':' {
			return b, false
		}
		b = b[1:]

		b, ok = validValue(b)
		if !ok {
			return b, ok
		}

		init = true
	}
}

// validArray reports whether the next JSON array is in its canonical form.
// It consume the leading value, and returns the remaining bytes.
func validArray(b []byte) ([]byte, bool) {
	if len(b) == 0 || b[0] != '[' {
		return b, false
	}
	b = b[1:]

	var init, ok bool
	for {
		if len(b) > 0 && b[0] == ']' {
			return b[1:], true
		}

		if init {
			if len(b) == 0 || b[0] != ',' {
				return b, false
			}
			b = b[1:]
		}

		b, ok = validValue(b)
		if !ok {
			return b, ok
		}

		init = true
	}
}

// validString reports whether the next JSON string is in its canonical form
// per RFC 8785, section 3.2.2.2.
// It consume the leading value, and returns the remaining bytes.
func validString(b []byte) ([]byte, bool) {
	if len(b) == 0 || b[0] != '"' {
		return b, false
	}

	// Fast-path optimization for unescaped ASCII.
	for b := b[1:]; len(b) > 0; b = b[1:] {
		if b[0] == '"' {
			return b[1:], true
		}
		if !(0x20 <= b[0] && b[0] < 0x80 && b[0] != '"' && b[0] != '\\') {
			break
		}
	}

	s, b2, err := decodeString(b)
	got := b[:len(b)-len(b2)]
	want, _ := formatString(nil, s)
	return b2, bytes.Equal(got, want) && err == nil
}

// validNumber reports whether the next JSON number is in its canonical form
// per RFC 8785, section 3.2.2.3.
// It consume the leading value, and returns the remaining bytes.
func validNumber(b []byte) ([]byte, bool) {
	if len(b) == 0 || !(b[0] == '-' || ('0' <= b[0] && b[0] <= '9')) {
		return b, false
	}

	// Fast-path optimization for integers.
	// Integer values in the range of ±2⁵³ are represented in decimal,
	// which is encoded using up to 16 digits (excluding the sign).
	{
		b := b
		var neg bool
		if len(b) > 0 && b[0] == '-' {
			b = b[1:]
			neg = true
		}
		switch {
		case len(b) == 0:
			break
		case b[0] == '0':
			b = b[1:]
			if neg {
				break // -0 is not permitted
			}
			if len(b) > 0 && (b[0] == '.' || b[0] == 'e' || b[0] == 'E') {
				break // number is not yet terminated
			}
			return b, true
		case '1' <= b[0] && b[0] <= '9':
			var n int
			b = b[1:]
			n++
			for len(b) > 0 && ('0' <= b[0] && b[0] <= '9') {
				b = b[1:]
				n++
			}
			if n >= 16 {
				break // possibly exceeds ±2⁵³
			}
			if len(b) > 0 && (b[0] == '.' || b[0] == 'e' || b[0] == 'E') {
				break // number is not yet terminated
			}
			return b, true
		}
	}

	f, b2, err := decodeNumber(b)
	got := b[:len(b)-len(b2)]
	want, _ := formatNumber(nil, f)
	return b2, bytes.Equal(got, want) && err == nil
}
