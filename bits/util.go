// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"

// Lookup table that match each byte to the number of one bits in that byte.
var cntOnes = [256]byte{
	0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	4, 5, 5, 6, 5, 6, 6, 7, 5, 6, 6, 7, 6, 7, 7, 8,
}

// Lookup table that reverses each bit in a byte.
var revByte = [256]byte{
	0x00, 0x80, 0x40, 0xc0, 0x20, 0xa0, 0x60, 0xe0, 0x10, 0x90, 0x50, 0xd0, 0x30, 0xb0, 0x70, 0xf0,
	0x08, 0x88, 0x48, 0xc8, 0x28, 0xa8, 0x68, 0xe8, 0x18, 0x98, 0x58, 0xd8, 0x38, 0xb8, 0x78, 0xf8,
	0x04, 0x84, 0x44, 0xc4, 0x24, 0xa4, 0x64, 0xe4, 0x14, 0x94, 0x54, 0xd4, 0x34, 0xb4, 0x74, 0xf4,
	0x0c, 0x8c, 0x4c, 0xcc, 0x2c, 0xac, 0x6c, 0xec, 0x1c, 0x9c, 0x5c, 0xdc, 0x3c, 0xbc, 0x7c, 0xfc,
	0x02, 0x82, 0x42, 0xc2, 0x22, 0xa2, 0x62, 0xe2, 0x12, 0x92, 0x52, 0xd2, 0x32, 0xb2, 0x72, 0xf2,
	0x0a, 0x8a, 0x4a, 0xca, 0x2a, 0xaa, 0x6a, 0xea, 0x1a, 0x9a, 0x5a, 0xda, 0x3a, 0xba, 0x7a, 0xfa,
	0x06, 0x86, 0x46, 0xc6, 0x26, 0xa6, 0x66, 0xe6, 0x16, 0x96, 0x56, 0xd6, 0x36, 0xb6, 0x76, 0xf6,
	0x0e, 0x8e, 0x4e, 0xce, 0x2e, 0xae, 0x6e, 0xee, 0x1e, 0x9e, 0x5e, 0xde, 0x3e, 0xbe, 0x7e, 0xfe,
	0x01, 0x81, 0x41, 0xc1, 0x21, 0xa1, 0x61, 0xe1, 0x11, 0x91, 0x51, 0xd1, 0x31, 0xb1, 0x71, 0xf1,
	0x09, 0x89, 0x49, 0xc9, 0x29, 0xa9, 0x69, 0xe9, 0x19, 0x99, 0x59, 0xd9, 0x39, 0xb9, 0x79, 0xf9,
	0x05, 0x85, 0x45, 0xc5, 0x25, 0xa5, 0x65, 0xe5, 0x15, 0x95, 0x55, 0xd5, 0x35, 0xb5, 0x75, 0xf5,
	0x0d, 0x8d, 0x4d, 0xcd, 0x2d, 0xad, 0x6d, 0xed, 0x1d, 0x9d, 0x5d, 0xdd, 0x3d, 0xbd, 0x7d, 0xfd,
	0x03, 0x83, 0x43, 0xc3, 0x23, 0xa3, 0x63, 0xe3, 0x13, 0x93, 0x53, 0xd3, 0x33, 0xb3, 0x73, 0xf3,
	0x0b, 0x8b, 0x4b, 0xcb, 0x2b, 0xab, 0x6b, 0xeb, 0x1b, 0x9b, 0x5b, 0xdb, 0x3b, 0xbb, 0x7b, 0xfb,
	0x07, 0x87, 0x47, 0xc7, 0x27, 0xa7, 0x67, 0xe7, 0x17, 0x97, 0x57, 0xd7, 0x37, 0xb7, 0x77, 0xf7,
	0x0f, 0x8f, 0x4f, 0xcf, 0x2f, 0xaf, 0x6f, 0xef, 0x1f, 0x9f, 0x5f, 0xdf, 0x3f, 0xbf, 0x7f, 0xff,
}

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

// Invert all bits in the slice.
func Invert(data []byte) {
	for idx := range data {
		data[idx] ^= 0xff
	}
}

// Count the number of one bits in the slice.
func Count(data []byte) (ones int) {
	for _, val := range data {
		ones += int(cntOnes[val])
	}
	return ones
}

// Count the number of one bits in the byte.
func CountByte(val byte) (ones int) {
	return int(cntOnes[val])
}

// Count the number of one bits in the uint.
func CountUint(val uint) (ones int) {
	for val > 0 {
		ones += int(cntOnes[val&0xff])
		val >>= 8
	}
	return ones
}

// Reverse the bits of every byte in the slice.
func Reverse(data []byte) {
	for i, b := range data {
		data[i] = revByte[b]
	}
}

// Reverse the bits of a byte.
func ReverseByte(val byte) byte {
	return revByte[val]
}

// Reverse the bits of a uint.
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

// Efficiently write the same bit.
func WriteSameBit(bw BitsWriter, val bool, num int) (cnt int, err error) {
	var mask = MinUint
	if val {
		mask = MaxUint
	}

	for num > 0 && err == nil {
		wrCnt := num
		if wrCnt > NumUintBits {
			wrCnt = NumUintBits
		}

		wrCnt, err = bw.WriteBits(mask, wrCnt)
		num -= wrCnt
		cnt += wrCnt
	}
	return
}

// This function allows a BitReader to easily satisfy the BitsReader interface.
func ReadBits(br BitReader, num int) (val uint, cnt int, err error) {
	// Since br.ReadBit is called fairly often and function calls become a
	// bottleneck, this logic is manually inlined in other bit readers.
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

// This function allows a BitWriter to easily satisfy the BitsWriter interface.
func WriteBits(bw BitWriter, val uint, num int) (cnt int, err error) {
	// Since bw.WriteBit is called fairly often and function calls become a
	// bottleneck, this logic is manually inlined in other bit writers.
	var bit bool
	for cnt = 0; cnt < num; cnt++ {
		bit = Itob(val & (1 << uint(cnt)))
		if err = bw.WriteBit(bit); err != nil {
			return
		}
	}
	return
}
