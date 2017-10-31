// Copyright 2017, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsonfmt

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"unicode/utf8"
)

var (
	stringRx  = `"(\\(["\\\/bfnrt]|u[a-fA-F0-9]{4})|[^"\\\x00-\x1f\x7f]+)*"`
	numberRx  = `-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?`
	literalRx = `(true|false|null)`

	commentRx = `(/\*([^\n]|\n)*?\*/|//[^\n]*\n?)`
	spaceRx   = `[ \r\n\t]*`

	stringRegex  = regexp.MustCompile("^" + stringRx)
	numberRegex  = regexp.MustCompile("^" + numberRx)
	literalRegex = regexp.MustCompile("^" + literalRx)

	commentRegex = regexp.MustCompile("^" + commentRx)
	spaceRegex   = regexp.MustCompile("^" + spaceRx)
)

func (s *state) parse(in []byte) (err error) {
	defer func() {
		if ex := recover(); ex != nil {
			if je, ok := ex.(jsonError); ok {
				// Insert line/column information to the error message.
				parsed := in[:len(in)-len(s.in)]
				je.line = bytes.Count(parsed, newlineBytes) + 1
				if i := bytes.LastIndexByte(parsed, '\n'); i >= 0 {
					parsed = parsed[i+len("\n"):]
				}
				je.column = len(parsed) + 1
				err = je

				// Insert the remainder input to the last node.
				switch js := s.last.(type) {
				case *jsonValue:
					*js = jsonInvalid(s.in)
				case *jsonMeta:
					*js = append(*js, jsonInvalid(s.in))
				}
				s.in = nil
			} else {
				panic(ex)
			}
		}
	}()

	s.in = in
	s.parseMeta(&s.preVal)
	s.parseValue(&s.val)
	s.parseMeta(&s.postVal)
	if len(s.in) > 0 {
		panic(s.errorf("unexpected trailing input: %s", bytesPreview(s.in)))
	}
	return
}

func (s *state) parseValue(js *jsonValue) {
	s.last = js
	switch s.nextChar() {
	case '{':
		s.parseObject(js)
	case '[':
		s.parseArray(js)
	case '"':
		s.parseString(js)
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.parseNumber(js)
	case 't', 'f', 'n':
		s.parseLiteral(js)
	default:
		if len(s.in) == 0 {
			panic(s.errorf("unable to parse value: unexpected EOF"))
		}
		panic(s.errorf("unable to parse value: unexpected %q", s.in[0]))
	}
}

func (s *state) parseObject(js *jsonValue) {
	obj := new(jsonObject)
	*js = obj

	var trailingComma bool
	s.parseChar('{', "object")
	s.parseMeta(&obj.preRecords)
	for s.nextChar() != '}' {
		obj.records = append(obj.records, jsonRecord{})
		rec := &obj.records[len(obj.records)-1]

		s.parseString(&rec.key)
		s.parseMeta(&rec.postKey)
		s.parseChar(':', "object")
		s.parseMeta(&rec.preVal)
		s.parseValue(&rec.val)
		s.parseMeta(&rec.postVal)

		if s.nextChar() == '}' {
			rec.postVal, rec.postComma = nil, rec.postVal
			trailingComma = false
			break
		}
		s.parseChar(',', "object")
		s.parseMeta(&rec.postComma)
		trailingComma = true
	}
	s.parseChar('}', "object")
	s.trailingComma = s.trailingComma || trailingComma

	var meta jsonMeta
	obj.preRecords, meta = splitMeta(obj.preRecords)
	for i := range obj.records {
		obj.records[i].preKey = meta
		obj.records[i].postComma, meta = splitMeta(obj.records[i].postComma)
	}
	obj.postRecords = meta
}

func (s *state) parseArray(js *jsonValue) {
	arr := new(jsonArray)
	*js = arr

	var trailingComma bool
	s.parseChar('[', "array")
	s.parseMeta(&arr.preElems)
	for s.nextChar() != ']' {
		arr.elems = append(arr.elems, jsonElement{})
		elem := &arr.elems[len(arr.elems)-1]

		s.parseValue(&elem.val)
		s.parseMeta(&elem.postVal)

		if s.nextChar() == ']' {
			elem.postVal, elem.postComma = nil, elem.postVal
			trailingComma = false
			break
		}
		s.parseChar(',', "array")
		s.parseMeta(&elem.postComma)
		trailingComma = true
	}
	s.parseChar(']', "array")
	s.trailingComma = s.trailingComma || trailingComma

	var meta jsonMeta
	arr.preElems, meta = splitMeta(arr.preElems)
	for i := range arr.elems {
		arr.elems[i].preVal = meta
		arr.elems[i].postComma, meta = splitMeta(arr.elems[i].postComma)
	}
	arr.postElems = meta
}

func (s *state) parseString(js *jsonValue) {
	b := stringRegex.Find(s.in)
	if len(b) == 0 {
		panic(s.errorf("unable to parse string: %s", bytesPreview(s.in)))
	}
	*js, s.in = jsonString(recodeString(b)), s.in[len(b):]
}

