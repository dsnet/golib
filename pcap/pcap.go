// Copyright 2022, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package pcap implements reader and writer for the pcap file format.
package pcap

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/bits"
	"time"
)

const (
	magicMicros = 0xa1b2c3d4
	magicNanos  = 0xa1b23c4d

	globalHeaderSize = 24
	packetHeaderSize = 16

	defaultBufSize = 64 << 10
)

// LinkType specifies the type of each Packet.Data.
type LinkType uint32

const (
	EthernetLinkType LinkType = 1   // ethernet; IEEE 802.3
	RawLinkType      LinkType = 101 // raw IP; either IPv4 or IPv6
)

func (t LinkType) String() string {
	switch t {
	case EthernetLinkType:
		return "EthernetLinkType"
	case RawLinkType:
		return "RawLinkType"
	default:
		return fmt.Sprintf("LinkType(%d)", t)
	}
}

// Header contains information from the pcap stream header.
type Header struct {
	// SnapLen is the snapshot length used for capturing each packet.
	// The Packet.Data length must not exceed this value.
	SnapLen int

	// LinkType specifies the type of each Packet.Data.
	LinkType LinkType
}

// Packet describes a captured packet.
type Packet struct {
	// Timestamp is the timestamp of the packet.
	Timestamp time.Time

	// OrigLen is the original packet length.
	// Due to a smaller Header.SnapLen, the Data may be truncated.
	OrigLen int

	// Data is the captured contents of the packet,
	// where the type is determined by Header.LinkType.
	// Its length may be less than OrigLen if truncated and
	// must not exceed Header.SnapLen.
	Data []byte
}

// Clone returns a deep copy of Packet.
func (p Packet) Clone() Packet {
	p.Data = append([]byte(nil), p.Data...) // TODO: Use bytes.Clone.
	return p
}

// Writer implements a streaming pcap writer.
type Writer struct {
	header Header
	writer io.Writer

	scratch [packetHeaderSize]byte

	err error
}

// NewWriter writes the pcap header and returns a Writer that
// writes each subsequent packet after the header.
//
// For performance, it is recommended that the io.Writer be buffered,
// such as provided by bufio.Writer.
func NewWriter(wr io.Writer, h Header) (*Writer, error) {
	var header [globalHeaderSize]byte
	binary.LittleEndian.PutUint32(header[0:4], magicNanos)           // magicNumber
	binary.LittleEndian.PutUint16(header[4:6], 2)                    // versionMinor
	binary.LittleEndian.PutUint16(header[6:8], 4)                    // versionMajor
	binary.LittleEndian.PutUint32(header[16:20], uint32(h.SnapLen))  // snapLen
	binary.LittleEndian.PutUint32(header[20:24], uint32(h.LinkType)) // linkType
	if _, err := wr.Write(header[:]); err != nil {
		return nil, err
	}
	return &Writer{header: h, writer: wr}, nil
}

// WriteNext writes the next packet in the stream.
// The Packet.Data is written immediately and not retained by the call.
func (w *Writer) WriteNext(p Packet) error {
	if w.err != nil {
		return w.err
	}
	if len(p.Data) > p.OrigLen {
		return fmt.Errorf("capture length exceeds packet length: %d > %d", len(p.Data), p.OrigLen)
	}
	if len(p.Data) > w.header.SnapLen {
		return fmt.Errorf("capture length exceeds snapshot length: %d > %d", len(p.Data), w.header.SnapLen)
	}
	var b []byte
	if wb, _ := w.writer.(*bufio.Writer); wb != nil {
		b = wb.AvailableBuffer()
	} else {
		b = w.scratch[:0]
	}
	b = binary.LittleEndian.AppendUint32(b, uint32(p.Timestamp.Unix()))
	b = binary.LittleEndian.AppendUint32(b, uint32(p.Timestamp.Nanosecond()))
	b = binary.LittleEndian.AppendUint32(b, uint32(len(p.Data)))
	b = binary.LittleEndian.AppendUint32(b, uint32(p.OrigLen))
	if _, w.err = w.writer.Write(b); w.err != nil {
		return w.err
	}
	if _, w.err = w.writer.Write(p.Data); w.err != nil {
		return w.err
	}
	return nil
}

// Reader implements a streaming pcap reader.
type Reader struct {
	Header // must not be modified

	reader    io.Reader
	bufReader *bufio.Reader
	nanos     bool
	swapped   bool
	location  *time.Location
}

