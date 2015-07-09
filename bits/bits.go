// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "unsafe"
import "bitbucket.org/rawr/golib/errs"

const (
	MaxUint     = ^uint(0)
	NumUintBits = int(8 * unsafe.Sizeof(uint(0)))
)

var ErrAlign = errs.New("golib/bits: offset is not byte-aligned")

type BitReader interface {
	// Read a single bit.
	// This returns an error if and only if no bit is read.
	ReadBit() (val bool, err error)
}

type BitCountReader interface {
	BitReader

	// Total number of bits read.
	BitsRead() int64
}

type BitWriter interface {
	// Write a single bit.
	// This returns an error if and only if no bit is written.
	WriteBit(bool) error
}

type BitCountWriter interface {
	BitWriter

	// Total number of bits written.
	BitsWritten() int64
}
