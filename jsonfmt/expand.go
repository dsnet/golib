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
// if the line (excluding indentation) exceeds the column limit.
func expandAST(js jsonValue, limit int) {
	m := make(map[int][]jsonValue)
	for d := exploreAST(js, m, 0); d >= 0; d-- {
		for _, js := range m[d] {
			switch js := js.(type) {
			case *jsonObject:
				tryExpandObject(js, limit)
			case *jsonArray:
				tryExpandArray(js, limit)
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
func tryExpandObject(js *jsonObject, limit int) {
	if len(js.records) == 0 {
		return
	}

	// Only allow batching if newlines already exist in record meta nodes
	// and at least one record entirely lacks newlines.
	_, mayBatch := lineLength(+1, js.preRecords)
	hasBatch := false
	for _, rec := range js.records {
		_, multi := lineLength(+1, rec.postKey, rec.preVal, rec.postVal, rec.postComma)
		mayBatch = mayBatch || multi
		hasBatch = hasBatch || !multi
	}
	batch := mayBatch && hasBatch

	for i := range js.records {
		prevLen, multi1 := lineLength(-1, js.records[:i])
		nextLen, multi2 := lineLength(+1, js.records[i:])
		expandMulti := !batch && len(js.records) > 1 && (multi1 || multi2)
		if prevLen+nextLen > limit || expandMulti {
			expandObject(js, batch, limit)
			return
		}
	}
}
func expandObject(js *jsonObject, batch bool, limit int) {
	n := len(js.records)
	js.preRecords = appendNewline(js.preRecords)
	for i := 1; i < n; i++ {
		prevLen, _ := lineLength(-1, js.records[:i])
		nextLen, _ := lineLength(+1, js.records[i:][:1])
		if !batch || prevLen+nextLen > limit {
			js.records[i-1].postComma = appendNewline(js.records[i-1].postComma)
		}
	}
	js.records[n-1].postComma = appendNewline(js.records[n-1].postComma)
}

// tryExpandArray determines whether to expand jsonArray, and expands if so.
func tryExpandArray(js *jsonArray, limit int) {
	if len(js.elems) == 0 {
		return
	}

	// Only allow batching if newlines already exist in element meta nodes
	// and at least one element entirely lacks newlines.
	_, mayBatch := lineLength(+1, js.preElems)
	hasBatch := false
	for _, elem := range js.elems {
		_, multi := lineLength(+1, elem.postVal, elem.postComma)
		mayBatch = mayBatch || multi
		hasBatch = hasBatch || !multi
	}
	batch := mayBatch && hasBatch

	for i := range js.elems {
		prevLen, _ := lineLength(-1, js.elems[:i])
		nextLen, _ := lineLength(+1, js.elems[i:])
		if prevLen+nextLen > limit {
			expandArray(js, batch, limit)
			return
		}
	}
}
func expandArray(js *jsonArray, batch bool, limit int) {
	n := len(js.elems)
	js.preElems = appendNewline(js.preElems)
	for i := 1; i < n; i++ {
		// Always batch primitive values together.
		isPrim := isPrimitive(js.elems[i-1].val) && isPrimitive(js.elems[i].val)
		prevLen, _ := lineLength(-1, js.elems[:i])
		nextLen, _ := lineLength(+1, js.elems[i:][:1])
		if !(batch || isPrim) || (prevLen+nextLen) > limit {
			js.elems[i-1].postComma = appendNewline(js.elems[i-1].postComma)
		}
	}
	js.elems[n-1].postComma = appendNewline(js.elems[n-1].postComma)
}

// lineLength reports the upcoming line length in the sequence of AST nodes.
// It reports the length, and whether a newline was encountered.
// The direction dir may only be +1 or -1 to control walking in a forward or
// reverse direction.
func lineLength(dir int, vs ...interface{}) (n int, multi bool) {
	length := func(v interface{}) (int, bool) {
		switch v := v.(type) {
		case *jsonObject:
			return lineLength(dir, '{', v.preRecords, v.records, v.postRecords, '}')
		case *jsonArray:
			return lineLength(dir, '[', v.preElems, v.elems, v.postElems, ']')
		case jsonRecord:
			return lineLength(dir, v.preKey, v.key, v.postKey, ':', v.preVal, v.val, v.postVal, ',', v.postComma)
		case jsonElement:
			return lineLength(dir, v.preVal, v.val, v.postVal, ',', v.postComma)
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
			// Add 1 to the returned result to account for preceding space.
			// This accounting is one too few for first block comment.
			b := reflect.ValueOf(v).Bytes()
			switch dir {
			case +1:
				if i := bytes.IndexByte(b, '\n'); i >= 0 {
					return 1 + i, true
				}
			case -1:
				if i := bytes.LastIndexByte(b, '\n'); i >= 0 {
					return 1 + len(b) - (i + 1), true
				}
			}
			return 1 + len(b), false
		case jsonNewlines:
			return 0, v > 0
		case rune:
			if v == ':' || v == ',' {
				return 2, false // Implicit space added
			}
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

func isPrimitive(js jsonValue) bool {
	switch js.(type) {
	case jsonNumber, jsonString, jsonLiteral:
		return true
	default:
		return false
	}
}
