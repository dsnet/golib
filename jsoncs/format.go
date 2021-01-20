// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package jsoncs implements of the JSON Canonicalization Scheme (JCS)
// as specified in RFC 8785.
package jsoncs

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/bits"
	"sort"
	"strconv"
	"unicode/utf8"
)

// Format transforms the JSON input to its canonical form.
// The input must comply with RFC 7493.
// It reuses the provided input buffer.
func Format(b []byte) ([]byte, error) {
	if Valid(b) {
		return b, nil
	}
	v, err := decode(b)
	if err != nil {
		return nil, err
	}
	return formatValue(b[:0], v)
}

// formatValue canonically marshals a JSON value.
func formatValue(b []byte, v jsonValue) ([]byte, error) {
	switch v := v.(type) {
	case jsonObject:
		return formatObject(b, v)
	case jsonArray:
		return formatArray(b, v)
	case string:
		return formatString(b, v)
	case float64:
		return formatNumber(b, v)
	case nil:
		return append(b, nullLiteral...), nil
	case bool:
		switch v {
		case true:
			return append(b, trueLiteral...), nil
		case false:
			return append(b, falseLiteral...), nil
		}
	}
	return nil, fmt.Errorf("invalid type: %T", v)
}

// formatObject canonically marshals a JSON object per RFC 8785, section 3.2.3.
func formatObject(b []byte, obj jsonObject) ([]byte, error) {
	var ks []string
	for k := range obj {
		ks = append(ks, k)
	}
	sort.Slice(ks, func(i, j int) bool {
		return lessUTF16(ks[i], ks[j])
	})

	var err error
	b = append(b, '{')
	for _, k := range ks {
		b, err = formatString(b, k)
		if err != nil {
			return nil, err
		}
		b = append(b, ':')
		b, err = formatValue(b, obj[k])
		if err != nil {
			return nil, err
		}
		b = append(b, ',')
	}
	b = bytes.TrimRight(b, ",")
	b = append(b, '}')
	return b, nil
}

// formatArray canonically marshals a JSON array.
func formatArray(b []byte, arr jsonArray) ([]byte, error) {
	var err error
	b = append(b, '[')
	for _, v := range arr {
		b, err = formatValue(b, v)
		if err != nil {
			return nil, err
		}
		b = append(b, ',')
	}
	b = bytes.TrimRight(b, ",")
	b = append(b, ']')
	return b, nil
}

// formatArray canonically marshals a JSON string per RFC 8785, section 3.2.2.2.
func formatString(b []byte, s string) ([]byte, error) {
	// indexNeedEscape returns the index of the character that needs escaping.
	indexNeedEscape := func(s string) int {
		for i, r := range s {
			if r < ' ' || r == '\\' || r == '"' || r == utf8.RuneError {
				return i
			}
		}
		return len(s)
	}

	b = append(b, '"')
	i := indexNeedEscape(s)
	s, b = s[i:], append(b, s[:i]...)
	for len(s) > 0 {
		switch r, n := utf8.DecodeRuneInString(s); {
		case r == utf8.RuneError && n == 1:
			return nil, errors.New("invalid UTF-8 in string")
		case r < ' ' || r == '"' || r == '\\':
			b = append(b, '\\')
			switch r {
			case '"', '\\':
				b = append(b, byte(r))
			case '\b':
				b = append(b, 'b')
			case '\f':
				b = append(b, 'f')
			case '\n':
				b = append(b, 'n')
			case '\r':
				b = append(b, 'r')
			case '\t':
				b = append(b, 't')
			default:
				b = append(b, 'u')
				b = append(b, "0000"[1+(bits.Len32(uint32(r))-1)/4:]...)
				b = strconv.AppendUint(b, uint64(r), 16)
			}
			s = s[n:]
		default:
			i := indexNeedEscape(s[n:])
			s, b = s[n+i:], append(b, s[:n+i]...)
		}
	}
	b = append(b, '"')
	return b, nil
}

// formatNumber canonically marshals a JSON number per RFC 8785, section 3.2.2.3.
func formatNumber(b []byte, f float64) ([]byte, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return nil, fmt.Errorf("invalid float value: %v", f)
	}
	return formatES6(b, f), nil
}
