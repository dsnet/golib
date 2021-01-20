// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsoncs

import (
	"math"
	"testing"
)

func TestFormatES6(t *testing.T) {
	tests := []struct {
		in   float64
		want string
	}{
		{math.NaN(), "NaN"},
		{math.Inf(-1), "-Infinity"},
		{math.Inf(+1), "Infinity"},
		{math.Copysign(0, -1), "0"},
		{math.Copysign(0, +1), "0"},

		// The following cases are from RFC 8785, Appendix B.
		{math.Float64frombits(0x0000000000000000), "0"},
		{math.Float64frombits(0x8000000000000000), "0"},
		{math.Float64frombits(0x0000000000000001), "5e-324"},
		{math.Float64frombits(0x8000000000000001), "-5e-324"},
		{math.Float64frombits(0x7fefffffffffffff), "1.7976931348623157e+308"},
		{math.Float64frombits(0xffefffffffffffff), "-1.7976931348623157e+308"},
		{math.Float64frombits(0x4340000000000000), "9007199254740992"},
		{math.Float64frombits(0xc340000000000000), "-9007199254740992"},
		{math.Float64frombits(0x4430000000000000), "295147905179352830000"},
		{math.Float64frombits(0x7fffffffffffffff), "NaN"},
		{math.Float64frombits(0x7ff0000000000000), "Infinity"},
		{math.Float64frombits(0x44b52d02c7e14af5), "9.999999999999997e+22"},
		{math.Float64frombits(0x44b52d02c7e14af6), "1e+23"},
		{math.Float64frombits(0x44b52d02c7e14af7), "1.0000000000000001e+23"},
		{math.Float64frombits(0x444b1ae4d6e2ef4e), "999999999999999700000"},
		{math.Float64frombits(0x444b1ae4d6e2ef4f), "999999999999999900000"},
		{math.Float64frombits(0x444b1ae4d6e2ef50), "1e+21"},
		{math.Float64frombits(0x3eb0c6f7a0b5ed8c), "9.999999999999997e-7"},
		{math.Float64frombits(0x3eb0c6f7a0b5ed8d), "0.000001"},
		{math.Float64frombits(0x41b3de4355555553), "333333333.3333332"},
		{math.Float64frombits(0x41b3de4355555554), "333333333.33333325"},
		{math.Float64frombits(0x41b3de4355555555), "333333333.3333333"},
		{math.Float64frombits(0x41b3de4355555556), "333333333.3333334"},
		{math.Float64frombits(0x41b3de4355555557), "333333333.33333343"},
		{math.Float64frombits(0xbecbf647612f3696), "-0.0000033333333333333333"},
		{math.Float64frombits(0x43143ff3c1cb0959), "1424953923781206.2"},
	}

	for _, tt := range tests {
		got := string(formatES6(nil, tt.in))
		if got != tt.want {
			t.Errorf("formatES6(%016x) = %v, want %v", math.Float64bits(tt.in), got, tt.want)
		}
	}
}