// NewReader parses the pcap header and returns a Reader that
// parses each subsequent packet after the header.
//
// The first call to Reader.ReadNext allocates a buffer
// proportional to Header.SnapLen, which may be arbitrarily large.
// It is the caller's responsibility to check whether the amount of memory
// needed is reasonable.
func NewReader(rd io.Reader) (*Reader, error) {
	var header [globalHeaderSize]byte
	if _, err := io.ReadFull(rd, header[:]); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	var (
		magicNumber  = binary.LittleEndian.Uint32(header[0:4])
		versionMajor = binary.LittleEndian.Uint16(header[4:6])
		versionMinor = binary.LittleEndian.Uint16(header[6:8])
		zoneOffset   = binary.LittleEndian.Uint32(header[8:12])
		snapLen      = binary.LittleEndian.Uint32(header[16:20])
		linkType     = binary.LittleEndian.Uint32(header[20:24])
	)

	r := new(Reader)
	switch magicNumber {
	case magicMicros:
		r.nanos, r.swapped = false, false
	case magicNanos:
		r.nanos, r.swapped = true, false
	case bits.ReverseBytes32(magicMicros):
		r.nanos, r.swapped = false, true
	case bits.ReverseBytes32(magicNanos):
		r.nanos, r.swapped = true, true
	default:
		return nil, fmt.Errorf("unknown magic number 0x%08x", magicNumber)
	}
	if r.swapped {
		versionMajor = bits.ReverseBytes16(versionMajor)
		versionMinor = bits.ReverseBytes16(versionMinor)
		zoneOffset = bits.ReverseBytes32(zoneOffset)
		snapLen = bits.ReverseBytes32(snapLen)
		linkType = bits.ReverseBytes32(linkType)
	}
	if versionMajor != 2 && versionMinor != 4 {
		return nil, fmt.Errorf("unsupported version %d.%d", versionMajor, versionMinor)
	}
	if zoneOffset != 0 {
		r.location = time.FixedZone("", -int(int32(zoneOffset)))
	}
	if int32(snapLen) < 0 {
		return nil, fmt.Errorf("snapshot length overflows 32-bit signed integer")
	}
	r.SnapLen = int(snapLen)
	r.LinkType = LinkType(linkType)
	r.reader = rd
	return r, nil
}

// ReadNext reads the next packet in the stream.
// The contents of Packet.Data is only valid until the next ReadNext call.
// Call Packet.Clone if a copy of the packet is needed.
// If there are no more packets in the stream, it returns io.EOF.
func (r *Reader) ReadNext() (Packet, error) {
	// Initialize the bufio.Reader if necessary.
	if r.bufReader == nil {
		bufSize := 2 * (packetHeaderSize + r.SnapLen)
		if bufSize < defaultBufSize {
			bufSize = defaultBufSize
		}
		r.bufReader = bufio.NewReaderSize(r.reader, bufSize)
	}

	// Read and parse the packet header.
	header, err := r.bufReader.Peek(packetHeaderSize)
	if err != nil {
		if err == io.EOF && len(header) > 0 {
			err = io.ErrUnexpectedEOF
		}
		return Packet{}, err
	}
	var (
		timeSec = binary.LittleEndian.Uint32(header[0:4])
		timeSub = binary.LittleEndian.Uint32(header[4:8])
		capLen  = binary.LittleEndian.Uint32(header[8:12])
		origLen = binary.LittleEndian.Uint32(header[12:16])
	)
	if r.swapped {
		timeSec = bits.ReverseBytes32(timeSec)
		timeSub = bits.ReverseBytes32(timeSub)
		capLen = bits.ReverseBytes32(capLen)
		origLen = bits.ReverseBytes32(origLen)
	}
	var t time.Time
	if r.nanos {
		t = time.Unix(int64(timeSec), int64(timeSub))
	} else {
		t = time.Unix(int64(timeSec), 1000*int64(timeSub))
	}
	if r.location != nil {
		t = t.In(r.location)
	}
	if capLen > origLen {
		return Packet{}, fmt.Errorf("capture length exceeds packet length: %d > %d", capLen, origLen)
	}
	if capLen > uint32(r.SnapLen) {
		return Packet{}, fmt.Errorf("capture length exceeds snapshot length: %d > %d", capLen, r.SnapLen)
	}

	// Read the packet data.
	b, err := r.bufReader.Peek(packetHeaderSize + int(capLen))
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return Packet{}, err
	}

	r.bufReader.Discard(packetHeaderSize + int(capLen))
	return Packet{t, int(origLen), b[packetHeaderSize:]}, nil
}
