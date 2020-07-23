// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsoncs

import (
	"unicode/utf16"
	"unicode/utf8"
)

// lessUTF16 reports whether x is lexicographically less than y according
// to the UTF-16 codepoints of the UTF-8 encoded input strings.
// This implements the ordering specified in RFC 8785, section 3.2.3.
func lessUTF16(x, y string) bool {
	for {
		if len(x) == 0 || len(y) == 0 {
			return len(x) < len(y)
		}

		// ASCII fast-path.
		if x[0] < utf8.RuneSelf || y[0] < utf8.RuneSelf {
			if x[0] != y[0] {
				return x[0] < y[0]
			}
			x, y = x[1:], y[1:]
			continue
		}

		// Decode next pair of runes as UTF-8.
		rx, nx := utf8.DecodeRuneInString(x)
		ry, ny := utf8.DecodeRuneInString(y)
		switch {

		// Both runes encode as either a single or surrogate pair
		// of UTF-16 codepoints.
		case isUTF16Self(rx) == isUTF16Self(ry):
			if rx != ry {
				return rx < ry
			}

		// The x rune is a single UTF-16 codepoint, while
		// the y rune is a surrogate pair of UTF-16 codepoints.
		case isUTF16Self(rx):
			ry, _ := utf16.EncodeRune(ry)
			if rx != ry {
				return rx < ry
			}
			panic("invalid UTF-8") // implies rx is an unpaired surrogate half

		// The y rune is a single UTF-16 codepoint, while
		// the x rune is a surrogate pair of UTF-16 codepoints.
		case isUTF16Self(ry):
			rx, _ := utf16.EncodeRune(rx)
			if rx != ry {
				return rx < ry
			}
			panic("invalid UTF-8") // implies ry is an unpaired surrogate half
		}
		x, y = x[nx:], y[ny:]
	}
}

func isUTF16Self(r rune) bool {
	return ('\u0000' <= r && r <= '\uD7FF') || ('\uE000' <= r && r <= '\uFFFF')
}
