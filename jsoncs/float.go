// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsoncs

import (
	"bytes"
	"math"
	"strconv"
)

// BUG(https://golang.org/issue/29491): On Go 1.12 and earlier, float formatting
// produces incorrect results in some rare cases due to a rounding error.

// formatES6 appends a canonically formatted float according to ES6,
// formally defined in ECMA-262, 6th edition, section 7.1.12.1.
// This implements the formatting specified in RFC 8785, section 3.2.2.3,
// except that this handles NaNs and Infinity.
func formatES6(b []byte, m float64) []byte {
	switch {
	// 1. If m is NaN, return the String "NaN".
	case math.IsNaN(m):
		return append(b, "NaN"...)
	// 2. If m is +0 or -0, return the String "0".
	case m == 0:
		return append(b, "0"...)
	// 3. If m is less than zero, return the String concatenation of the String "-" and ToString(-m).
	case m < 0:
		return formatES6(append(b, '-'), -m)
	// 4. If m is +∞, return the String "Infinity".
	case math.IsInf(m, +1):
		return append(b, "Infinity"...)
	}

	// 5. Otherwise, let n, k, and s be integers such that k ≥ 1, 10^(k-1) ≤ s < 10^k,
	// the Number value for s ⨯ 10^(n-k) is m, and k is as small as possible.
	// If there are multiple possibilities for s, choose the value of s for which s ⨯ 10^(n-k)
	// is closest in value to m. If there are two such possible values of s, choose the one that is even.
	// Note that k is the number of digits in the decimal representation of s
	// and that s is not divisible by 10.
	var n, k int64
	var s []byte // decimal representation of s
	{
		// Unfortunately strconv.FormatFloat does not directly expose n, k, s,
		// nor is the output defined as stable in any way.
		// However, it's implementation is guaranteed to produce precise n, k, s values.
		// Format a float with the 'e' format and derive n, k, s from the output.
		var arr [32]byte
		b := strconv.AppendFloat(arr[:0], m, 'e', -1, 64) // e.g., "d.dddde±dd"

		// Parse the exponent.
		i := bytes.IndexByte(b, 'e')
		nk, err := strconv.ParseInt(string(b[i+len("e"):]), 10, 64) // i.e., n-k
		if err != nil {
			panic("BUG: unexpected strconv.ParseInt error: " + err.Error())
		}

		// Format the significand.
		s = b[:i]
		if len(b) > 1 && b[1] == '.' {
			s[1], s = s[0], s[1:] // e.g., "d.dddd" => "ddddd"
			for len(s) > 1 && s[len(s)-1] == '0' {
				s = s[:len(s)-1] // trim trailing zeros
			}
			nk -= int64(len(s) - 1)
		}

		k = int64(len(s)) // k is the number of digits in the decimal representation of s
		n = nk + k        // nk=n-k => n=nk+k
	}

	const zeros = "000000000000000000000"
	switch {
	// 6. If k ≤ n ≤ 21, …
	case k <= n && n <= 21:
		// … return the String consisting of the code units of the k digits of
		//   the decimal representation of s (in order, with no leading zeros), …
		b = append(b, s...)
		// … followed by n-k occurrences of the code unit 0x0030 (DIGIT ZERO).
		b = append(b, zeros[:n-k]...)

	// 7. If 0 < n ≤ 21, …
	case 0 < n && n <= 21:
		// … return the String consisting of the code units of the
		//   most significant n digits of the decimal representation of s, …
		b = append(b, s[:n]...)
		// … followed by the code unit 0x002E (FULL STOP), …
		b = append(b, '.')
		// … followed by the code units of the remaining k-n digits of the decimal representation of s.
		b = append(b, s[n:]...)

	// 8. If -6 < n ≤ 0, …
	case -6 < n && n <= 0:
		// … return the String consisting of the code unit 0x0030 (DIGIT ZERO), …
		b = append(b, '0')
		// … followed by the code unit 0x002E (FULL STOP), …
		b = append(b, '.')
		// … followed by -n occurrences of the code unit 0x0030 (DIGIT ZERO), …
		b = append(b, zeros[:-n]...)
		// … followed by the code units of the k digits of the decimal representation of s.
		b = append(b, s...)

	// 9. If k = 1, …
	case k == 1:
		// … return the String consisting of the code unit of the single digit of s, …
		b = append(b, s...)
		// … followed by code unit 0x0065 (LATIN SMALL LETTER E), …
		b = append(b, 'e')
		// … followed by code unit 0x002B (PLUS SIGN) or the code unit 0x002D (HYPHEN-MINUS) according to whether n-1 is positive or negative, …
		// … followed by the code units of the decimal representation of the integer abs(n-1) (with no leading zeroes).
		switch {
		case n-1 > 0:
			b = strconv.AppendInt(append(b, '+'), n-1, 10)
		case n-1 < 0:
			b = strconv.AppendInt(append(b, '-'), 1-n, 10)
		}

	// 10. Otherwise …
	default:
		// … return the String consisting of the code units of the
		//   most significant digit of the decimal representation of s, …
		b = append(b, s[0])
		// … followed by code unit 0x002E (FULL STOP), …
		b = append(b, '.')
		// … followed by the code units of the remaining k-1 digits of the decimal representation of s, …
		b = append(b, s[1:]...)
		// … followed by code unit 0x0065 (LATIN SMALL LETTER E), …
		b = append(b, 'e')
		// … followed by code unit 0x002B (PLUS SIGN) or the code unit 0x002D (HYPHEN-MINUS) according to whether n-1 is positive or negative, …
		// … followed by the code units of the decimal representation of the integer abs(n-1) (with no leading zeroes).
		switch {
		case n-1 > 0:
			b = strconv.AppendInt(append(b, '+'), n-1, 10)
		case n-1 < 0:
			b = strconv.AppendInt(append(b, '-'), 1-n, 10)
		}
	}
	return b
}
