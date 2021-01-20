// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsoncs

import (
	"math/rand"
	"sort"
	"testing"
	"unicode/utf16"

	"github.com/google/go-cmp/cmp"
)

func TestLessUTF16(t *testing.T) {
	want := []string{"", "\r", "1", "\u0080", "\u00f6", "\u20ac", "\U0001f600", "\ufb33"}

	got1 := append([]string(nil), want...)
	got2 := append([]string(nil), want...)
	for i, j := range rand.Perm(len(want)) {
		got1[i], got1[j] = got1[j], got1[i]
		got2[i], got2[j] = got2[j], got2[i]
	}

	// Sort using optimized lessUTF16 implementation.
	sort.Slice(got1, func(i, j int) bool {
		return lessUTF16(got1[i], got1[j])
	})
	if diff := cmp.Diff(want, got1); diff != "" {
		t.Errorf("sort.Slice(LessUTF16.Optimized) mismatch (-want +got)\n%s", diff)
	}

	// Sort using simple, but slow lessUTF16 implementation.
	lessUTF16 := func(x, y string) bool {
		ux := utf16.Encode([]rune(x))
		uy := utf16.Encode([]rune(y))
		for {
			if len(ux) == 0 || len(uy) == 0 {
				return len(ux) < len(uy)
			}
			if ux[0] != uy[0] {
				return ux[0] < uy[0]
			}
			ux, uy = ux[1:], uy[1:]
		}
	}
	sort.Slice(got2, func(i, j int) bool {
		return lessUTF16(got2[i], got2[j])
	})
	if diff := cmp.Diff(want, got2); diff != "" {
		t.Errorf("sort.Slice(LessUTF16.Simplified) mismatch (-want +got)\n%s", diff)
	}
}
