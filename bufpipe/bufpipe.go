// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package bufpipe implements a buffered pipe.
package bufpipe

import "io"
import "sync"

const (
	LineMonoIO = iota
	LineDualIO
	RingDualIO
)

type BufferPipe struct {
	buf    []byte
	mode   int
	rdCnt  int64
	wrCnt  int64
	closed bool
	mutex  sync.Mutex
	rdCond sync.Cond
	wrCond sync.Cond
}

// BufferPipe is similar in operation to io.Pipe and is intended to be the
// communication channel between producer and consumer routines. There are some
// specific use cases for BufferPipes over io.Pipe.
//
// First, in cases where a writer may need to go back a modify a portion of the
// past "written" data before allowing the reader to consume it.
// Second, for performance applications, where the cost of copying of data is
// noticeable. Thus, it would be more efficient to read/write data from/to the
// internal buffer directly.
//
// The BufferPipe allows for several modes:
//  * LineMonoIO: Acts like a linear buffer. A writer can produce at most as
//    much data as the size of the internal buffer. Also, a reader is blocked
//    on reading until the writer explicitly closes the pipe.
//  * LineDualIO: Operates like the previous mode, but only blocks readers until
//    there is at least some data.
//  * RingDualIO: Operates in a similar fashion as io.Pipe.
func NewBufferPipe(buf []byte, mode int) *BufferPipe {
	switch mode {
	case LineMonoIO, LineDualIO, RingDualIO:
	default:
		panic("unknown buffer IO mode")
	}
	if len(buf) == 0 {
		panic("buffer may not be empty")
	}

	b := new(BufferPipe)
	b.buf = buf
	b.mode = mode
	b.rdCond.L = &b.mutex
	b.wrCond.L = &b.mutex
	return b
}

// The total number of bytes the buffer can store.
func (b *BufferPipe) Capacity() int {
	return len(b.buf)
}

// The number of valid bytes that can be read.
func (b *BufferPipe) Length() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return int(b.wrCnt - b.rdCnt)
}

func (b *BufferPipe) writeWait() int {
	var offLo int64 // For Linear buffer, this is always zero
	if b.mode == RingDualIO {
		for !b.closed && len(b.buf) == int(b.wrCnt-b.rdCnt) {
			b.wrCond.Wait()
		}
		offLo = b.rdCnt
	}
	if b.closed {
		return 0 // Closed buffer is never available
	}
	return len(b.buf) - int(b.wrCnt-offLo)
}

// Slice of available buffer that can be written to. This does not advance the
// internal write pointer.
//
// In linear buffers, the slice obtained is guaranteed to be the entire
// available writable buffer slice.
//
// In LineMonoIO mode only, slices obtained may still be modified even after
// WriteMark() has been called and before Close(). This is valid because this
// mode blocks readers until the buffer has been closed.
//
// In ring buffers, the slice obtained may not represent all of the available
// buffer space since this method always returns contiguous pieces of memory.
//
// In the RingDualIO mode only, this method blocks until there is available
// space in the buffer to write. Other modes do not block and will return 0 if
// the buffer is full.
func (b *BufferPipe) WriteSlice() []byte {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.writeSlice()
}

func (b *BufferPipe) writeSlice() []byte {
	if b.mode == RingDualIO {
		availCnt := b.writeWait() // Block until there is available buffer

		offLo := int(b.wrCnt) % len(b.buf)
		offHi := offLo + availCnt
		if offHi > len(b.buf) { // If available slice is split, take bottom
			offHi = len(b.buf)
		}
		return b.buf[offLo:offHi] // Ring buffer
	}
	return b.buf[b.wrCnt:] // Linear buffer
}

// Advances the write pointer.
//
// The amount that can be advanced must be non-negative and be less than the
// length of the slice returned by the previous WriteSlice(). Calls to Write()
// may not be done between these two calls. Also, another call to WriteMark()
// is invalid until WriteSlice() has been called again.
//
// If WriteMark() is being used, only one writer routine is allowed.
func (b *BufferPipe) WriteMark(cnt int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.writeMark(cnt)
}

