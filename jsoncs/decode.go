// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsoncs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

type (
	jsonValue  = interface{}
	jsonArray  = []jsonValue
	jsonObject = map[string]jsonValue
)

var (
	nullLiteral  = []byte("null")
	trueLiteral  = []byte("true")
	falseLiteral = []byte("false")
)

// decode unmarshals the JSON input as a JSON value.
//
// It is similar to:
//
//	var v interface{}
//	err := json.Unmarshal(b, &v)
//	return v, err
//
// However, it is stricter than the standard library implementation in that
// it rejects JSON objects with duplicate keys and invalid UTF-8
// in adherence of RFC 7493.
func decode(b []byte) (jsonValue, error) {
	v, b, err := decodeValue(b)
	if err != nil {
		return nil, err
	}
	if len(trimSpace(b)) > 0 {
		return nil, errors.New("invalid data after top-level value")
	}
	return v, nil
}

// decodeValue unmarshals the next JSON value per RFC 7159, section 3.
// It consume the leading value, and returns the remaining bytes.
func decodeValue(b []byte) (jsonValue, []byte, error) {
	b = trimSpace(b)
	switch {
	case len(b) == 0:
		return nil, nil, io.ErrUnexpectedEOF
	case b[0] == '{':
		return decodeObject(b)
	case b[0] == '[':
		return decodeArray(b)
	case b[0] == '"':
		return decodeString(b)
	case (b[0] == '-' || ('0' <= b[0] && b[0] <= '9')):
		return decodeNumber(b)
	case bytes.HasPrefix(b, nullLiteral):
		return nil, b[len(nullLiteral):], nil
	case bytes.HasPrefix(b, trueLiteral):
		return true, b[len(trueLiteral):], nil
	case bytes.HasPrefix(b, falseLiteral):
		return false, b[len(falseLiteral):], nil
	default:
		return nil, b, errors.New("expected next JSON value")
	}
}

// decodeObject unmarshals the next JSON object per RFC 7159, section 4,
// with special attention paid to RFC 7493, section 2.3
// regarding rejection of duplicate entry names.
// It consume the leading value, and returns the remaining bytes.
func decodeObject(b []byte) (jsonObject, []byte, error) {
	b = trimSpace(b)
	switch {
	case len(b) == 0:
		return nil, b, io.ErrUnexpectedEOF
	case b[0] == '{':
		b = b[1:]
	default:
		return nil, b, errors.New("expected '{' character in JSON object")
	}

	var init bool
	var obj = make(jsonObject)
	var err error
	for {
		b = trimSpace(b)
		if len(b) > 0 && b[0] == '}' {
			return obj, b[1:], nil
		}

		if init {
			b = trimSpace(b)
			switch {
			case len(b) == 0:
				return nil, b, io.ErrUnexpectedEOF
			case b[0] == ',':
				b = b[1:]
			default:
				return nil, b, errors.New("expected ',' character in JSON object")
			}
		}

		var k string
		k, b, err = decodeString(b)
		if err != nil {
			return nil, b, err
		}
		if _, ok := obj[k]; ok {
			return nil, b, fmt.Errorf("duplicate key %q in JSON object", k)
		}

		b = trimSpace(b)
		switch {
		case len(b) == 0:
			return nil, b, io.ErrUnexpectedEOF
		case b[0] == ':':
			b = b[1:]
		default:
			return nil, b, errors.New("expected ':' character in JSON object")
		}

		var v jsonValue
		v, b, err = decodeValue(b)
		if err != nil {
			return nil, b, err
		}

		obj[k] = v
		init = true
	}
}

// decodeArray unmarshals the next JSON object per RFC 7159, section 5.
// It consume the leading value, and returns the remaining bytes.
func decodeArray(b []byte) (jsonArray, []byte, error) {
	b = trimSpace(b)
	switch {
	case len(b) == 0:
		return nil, b, io.ErrUnexpectedEOF
	case b[0] == '[':
		b = b[1:]
	default:
		return nil, b, errors.New("expected '[' character in JSON array")
	}

	var init bool
	var arr jsonArray
	var err error
	for {
		b = trimSpace(b)
		if len(b) > 0 && b[0] == ']' {
			return arr, b[1:], nil
		}

		if init {
			b = trimSpace(b)
			switch {
			case len(b) == 0:
				return nil, b, io.ErrUnexpectedEOF
			case b[0] == ',':
				b = b[1:]
			default:
				return nil, b, errors.New("expected ',' character in JSON array")
			}
		}

		var v jsonValue
		v, b, err = decodeValue(b)
		if err != nil {
			return nil, b, err
		}

		arr = append(arr, v)
		init = true
	}
}