func (s *state) parseNumber(js *jsonValue) {
	b := numberRegex.Find(s.in)
	if len(b) == 0 {
		panic(s.errorf("unable to parse number: %s", bytesPreview(s.in)))
	}
	*js, s.in = jsonNumber(recodeNumber(b)), s.in[len(b):]
}

func (s *state) parseLiteral(js *jsonValue) {
	b := literalRegex.Find(s.in)
	if len(b) == 0 {
		panic(s.errorf("unable to parse literal: %s", bytesPreview(s.in)))
	}
	*js, s.in = jsonLiteral(b), s.in[len(b):]
}

func (s *state) parseMeta(js *jsonMeta) {
	s.last = js
	for {
		if b := commentRegex.Find(s.in); len(b) > 0 {
			b = bytes.TrimRight(b, "\n")
			if !s.standardize {
				*js = append(*js, jsonComment(b))
			}
			s.in = s.in[len(b):]
			continue
		}
		if b := spaceRegex.Find(s.in); len(b) > 0 {
			if !s.minify {
				if n := bytes.Count(b, newlineBytes); n > 0 {
					if len(*js) == 0 {
						*js = append(*js, jsonNewlines(n))
					} else if nl, ok := (*js)[len(*js)-1].(jsonNewlines); ok {
						(*js)[len(*js)-1] = nl + jsonNewlines(n)
					} else {
						*js = append(*js, jsonNewlines(n))
					}
				}
			}
			s.in = s.in[len(b):]
			continue
		}
		break
	}
}

// parseChar parses the next character for want.
// If the next char is not want, this panics with a descriptive jsonError.
func (s *state) parseChar(want byte, what string) {
	if got := s.nextChar(); got != want {
		if len(s.in) == 0 {
			panic(s.errorf("unable to parse %s: unexpected EOF", what))
		}
		panic(s.errorf("unable to parse %s: got %q, expected %q", what, got, want))
	}
	s.in = s.in[1:]
}

// nextChat reports the next character in the input, returning 0 if EOF.
func (s *state) nextChar() byte {
	if len(s.in) == 0 {
		return 0
	}
	return s.in[0]
}

func (s *state) errorf(f string, x ...interface{}) error {
	return jsonError{message: fmt.Sprintf(f, x...)}
}

var hexLUT = [256]rune{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15,
}

// recodeString decodes a JSON string and re-encodes it such that it uses
// UTF-8 when possible instead of the \uXXXX notation.
// This assumes that the input has already been validated.
// This does not mutate the input, but may alias it.
func recodeString(in []byte) (out []byte) {
	// TODO: Support other string encoding options:
	//	* Verbatim: return input as is.
	//	* SafeJS: escapes U+2028 and U+2029 with the \u notation.
	//	* SafeHTML: escapes '<', '>', and '&' with the \u notation.
	//	* OnlyASCII: escapes all Unicode with \u notation.

	var rb [utf8.UTFMax]byte
	out = append(out, '"')
	in = in[1 : len(in)-1]
	for len(in) > 0 {
		switch c := in[0]; {
		case c == '\\':
			switch in[1] {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				out = append(out, in[:2]...) // Copy ad-verbatim
				in = in[2:]
			case 'u':
				r := (hexLUT[in[2]] << 12) + (hexLUT[in[3]] << 8) + (hexLUT[in[4]] << 4) + (hexLUT[in[5]] << 0)
				out = append(out, rb[:utf8.EncodeRune(rb[:], r)]...)
				in = in[6:]
			}
			continue
		case c < utf8.RuneSelf:
			out = append(out, c) // Copy ad-verbatim
			in = in[1:]
		default:
			r, n := utf8.DecodeRune(in) // Copy ad-verbatim except for rune errors
			out = append(out, rb[:utf8.EncodeRune(rb[:], r)]...)
			in = in[n:]
		}
	}
	return append(out, '"')
}

// recodeNumber re-encodes the number is a possibly shorter form.
// This does not mutate the input, but may alias it.
func recodeNumber(in []byte) (out []byte) {
	// TODO: Support other float encoding options:
	//	* Verbatim: return input as is.
	//	* Precision: control bit precision.

	f, err := strconv.ParseFloat(string(in), 64)
	if err != nil {
		return in
	}
	if abs := math.Abs(f); abs != 0 && abs < 1e-6 || abs >= 1e21 {
		out = strconv.AppendFloat(nil, f, 'e', -1, 64)
	} else {
		out = strconv.AppendFloat(nil, f, 'f', -1, 64)
	}
	if len(in) < len(out) {
		return in
	}
	return out
}

// splitMeta splits the input such that any meta nodes without newlines are
// split off as the second pair.
func splitMeta(js jsonMeta) (jsonMeta, jsonMeta) {
	for i := len(js) - 1; i >= 0; i-- {
		if _, ok := js[i].(jsonNewlines); ok {
			return js[:i+1], js[i+1:]
		}
	}
	return nil, js
}

func bytesPreview(b []byte) string {
	const prevLen = 8
	if len(b) > prevLen {
		return fmt.Sprintf("%q...", b[:prevLen])
	}
	return fmt.Sprintf("%q", b)
}
