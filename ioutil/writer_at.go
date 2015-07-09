// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"
import "os"
import "sync"

// Given a WriteSeeker, create a new object that satisfies the WriterAt
// interface. Concurrent WriteAt operations are permitted due to the use of a
// mutex to ensure synchronization. However, any regular Write and Seek calls
// must also be made through this class to ensure safe concurrent action.
type WriterAt struct {
	io.WriteSeeker
	L *sync.Mutex
}

// Create a new WriterAt from the given WriteSeeker.
func NewWriterAt(rd io.WriteSeeker) *WriterAt {
	return &WriterAt{rd, new(sync.Mutex)}
}

func (w *WriterAt) WriteAt(data []byte, off int64) (cnt int, err error) {
	w.L.Lock()
	defer w.L.Unlock()

	var pos int64
	if pos, err = w.Seek(off, os.SEEK_SET); err != nil {
		return
	}
	defer func() {
		if _, skErr := w.Seek(pos, os.SEEK_SET); skErr != nil {
			err = skErr
		}
	}()

	return w.Write(data)
}

func (w *WriterAt) Write(data []byte) (cnt int, err error) {
	w.L.Lock()
	defer w.L.Unlock()
	return w.WriteSeeker.Write(data)
}

func (w *WriterAt) Seek(off int64, whence int) (int64, error) {
	w.L.Lock()
	defer w.L.Unlock()
	return w.WriteSeeker.Seek(off, whence)
}
