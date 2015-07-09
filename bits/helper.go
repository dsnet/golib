// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"

// Convert boolean false to 0 and true to 1.
func Btoi(b bool) uint {
	if b {
		return 1
	}
	return 0
}

// Convert integer to boolean false if equal to 0, otherwise true.
func Itob(i uint) bool {
	return i != 0
}

// Get the bit at a given offset.
func Get(buf []byte, ofs int) (val bool) {
	d := buf[ofs/8]
	m := byte(1 << uint(ofs%8))
	return d&m > 0
}

// Get the value of cnt bits at the given offset.
// Value is read in LSB-first order.
func GetN(buf []byte, cnt, ofs int) (val uint) {
	i := ofs / 8
	m := byte(1 << uint(ofs%8))
	for idx := 0; idx < cnt; idx++ {
		if buf[i]&m > 0 {
			val |= (0x01 << uint(idx))
		}
		m <<= 1
		if m == 0 {
			m = 0x01
			i++
		}
	}
	return
}

// Set the value of the bit at a given offset.
func Set(data []byte, val bool, ofs int) {
	if m := byte(1 << uint(ofs%8)); val {
		data[ofs/8] |= m
	} else {
		data[ofs/8] &= ^m
	}
}

// Set the value of cnt bits at the given offset.
// Value is written in LSB-first order.
func SetN(data []byte, val uint, cnt, ofs int) {
	i := ofs / 8
	m := byte(1 << uint(ofs%8))
	for idx := 0; idx < cnt; idx++ {
		if val&(0x01<<uint(idx)) > 0 {
			data[i] |= m
		} else {
			data[i] &= ^m
		}
		m <<= 1
		if m == 0 {
			m = 0x01
			i++
		}
	}
}

// Count the number of zeros and ones in the slice.
func Count(data []byte) (zeros, ones int) {
	for _, val := range data {
		for idx := 0; idx < 8 && val > 0; idx++ {
			ones += int(val & 1)
			val >>= 1
		}
	}
	return 8*len(data) - ones, ones
}

// Invert all bits in the slice.
func Invert(data []byte) {
	for idx := range data {
		data[idx] ^= 0xff
	}
}

// Reverse the bits of val.
func ReverseUint(val uint) uint {
	var w = uint(NumUintBits)
	var m = uint(MaxUint)
	for w > 1 {
		w >>= 1
		m ^= m >> w
		val = ((val & m) >> w) | ((val &^ m) << w)
	}
	return val
}

// Reverse the lower num bits of val. The upper bits will be zero.
func ReverseUintN(val uint, num int) uint {
	val = ReverseUint(val)
	val >>= uint(NumUintBits - num)
	return val
}

// Read num bits from br in LSB order. That is, the first bits read from br are
// packed into val as the LSB. The behavior is undefined if a read is attempted
// on more bits than fits in an uint.
//
// If an error is encountered while reading a bit, then an error will be
// returned along with the number of bits read thus far. The error returned will
// be io.EOF only if the cnt is 0. Otherwise, io.ErrUnexpectedEOF will be used.
// This is done to match the behavior of io.ReadFull.
func ReadBits(br BitReader, num int) (val uint, cnt int, err error) {
	var bit bool
	for cnt = 0; cnt < num; cnt++ {
		if bit, err = br.ReadBit(); err != nil {
			if err == io.EOF && cnt > 0 {
				err = io.ErrUnexpectedEOF
			}
			return
		}
		val |= (Btoi(bit) << uint(cnt))
	}
	return
}

// Write num bits to bw in LSB order. That is, the LSB bits of val will be the
// first bits to be written to bw. The behavior is undefined if a read is
// attempted on more bits than fits in an uint.
//
// If an error is encountered while writing a bit, then an error will be
// returned along with the number of bits written thus far.
func WriteBits(bw BitWriter, val uint, num int) (cnt int, err error) {
	var bit bool
	for cnt = 0; cnt < num; cnt++ {
		bit = Itob(val & (1 << uint(cnt)))
		if err = bw.WriteBit(bit); err != nil {
			return
		}
	}
	return
}