func (b *BufferPipe) writeMark(cnt int) {
	availCnt := b.writeWait()
	if cnt < 0 || cnt > availCnt {
		panic("invalid mark increment value")
	}
	b.wrCnt += int64(cnt)

	b.rdCond.Signal()
}

// Write data to the buffer.
//
// In linear buffers, the length of the data slice may not exceed the capacity
// of the buffer. Otherwise, an ErrShortWrite error will be returned.
//
// In ring buffers, the length of the data slice may exceed the capacity.
// The operation will block until all data has been written. If there is no
// consumer of the data, then this method may block forever.
func (b *BufferPipe) Write(data []byte) (cnt int, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for cnt < len(data) {
		buf := b.writeSlice()
		if len(buf) == 0 {
			if b.closed {
				return cnt, io.ErrClosedPipe
			} else {
				return cnt, io.ErrShortWrite
			}
		}

		copyCnt := copy(buf, data[cnt:])
		b.writeMark(copyCnt)
		cnt += copyCnt
	}
	return cnt, nil
}

func (b *BufferPipe) readWait() int {
	for !b.closed && (b.rdCnt == b.wrCnt || b.mode == LineMonoIO) {
		b.rdCond.Wait()
	}
	return int(b.wrCnt - b.rdCnt)
}

// Slice of valid data that can be read. This does not advance the internal
// read pointer.
//
// In linear buffers, the slice obtained is guaranteed to be the entire
// valid readable buffer slice.
//
// In ring buffers, the slice obtained may not represent all of the valid buffer
// data since this method always returns contiguous pieces of memory.
//
// In all modes, this method blocks until there is some valid data to read.
// The LineMonoIO mode is special in that it will block until the buffer has
// been closed. Other modes just block until some data is available.
func (b *BufferPipe) ReadSlice() []byte {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.readSlice()
}

func (b *BufferPipe) readSlice() []byte {
	if b.mode == RingDualIO {
		validCnt := b.readWait() // Block until there is valid buffer

		offLo := int(b.rdCnt) % len(b.buf)
		offHi := offLo + validCnt
		if offHi > len(b.buf) { // If valid slice is split, take bottom
			offHi = len(b.buf)
		}
		return b.buf[offLo:offHi] // Ring buffer
	}
	return b.buf[b.rdCnt:b.wrCnt] // Linear buffer
}

// Advances the read pointer.
//
// The amount that can be advanced must be non-negative and be less than the
// length of the slice returned by the previous ReadSlice(). Calls to Read()
// may not be done between these two calls. Also, another call to ReadMark()
// is invalid until ReadSlice() has been called again.
//
// If ReadMark() is being used, only one writer routine is allowed.
func (b *BufferPipe) ReadMark(cnt int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.readMark(cnt)
}

func (b *BufferPipe) readMark(cnt int) {
	validCnt := b.readWait()
	if cnt < 0 || cnt > validCnt {
		panic("invalid mark increment value")
	}
	b.rdCnt += int64(cnt)

	b.wrCond.Signal()
}

// Read data from the buffer.
//
// In all modes, the length of the data slice may exceed the capacity of
// the buffer. The operation will block until all data has been read or until
// the EOF is hit. If there is no producer of the data, then this method may
// block forever.
func (b *BufferPipe) Read(data []byte) (cnt int, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for cnt < len(data) {
		buf := b.readSlice()
		if len(buf) == 0 {
			return cnt, io.EOF
		}

		copyCnt := copy(data[cnt:], buf)
		b.readMark(copyCnt)
		cnt += copyCnt
	}
	return cnt, nil
}

// Close the buffer down.
//
// All write operations have no effect after this, while all read operations
// will drain remaining data in the buffer. This operation is somewhat similar
// to how Go channels operation.
//
// Writers should close the buffer to indicate to readers to mark end-of-stream.
//
// Readers should only close the buffer in the event of unexpected termination.
// The mechanism allows readers to inform writers of consumer termination and
// prevents the producer from potentially being blocked forever.
func (b *BufferPipe) Close() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.closed = true
	b.rdCond.Signal()
	b.wrCond.Signal()
	return nil
}

// Makes the buffer ready for use again.
func (b *BufferPipe) Reset() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.wrCnt, b.rdCnt = 0, 0
	b.closed = false
}
