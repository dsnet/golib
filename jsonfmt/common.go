// Copyright 2017, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// Package jsonfmt provides functionality for formatting JSON.
package jsonfmt

import "fmt"

var (
	newlineBytes = []byte{'\n'}
	spaceBytes   = []byte{' '}
)

type (
	jsonValue interface {
		isValue() // Satisfied by (*jsonObject | *jsonArray | jsonString | jsonNumber | jsonLiteral | jsonInvalid)
	}

	jsonObject struct {
		// '{'
		preRecords  jsonMeta // If non-empty, ends with jsonNewlines
		records     []jsonRecord
		postRecords jsonMeta // Never contains jsonNewlines
		// '}'
	}
	jsonRecord struct {
		preKey  jsonMeta // Never contains jsonNewlines
		key     jsonValue
		postKey jsonMeta
		// ':'
		preVal  jsonMeta
		val     jsonValue
		postVal jsonMeta
		// ','
		postComma jsonMeta // If non-empty, ends with jsonNewlines
	}

	jsonArray struct {
		// '['
		preElems  jsonMeta // If non-empty, ends with jsonNewlines
		elems     []jsonElement
		postElems jsonMeta // Never contains jsonNewlines
		// ']'
	}
	jsonElement struct {
		preVal  jsonMeta // Never contains jsonNewlines
		val     jsonValue
		postVal jsonMeta
		// ','
		postComma jsonMeta // If non-empty, ends with jsonNewlines
	}

	jsonString  []byte // Quoted string
	jsonNumber  []byte // Numeric value
	jsonLiteral []byte // "true" | "false" | "null"

	jsonMeta []interface {
		isMeta() // Implemented by (jsonComment | jsonNewlines | jsonInvalid)
	}
	jsonComment  []byte // Comment of either "//" or "/**/" form without trailing newlines
	jsonNewlines int    // Number of newlines

	// When a parsing error occurs, then the remainder of the input is stored
	// as jsonInvalid in the AST.
	jsonInvalid []byte // May possibly be an empty string
)

func (*jsonObject) isValue() {}
func (*jsonArray) isValue()  {}
func (jsonString) isValue()  {}
func (jsonNumber) isValue()  {}
func (jsonLiteral) isValue() {}
func (jsonInvalid) isValue() {}

func (jsonComment) isMeta()  {}
func (jsonNewlines) isMeta() {}
func (jsonInvalid) isMeta()  {}

// Option configures how to format JSON.
type Option interface {
	option()
}

type (
	minify      struct{ Option }
	standardize struct{ Option }
)

// TODO: Make these an user Option?
const defaultColumnLimit = 80
const defaultAlignLimit = 20

// Minify configures Format to produce the minimal representation of the input.
// If Format returns no error, then the output is guaranteed to be valid JSON,
func Minify() Option { return minify{} }

// Standardize configures Format to produce valid JSON according to ECMA-404.
// This strips any comments and trailing commas.
func Standardize() Option { return standardize{} }

// Format parses and formats the input JSON according to provided Options.
// If err is non-nil, then the output is a best effort at processing the input.
//
// This function accepts a superset of the JSON specification that allows
// comments and trailing commas after the last element in an object or array.
func Format(in []byte, opts ...Option) (out []byte, err error) {
	// Process the provided options.
	var st state
	for _, opt := range opts {
		switch opt.(type) {
		case minify:
			st.minify = true
			st.standardize = true
		case standardize:
			st.standardize = true
		default:
			panic(fmt.Sprintf("unknown option: %#v", opt))
		}
	}

	// Attempt to parse and format the JSON data.
	err = st.parse(in)
	if !st.minify {
		expandAST(st.val, defaultColumnLimit)
	}
	out = st.format()
	if !st.minify {
		out = alignJSON(out, defaultAlignLimit)
	}
	return out, err
}

type state struct {
	in      []byte
	preVal  jsonMeta
	val     jsonValue
	postVal jsonMeta
	last    interface{} // T where T = (*jsonValue | *jsonMeta)
	out     []byte

	// Parsing and formatting options.
	minify        bool // If set, implies standardize is set too
	standardize   bool // If set, output will be ECMA-404 compliant
	trailingComma bool // Set by parser if any trailing commas detected
	hasNewlines   bool // Set by formatter if any newlines are emitted

	newlines []byte // Pending newlines to output
	indents  []byte // Indents to output per line
}

func (s *state) pushIndent() {
	s.indents = append(s.indents, '\t')
}
func (s *state) popIndent() {
	s.indents = s.indents[:len(s.indents)-1]
}

type jsonError struct {
	line, column int
	message      string
}

func (e jsonError) Error() string {
	if e.line > 0 && e.column > 0 {
		return "jsonfmt: " + e.message
	}
	return fmt.Sprintf("jsonfmt: line %d, column %d: %v", e.line, e.column, e.message)
}
