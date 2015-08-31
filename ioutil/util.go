// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package ioutil is a collection of small io related implementations.
package ioutil

import "io"
import "os"
import "math"
import "github.com/dsnet/golib/errs"

// Define constants for seeking here so that package os is not needed as well.
const (
	SeekCur = os.SEEK_CUR
	SeekSet = os.SEEK_SET
	SeekEnd = os.SEEK_END
)

// Determine the size of a ReaderAt using a binary search. Given that file
// offsets are no larger than int64, there is an upper limit of 64 iterations
// before the EOF is found.
func ReaderAtSize(rd io.ReaderAt) (pos int64, err error) {
	defer errs.Recover(&err)

	// Function to check if the given position is at EOF
	buf := make([]byte, 2)
	checkEOF := func(pos int64) int {
		if pos > 0 {
			cnt, err := rd.ReadAt(buf[:2], pos-1)
			errs.Panic(errs.Ignore(err, io.EOF))
			return 1 - cnt // RetVal[Cnt] = {0: +1, 1: 0, 2: -1}
		} else { // Special case where position is zero
			cnt, err := rd.ReadAt(buf[:1], pos-0)
			errs.Panic(errs.Ignore(err, io.EOF))
			return 0 - cnt // RetVal[Cnt] = {0: 0, 1: -1}
		}
	}

	// Obtain the size via binary search O(log n) => 64 iterations
	posMin, posMax := int64(0), int64(math.MaxInt64)
	for posMax >= posMin {
		pos = (posMax + posMin) / 2
		switch checkEOF(pos) {
		case -1: // Below EOF
			posMin = pos + 1
		case 0: // At EOF
			return pos, nil
		case +1: // Above EOF
			posMax = pos - 1
		}
	}
	panic(errs.New("EOF is in a transient state"))
}

// Determine the size of a Seeker by seeking to the end. This function will
// attempt to bring the file pointer back to the original location.
func SeekerSize(sk io.Seeker) (pos int64, err error) {
	var curPos int64
	if curPos, err = sk.Seek(0, SeekCur); err != nil {
		return
	}
	if pos, err = sk.Seek(0, SeekEnd); err != nil {
		return
	}
	if _, err = sk.Seek(curPos, SeekSet); err != nil {
		return
	}
	return
}

// Performs like io.ReadFull, but uses a ByteReader instead.
func ByteReadFull(rd io.ByteReader, buf []byte) (cnt int, err error) {
	for cnt = 0; cnt < len(buf); cnt++ {
		buf[cnt], err = rd.ReadByte()
		if err == io.EOF && cnt > 0 {
			err = io.ErrUnexpectedEOF
		}
		if err != nil {
			return
		}
	}
	return
}

// Allows multiple bytes to be written to a buffer like io.Write.
func ByteWriteFull(wr io.ByteWriter, buf []byte) (cnt int, err error) {
	for cnt = 0; cnt < len(buf); cnt++ {
		err = wr.WriteByte(buf[cnt])
		if err != nil {
			return
		}
	}
	return
}

func ByteCopy(dst io.ByteWriter, src io.ByteReader) (cnt int64, err error) {
	var val byte
	for cnt = 0; true; cnt++ {
		if val, err = src.ReadByte(); err != nil {
			break
		}
		if err = dst.WriteByte(val); err != nil {
			break
		}
	}
	if err == io.EOF { // This is expected
		err = nil
	}
	return
}

func ByteCopyN(dst io.ByteWriter, src io.ByteReader, num int64) (cnt int64, err error) {
	var val byte
	for cnt = 0; cnt < num; cnt++ {
		if val, err = src.ReadByte(); err != nil {
			break
		}
		if err = dst.WriteByte(val); err != nil {
			break
		}
	}
	if cnt == num { // This is expected
		err = nil
	}
	return
}
