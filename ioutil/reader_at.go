// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package ioutil

import "io"
import "os"
import "sync"

// Given a ReadSeeker, create a new object that satisfies the ReaderAt
// interface. Concurrent ReadAt operations are permitted due to the use of a
// mutex to ensure synchronization. However, any regular Read and Seek calls
// must also be made through this class to ensure safe concurrent action.
type ReaderAt struct {
	io.ReadSeeker
	L *sync.Mutex
}

// Create a new ReaderAt from the given ReadSeeker.
func NewReaderAt(rd io.ReadSeeker) *ReaderAt {
	return &ReaderAt{rd, new(sync.Mutex)}
}

func (r *ReaderAt) ReadAt(data []byte, off int64) (cnt int, err error) {
	r.L.Lock()
	defer r.L.Unlock()

	var pos int64
	if pos, err = r.Seek(off, os.SEEK_SET); err != nil {
		return
	}
	defer func() {
		if _, skErr := r.Seek(pos, os.SEEK_SET); skErr != nil {
			err = skErr
		}
	}()

	return r.Read(data)
}

func (r *ReaderAt) Read(data []byte) (cnt int, err error) {
	r.L.Lock()
	defer r.L.Unlock()
	return r.ReadSeeker.Read(data)
}

func (r *ReaderAt) Seek(off int64, whence int) (int64, error) {
	r.L.Lock()
	defer r.L.Unlock()
	return r.ReadSeeker.Seek(off, whence)
}
