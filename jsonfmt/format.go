// Copyright 2017, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsonfmt

import (
	"bytes"
	"unicode"
)

func (s *state) format() (out []byte) {
	defer func() {
		if ex := recover(); ex != nil {
			if js, ok := ex.(jsonInvalid); ok {
				s.flushAppend(js...)
			} else {
				panic(ex)
			}
		}
		if s.hasNewlines {
			s.out = append(s.out, '\n')
		}
		out = s.out
	}()

	s.out = nil
	s.formatMeta(s.preVal)
	s.formatValue(s.val)
	s.formatMeta(s.postVal)
	return
}

func (s *state) formatValue(js jsonValue) {
	switch v := js.(type) {
	case *jsonObject:
		s.formatObject(v)
	case *jsonArray:
		s.formatArray(v)
	case jsonString:
		s.flushAppend(v...)
	case jsonNumber:
		s.flushAppend(v...)
	case jsonLiteral:
		s.flushAppend(v...)
	case jsonInvalid:
		panic(v)
	}
}

func (s *state) formatObject(js *jsonObject) {
	s.flushAppend('{')
	indentOuter := hasNewlines(js.preRecords)
	if indentOuter {
		s.pushIndent()
	}
	s.formatMeta(js.preRecords)
	for i, rec := range js.records {
		s.formatMeta(rec.preKey)
		s.formatValue(rec.key)
		s.formatMeta(rec.postKey)
		s.flushAppend(':')
		indentInner := hasNewlines(rec.preVal)
		if indentInner {
			s.pushIndent()
		}
		s.formatMeta(rec.preVal)
		s.formatValue(rec.val)
		s.formatMeta(rec.postVal)
		if i < len(js.records)-1 || s.emitTrailingComma(rec.postComma, js.postRecords) {
			s.flushAppend(',')
		}
		if indentInner {
			s.popIndent()
		}
		s.formatMeta(rec.postComma)
	}
	s.formatMeta(js.postRecords)
	if indentOuter {
		s.popIndent()
	}
	s.flushAppend('}')
}

func (s *state) formatArray(js *jsonArray) {
	s.flushAppend('[')
	indentOuter := hasNewlines(js.preElems)
	if indentOuter {
		s.pushIndent()
	}
	s.formatMeta(js.preElems)
	for i, rec := range js.elems {
		s.formatMeta(rec.preVal)
		s.formatValue(rec.val)
		s.formatMeta(rec.postVal)
		if i < len(js.elems)-1 || s.emitTrailingComma(rec.postComma, js.postElems) {
			s.flushAppend(',')
		}
		s.formatMeta(rec.postComma)
	}
	s.formatMeta(js.postElems)
	if indentOuter {
		s.popIndent()
	}
	s.flushAppend(']')
}

func (s *state) formatMeta(js jsonMeta) {
	for _, m := range js {
		switch m := m.(type) {
		case jsonComment:
			s.flushSpaces('/')
			s.formatComment(m)
		case jsonNewlines:
			for i := 0; i < int(m); i++ {
				s.newlines = append(s.newlines, '\n')
			}
		case jsonInvalid:
			panic(m)
		}
	}
}

func (s *state) formatComment(js jsonComment) {
	if s.standardize {
		return // Comments are not valid in ECMA-404
	}
	bs := bytes.Split(js, newlineBytes)
	prefix := indentPrefix(bs[len(bs)-1])
	allStars := len(bs) > 1
	for i, b := range bs {
		bs[i] = bytes.TrimRightFunc(bytes.TrimPrefix(b, prefix), unicode.IsSpace)
		if i > 0 {
			allStars = allStars && len(bs[i]) > 0 && bs[i][0] == '*'
		}
	}
	if allStars { // Line up stars in block comments
		for i, b := range bs {
			if i > 0 {
				bs[i] = append(spaceBytes, b...)
			}
		}
	}
	js = bytes.Join(bs, append(newlineBytes, s.indents...))
	s.hasNewlines = s.hasNewlines || len(bs) > 1
	s.out = append(s.out, js...)
}

// flushAppend calls flushSpaces before appending b to the output.
func (s *state) flushAppend(b ...byte) {
	if len(b) > 0 {
		s.flushSpaces(b[0])
		s.out = append(s.out, b...)
	}
}

// flushSpaces determines how many spaces and newlines to output using
// information from the previous and next non-whitespace characters.
func (s *state) flushSpaces(next byte) {
	if s.minify {
		return
	}
	var prev byte
	if len(s.out) > 0 {
		prev = s.out[len(s.out)-1]
	} else {
		s.newlines = s.newlines[:0] // Avoid leading empty lines
	}
	if len(s.newlines) > 2 {
		s.newlines = s.newlines[:2] // Avoid more than 1 empty line
	}
	if len(s.newlines) > 1 && (prev == '{' || prev == '[' || next == '}' || next == ']') {
		s.newlines = s.newlines[:1] // Avoid empty lines after open brace or before closing brace
	}
	if len(s.newlines) > 1 && prev == ':' {
		s.newlines = s.newlines[:1] // Avoid empty lines after lines ending with colon
	}
	if next == ':' || next == ',' {
		s.newlines = s.newlines[:0] // Avoid starting lines with a colon or comma
	}
	if (prev == '{' && next == '}') || (prev == '[' && next == ']') {
		s.newlines = s.newlines[:0] // Always collapse empty objects and arrays
	}
	if len(s.newlines) > 0 {
		s.hasNewlines = true
		s.out = append(s.out, s.newlines...)
		s.out = append(s.out, s.indents...)
		s.newlines = s.newlines[:0]
	} else if prev > 0 && (prev == ':' || prev == ',' || prev == '/' || next == '/') {
		s.out = append(s.out, ' ')
	}
}

// emitTrailingComma reports whether a trailing comma should be emitted.
// Only emit trailing commas if the source had trailing commas,
// and there is at least one newline until the closing brace.
func (s *state) emitTrailingComma(postComma, postVals jsonMeta) bool {
	if s.trailingComma && !s.standardize {
		return hasNewlines(postComma) || hasNewlines(postVals)
	}
	return false
}

func hasNewlines(js jsonMeta) bool {
	for _, m := range js {
		switch m := m.(type) {
		case jsonComment:
			if bytes.IndexByte(m, '\n') > 0 {
				return true
			}
		case jsonNewlines:
			if m > 0 {
				return true
			}
		}
	}
	return false
}

func indentPrefix(s []byte) []byte {
	n := len(s) - len(bytes.TrimLeft(s, " \t"))
	return s[:n]
}
