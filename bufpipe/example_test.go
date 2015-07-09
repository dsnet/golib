// Copyright 2014, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package bufpipe

import "io"
import "fmt"
import "time"
import "sync"
import "math/rand"

func randomChars(cnt int, rand *rand.Rand) []byte {
	data := make([]byte, cnt)
	for idx := range data {
		char := byte(rand.Intn(10+26+26))
		if char < 10 {
			data[idx] = '0' + char
		} else if char < 10+26 {
			data[idx] = 'A' + char - 10
		} else {
			data[idx] = 'a' + char - 36
		}
	}
	return data
}

func ExampleLineMonoPipe() {
	// The buffer is large enough such that the producer doesn't overfill it.
	buffer := bufpipe.NewBufferPipe(make([]byte, 4096), bufpipe.LineMonoIO)

	rand := rand.New(rand.NewSource(0))
	group := new(sync.WaitGroup)
	group.Add(2)

	// Producer routine.
	go func() {
		defer group.Done()
		defer buffer.Close()

		// In LineMonoIO mode only, it is safe to store a reference to written
		// data and modify later.
		header := buffer.WriteSlice()[:4]

		totalCnt := 0
		buffer.Write([]byte("#### "))
		for idx := 0; idx < 10; idx++ {
			data := randomChars(rand.Intn(64), rand)

			// So long as the amount of data written has not exceeded the size
			// of the buffer, Write() will never fail.
			buffer.Write([]byte(data))
			totalCnt += len(data)

			time.Sleep(100 * time.Millisecond)
		}

		// Write the header afterwards
		copy(header, fmt.Sprintf("%04d", totalCnt))
	}()

	// Consumer routine.
	go func() {
		defer group.Done()

		// In LineMonoIO mode only, a call to ReadSlice() is guaranteed to block
		// until the channel is closed. All written data will be made available.
		data := buffer.ReadSlice()
		buffer.ReadMark(len(data)) // Technically, this is optional

		fmt.Println(string(data))
	}()

	group.Wait()
}

func ExampleLineDualPipe() {
	// The buffer is large enough such that the producer doesn't overfill it.
	buffer := bufpipe.NewBufferPipe(make([]byte, 4096), bufpipe.LineDualIO)

	rand := rand.New(rand.NewSource(0))
	group := new(sync.WaitGroup)
	group.Add(2)

	// Producer routine.
	go func() {
		defer group.Done()
		defer buffer.Close()

		buffer.Write([]byte("#### ")) // Write a fake header
		for idx := 0; idx < 10; idx++ {
			data := randomChars(rand.Intn(64), rand)

			// So long as the amount of data written has not exceeded the size
			// of the buffer, Write() will never fail.
			buffer.Write([]byte(data))

			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Consumer routine.
	go func() {
		defer group.Done()
		for {
			// Reading can be also done using ReadSlice() and ReadMark() pairs.
			data := buffer.ReadSlice()
			if len(data) == 0 {
				break
			}
			buffer.ReadMark(len(data))
			fmt.Print(string(data))
		}
		fmt.Println()
	}()

	group.Wait()
}

func ExampleRingDualPipe() {
	// Intentionally small buffer to show that data written into the buffer
	// can exceed the size of the buffer itself.
	buffer := bufpipe.NewBufferPipe(make([]byte, 128), bufpipe.RingDualIO)

	rand := rand.New(rand.NewSource(0))
	group := new(sync.WaitGroup)
	group.Add(2)

	// Producer routine.
	go func() {
		defer group.Done()
		defer buffer.Close()

		buffer.Write([]byte("#### ")) // Write a fake header
		for idx := 0; idx < 10; idx++ {
			data := randomChars(rand.Intn(64), rand)

			// So long as the amount of data written has not exceeded the size
			// of the buffer, Write() will never fail.
			buffer.Write([]byte(data))

			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Consumer routine.
	go func() {
		defer group.Done()

		data := make([]byte, 64)
		for {
			// Reading can also be done using the Read() method.
			cnt, err := buffer.Read(data)
			fmt.Print(string(data[:cnt]))
			if err == io.EOF {
				break
			}
		}
		fmt.Println()
	}()

	group.Wait()
}
