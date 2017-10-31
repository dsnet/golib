// Copyright 2017, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsonfmt

import (
	"bytes"
	"fmt"
	"reflect"
)

// expandAST walks js and inserts newlines into jsonObject and jsonArray nodes
// if the non-whitespace of any line exceeds the column limit.
func expandAST(js jsonValue, columnLimit int) {
	m := make(map[int][]jsonValue)
	for d := exploreAST(js, m, 0); d >= 0; d-- {
		for _, js := range m[d] {
			switch js := js.(type) {
			case *jsonObject:
				tryExpandObject(js, columnLimit)
			case *jsonArray:
				tryExpandArray(js, columnLimit)
			}
		}
	}
}

// exploreAST walks the AST of js and inserts any jsonObjects or jsonArrays
// encountered into map m at the depth encountered.
// It returns the maximum depth encountered.
func exploreAST(js jsonValue, m map[int][]jsonValue, depth int) (maxDepth int) {
	maxDepth = depth
	switch v := js.(type) {
	case *jsonObject:
		m[depth] = append(m[depth], v)
		for _, rec := range v.records {
			if d := exploreAST(rec.val, m, depth+1); d > maxDepth {
				maxDepth = d
			}
		}
	case *jsonArray:
		m[depth] = append(m[depth], v)
		for _, elem := range v.elems {
			if d := exploreAST(elem.val, m, depth+1); d > maxDepth {
				maxDepth = d
			}
		}
	}
	return maxDepth
}

// tryExpandObject determines whether to expand jsonObject, and expands if so.
func tryExpandObject(js *jsonObject, columnLimit int) {
	for i := range js.records {
		prevLen, multi1 := lineLength(-1, js.records[:i])
		nextLen, multi2 := lineLength(+1, js.records[i:])
		expandMulti := len(js.records) > 1 && (multi1 || multi2)
		if prevLen+nextLen > columnLimit || expandMulti {
			expandObject(js)
			return
		}
	}
}
func expandObject(js *jsonObject) {
	js.preRecords = appendNewline(js.preRecords)
	for i, rec := range js.records {
		rec.postComma = appendNewline(rec.postComma)
		js.records[i] = rec
	}
	js.postRecords = appendNewline(js.postRecords)
}

// tryExpandArray determines whether to expand jsonArray, and expands if so.
func tryExpandArray(js *jsonArray, columnLimit int) {
	for i := range js.elems {
		prevLen, _ := lineLength(-1, js.elems[:i])
		nextLen, _ := lineLength(+1, js.elems[i:])
		if prevLen+nextLen > columnLimit {
			expandArray(js)
			return
		}
	}
}
func expandArray(js *jsonArray) {
	js.preElems = appendNewline(js.preElems)
	for i, elem := range js.elems {
		elem.postComma = appendNewline(elem.postComma)
		js.elems[i] = elem
	}
	js.postElems = appendNewline(js.postElems)
}

// lineLength reports the upcoming line length in the sequence of AST nodes.
// It reports the length, and whether a newline was encountered.
// The direction dir may only be +1 or -1 to control walking in a forward or
// reverse direction.
func lineLength(dir int, vs ...interface{}) (n int, multi bool) {
	type token byte
	length := func(v interface{}) (int, bool) {
		switch v := v.(type) {
		case *jsonObject:
			return lineLength(dir, token('{'), v.preRecords, v.records, v.postRecords, token('}'))
		case *jsonArray:
			return lineLength(dir, token('['), v.preElems, v.elems, v.postElems, token(']'))
		case jsonRecord:
			return lineLength(dir, v.preKey, v.key, v.postKey, token(':'), v.preVal, v.val, v.postVal, token(','), v.postComma)
		case jsonElement:
			return lineLength(dir, v.preVal, v.val, v.postVal, token(','), v.postComma)
		case []jsonRecord, []jsonElement, jsonMeta:
			var args []interface{}
			vv := reflect.ValueOf(v)
			for j := 0; j < vv.Len(); j++ {
				args = append(args, vv.Index(j).Interface())
			}
			return lineLength(dir, args...)
		case jsonString, jsonNumber, jsonLiteral:
			return reflect.ValueOf(v).Len(), false
		case jsonComment, jsonInvalid:
			b := reflect.ValueOf(v).Bytes()
			switch dir {
			case +1:
				if i := bytes.IndexByte(b, '\n'); i >= 0 {
					return i, true
				}
			case -1:
				if i := bytes.LastIndexByte(b, '\n'); i >= 0 {
					return len(b) - (i + 1), true
				}
			}
			return len(b), false
		case jsonNewlines:
			return 0, v > 0
		case token:
			return 1, false
		case nil:
			return 0, false
		default:
			panic(fmt.Sprintf("unable to handle type %T", v))
		}
	}

	var m int
	switch dir {
	case +1:
		for i := 0; i < len(vs) && !multi; i++ {
			m, multi = length(vs[i])
			n += m
		}
	case -1:
		for i := len(vs) - 1; i >= 0 && !multi; i-- {
			m, multi = length(vs[i])
			n += m
		}
	default:
		panic("invalid direction")
	}
	return n, multi
}

func appendNewline(js jsonMeta) jsonMeta {
	if len(js) > 0 {
		if nl, _ := js[len(js)-1].(jsonNewlines); nl > 0 {
			return js
		}
	}
	return append(js, jsonNewlines(1))
}
