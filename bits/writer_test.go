// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bits

import "io"
import "bytes"
import "runtime"
import "testing"
import "github.com/stretchr/testify/assert"

type buffer struct {
	bytes.Buffer
	fail bool
}

func (b *buffer) WriteByte(val byte) (err error) {
	if b.fail {
		return io.ErrShortWrite
	}
	return b.Buffer.WriteByte(val)
}

func TestWriter(t *testing.T) {
	type X struct {
		wa      bool
		wbn, wn int64
		buf     []byte
	}
	state := func(bw *Writer) X {
		return X{
			bw.WriteAligned(),
			bw.BytesWritten(), bw.BitsWritten(),
			nb(bw.wr.(*buffer).Bytes()),
		}
	}

	var cnt int
	var err error

	b := new(buffer)
	bw := NewWriter(b)
	assert.Equal(t, X{true, 0, 0, nil}, state(bw))

	err = bw.WriteBit(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, 0, 1, nil}, state(bw))

	err = bw.WriteBit(false)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, 0, 2, nil}, state(bw))

	err = bw.WriteBit(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, 0, 3, nil}, state(bw))

	err = bw.WriteByte(0xff)
	assert.Equal(t, ErrAlign, err)
	assert.Equal(t, X{false, 0, 3, nil}, state(bw))

	cnt, err = bw.WriteBits(0xab3, 12)
	assert.Equal(t, 12, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, 1, 15, []byte{0x9d}}, state(bw))

	err = bw.WriteBit(true)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, 2, 16, []byte{0x9d, 0xd5}}, state(bw))

	err = bw.WriteByte(0xa7)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, 3, 24, []byte{0x9d, 0xd5, 0xa7}}, state(bw))

	b.fail = true
	assert.Equal(t, io.ErrShortWrite, bw.WriteByte(0xff))
	assert.Equal(t, X{true, 3, 24, []byte{0x9d, 0xd5, 0xa7}}, state(bw))

	cnt, err = bw.WriteBits(0x74, 7)
	assert.Equal(t, 7, cnt)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{false, 3, 31, []byte{0x9d, 0xd5, 0xa7}}, state(bw))

	cnt, err = bw.WriteBits(0x3, 2)
	assert.Equal(t, 0, cnt)
	assert.Equal(t, io.ErrShortWrite, err)
	assert.Equal(t, X{false, 3, 31, []byte{0x9d, 0xd5, 0xa7}}, state(bw))

	err = bw.WriteBit(false)
	assert.Equal(t, io.ErrShortWrite, err)
	assert.Equal(t, X{false, 3, 31, []byte{0x9d, 0xd5, 0xa7}}, state(bw))

	b.fail = false
	err = bw.WriteBit(false)
	assert.Equal(t, nil, err)
	assert.Equal(t, X{true, 4, 32, []byte{0x9d, 0xd5, 0xa7, 0x74}}, state(bw))

	// Reset
	b.Reset()
	bw.Reset(b)
	assert.Equal(t, X{true, 0, 0, nil}, state(bw))
}

func BenchmarkWriter(b *testing.B) {
	cnt := 1 << 20 // 1 MiB
	buf := NewBuffer(make([]byte, 0, cnt))
	bw := NewWriter(buf)
	bww := BitWriter(bw)

	runtime.GC()
	b.ReportAllocs()
	b.SetBytes(int64(cnt))
	b.ResetTimer()

	for ni := 0; ni < b.N; ni++ {
		buf.Reset()
		bw.Reset(buf)
		for bi := 0; bi < cnt; bi++ {
			bww.WriteBit(true)
			bww.WriteBit(false)
			bww.WriteBit(true)
			bww.WriteBit(false)
			bww.WriteBit(true)
			bww.WriteBit(false)
			bww.WriteBit(true)
			bww.WriteBit(false)
		}
	}
}
