// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package bits implements bit-twiddling functionality.
package bits

import "unsafe"
import "errors"

const (
	MinUint     = uint(0)
	MaxUint     = ^MinUint
	NumUintBits = int(8 * unsafe.Sizeof(uint(0)))
)

// ErrAlign is returned if the function called cannot perform the requested
// operation because the internal offset is not aligned to a byte-boundary.
var ErrAlign = errors.New("golib/bits: offset is not byte-aligned")

type BitReader interface {
	// Read a single bit.
	// This returns an error if and only if no bit is read.
	ReadBit() (val bool, err error)
}

type BitsReader interface {
	BitReader

	// Read num bits in LSB order. That is, the first bits read are packed into
	// val as the LSB. The behavior is undefined if a read is attempted on more
	// bits than fits in an uint.
	//
	// If an error is encountered while reading a bit, then an error will be
	// returned along with the number of bits read thus far. The error returned
	// will be io.EOF only if the cnt is 0. Otherwise, io.ErrUnexpectedEOF will
	// be used. This is done to match the behavior of io.ReadFull.
	ReadBits(num int) (val uint, cnt int, err error)
}

type BitWriter interface {
	// Write a single bit.
	// This returns an error if and only if no bit is written.
	WriteBit(bool) error
}

type BitsWriter interface {
	BitWriter

	// Write num bits in LSB order. That is, the LSB bits of val will be the
	// first bits to be written. The behavior is undefined if a read is
	// attempted on more bits than fits in an uint.
	//
	// If an error is encountered while writing a bit, then an error will be
	// returned along with the number of bits written thus far.
	WriteBits(val uint, num int) (cnt int, err error)
}