// decodeString unmarshals the next JSON string per RFC 7159, section 7,
// with special attention paid to RFC 7493, section 2.1
// regarding rejection of unpaired surrogate halves.
// It consume the leading value, and returns the remaining bytes.
func decodeString(b []byte) (string, []byte, error) {
	// indexNeedEscape returns the index of the character that needs escaping.
	indexNeedEscape := func(b []byte) int {
		for i, r := range string(b) {
			if r < ' ' || r == '\\' || r == '"' || r == utf8.RuneError {
				return i
			}
		}
		return len(b)
	}

	b = trimSpace(b)
	switch {
	case len(b) == 0:
		return "", b, io.ErrUnexpectedEOF
	case b[0] == '"':
		b = b[1:]
	default:
		return "", b, errors.New("expected '\"' character in JSON string")
	}

	var s []byte
	i := indexNeedEscape(b)
	b, s = b[i:], append(s, b[:i]...)
	for len(b) > 0 {
		switch r, n := utf8.DecodeRune(b); {
		case r == utf8.RuneError && n == 1:
			return "", b, errors.New("invalid UTF-8 in JSON string")
		case r < ' ':
			return "", b, fmt.Errorf("invalid character %q in string", r)
		case r == '"':
			b = b[1:]
			return string(s), b, nil
		case r == '\\':
			if len(b) < 2 {
				return "", b, io.ErrUnexpectedEOF
			}
			switch r := b[1]; r {
			case '"', '\\', '/':
				b, s = b[2:], append(s, r)
			case 'b':
				b, s = b[2:], append(s, '\b')
			case 'f':
				b, s = b[2:], append(s, '\f')
			case 'n':
				b, s = b[2:], append(s, '\n')
			case 'r':
				b, s = b[2:], append(s, '\r')
			case 't':
				b, s = b[2:], append(s, '\t')
			case 'u':
				if len(b) < 6 {
					return "", b, io.ErrUnexpectedEOF
				}
				v, err := strconv.ParseUint(string(b[2:6]), 16, 16)
				if err != nil {
					return "", b, fmt.Errorf("invalid escape code %q in string", b[:6])
				}
				b = b[6:]

				r := rune(v)
				if utf16.IsSurrogate(r) {
					if len(b) < 6 {
						return "", b, io.ErrUnexpectedEOF
					}
					v, err := strconv.ParseUint(string(b[2:6]), 16, 16)
					r = utf16.DecodeRune(r, rune(v))
					if b[0] != '\\' || b[1] != 'u' || r == unicode.ReplacementChar || err != nil {
						return "", b, fmt.Errorf("invalid escape code %q in string", b[:6])
					}
					b = b[6:]
				}
				s = append(s, string(r)...)
			default:
				return "", b, fmt.Errorf("invalid escape code %q in string", b[:2])
			}
		default:
			i := indexNeedEscape(b[n:])
			b, s = b[n+i:], append(s, b[:n+i]...)
		}
	}
	return "", b, io.ErrUnexpectedEOF
}

// decodeNumber unmarshals the next JSON number per RFC 7159, section 6.
// It consume the leading value, and returns the remaining bytes.
func decodeNumber(b []byte) (float64, []byte, error) {
	b = trimSpace(b)
	b0 := b
	if len(b) > 0 && b[0] == '-' {
		b = b[1:]
	}
	switch {
	case len(b) == 0:
		return 0, b, io.ErrUnexpectedEOF
	case b[0] == '0':
		b = b[1:]
	case '1' <= b[0] && b[0] <= '9':
		b = b[1:]
		for len(b) > 0 && ('0' <= b[0] && b[0] <= '9') {
			b = b[1:]
		}
	default:
		return 0, nil, errors.New("expected digit character in JSON number")
	}
	if len(b) > 0 && b[0] == '.' {
		b = b[1:]
		switch {
		case len(b) == 0:
			return 0, b, io.ErrUnexpectedEOF
		case '0' <= b[0] && b[0] <= '9':
			b = b[1:]
		default:
			return 0, nil, errors.New("expected digit character in JSON number")
		}
		for len(b) > 0 && ('0' <= b[0] && b[0] <= '9') {
			b = b[1:]
		}
	}
	if len(b) > 0 && (b[0] == 'e' || b[0] == 'E') {
		b = b[1:]
		if len(b) > 0 && (b[0] == '-' || b[0] == '+') {
			b = b[1:]
		}
		switch {
		case len(b) == 0:
			return 0, b, io.ErrUnexpectedEOF
		case '0' <= b[0] && b[0] <= '9':
			b = b[1:]
		default:
			return 0, nil, errors.New("expected digit character in JSON number")
		}
		for len(b) > 0 && ('0' <= b[0] && b[0] <= '9') {
			b = b[1:]
		}
	}

	f, err := strconv.ParseFloat(string(b0[:len(b0)-len(b)]), 64)
	if err != nil {
		return 0, b, fmt.Errorf("invalid JSON number: %s", b0[:len(b0)-len(b)])
	}
	return f, b, nil
}

// trimSpace strips leading whitespace per RFC 7159, section 2.
func trimSpace(b []byte) []byte {
	return bytes.TrimLeft(b, " \t\n\r")
}
