// Copyright 2017, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsonfmt

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// tokenRegex is a regexp for any valid JSON token (not including comments).
var tokenRegex = regexp.MustCompile(fmt.Sprintf("^(%s)",
	strings.Join([]string{stringRx, numberRx, literalRx, spaceRx}, "|")))

// alignJSON inserts spaces into certain lines in the input to
// to achieve some form of vertical alignment across related lines.
func alignJSON(in []byte, alignLimit int) (out []byte) {
	if alignLimit <= 0 {
		return in
	}

	// Parse the JSON structure on a line-by-line basis.
	// TODO: Should we just record this information in format?
	lines := bytes.Split(in, newlineBytes)
	infos := make([]lineInfo, len(lines))
	var inBlockComment bool
	for i, line := range lines {
		nl := len(line)
		info := &infos[i]
		if inBlockComment {
			info.numIndents = -1 // Invalidate alignment for this line
		} else {
			info.numIndents = len(indentPrefix(line))
		}

		var hasTokens bool
		var lastComma bool
		for len(line) > 0 {
			if inBlockComment {
				if j := bytes.Index(line, []byte("*/")); j >= 0 {
					inBlockComment = false
					line = line[j+len("*/"):]
				} else {
					line = nil
				}
				continue
			}
			if b := tokenRegex.Find(line); len(b) > 0 {
				hasTokens = true
				line = line[len(b):] // Ignore any other JSON token
				continue
			}
			switch line[0] {
			case '/': // Handle comments
				inBlockComment = len(line) > 1 && line[1] == '*'
				if !inBlockComment {
					if hasTokens { // First comment is not "end" comment
						info.commentOffset = nl - len(line)
					}
					line = nil // Rest of line is comment
				}
				continue
			case '{', '[': // Record open braces
				info.bracesBalance++
			case '}', ']': // Record close braces
				info.bracesBalance--
			case ':', ',': // Record alignment markers
				if len(info.alignOffsets) > 0 && info.alignBracesLevel == info.bracesBalance {
					info.alignOffsets = append(info.alignOffsets, nl-len(line)+1)
					lastComma = line[0] == ','
				} else if len(info.alignOffsets) == 0 && line[0] == ':' {
					info.alignOffsets = append(info.alignOffsets, nl-len(line)+1)
					info.alignBracesLevel = info.bracesBalance
				}
			default:
				return in // Invalid JSON
			}
			line = line[1:]
		}
		if lastComma {
			info.alignOffsets = info.alignOffsets[:len(info.alignOffsets)-1]
		}
	}

	// For each indented group, align according to each marker and end comment.
	for i := 0; i < len(infos); {
		infoGroup := nextGroup(infos[i:])
		if len(infoGroup) > 1 {
			var maxAligns int // Maximum number of align markers
			for _, li := range infoGroup {
				if len(li.alignOffsets) > maxAligns {
					maxAligns = len(li.alignOffsets)
				}
			}
			lineGroup := lines[i : i+len(infoGroup)]
			for j := 0; j < maxAligns; j++ {
				alignPositions(lineGroup, infoGroup, j, alignLimit) // Align markers
			}
			alignPositions(lineGroup, infoGroup, -1, alignLimit) // Align comments
		}
		i += len(infoGroup)
	}
	return bytes.Join(lines, newlineBytes)
}

type lineInfo struct {
	numIndents       int   // Number of preceding indents
	bracesBalance    int   // Incremented for '{' or '['; decremented for '}' or ']'
	alignBracesLevel int   // Value of bracesBalance at first append
	alignOffsets     []int // Offset of ':' and ',' in same brace level
	commentOffset    int   // Offset of trailing comment
}

// alignOffset returns the offset into the line for the i-th align marker.
// A negative i returns the offset of the end comment.
// An offset of zero implies that it is invalid.
func (li *lineInfo) alignOffset(i int) int {
	if i < 0 {
		return li.commentOffset
	}
	if i < len(li.alignOffsets) {
		return li.alignOffsets[i]
	}
	return 0
}

// insertSpaces inserts n spaces at pos into line.
// This also updates lineInfo to account for the shifted offsets.
func (li *lineInfo) insertSpaces(line []byte, pos, n int) []byte {
	if n <= 0 || pos >= len(line) {
		return line
	}
	for i, p := range li.alignOffsets {
		if p > pos {
			li.alignOffsets[i] = p + n
		}
	}
	if li.commentOffset > pos {
		li.commentOffset += n
	}
	// TODO: This is inefficient because of repeated copying.
	line = append(line[:pos:pos], append(bytes.Repeat(spaceBytes, n), line[pos:]...)...)
	return line
}

// nextGroup returns the next set of infos that are grouped together
// based on their indent level. This always returns a non-empty slice
// if the input is non-empty.
func nextGroup(infos []lineInfo) []lineInfo {
	if len(infos) <= 1 {
		return infos
	}
	li0 := infos[0]
	n := 1
	if li0.numIndents < 0 || li0.bracesBalance != 0 {
		return infos[:n]
	}
	for _, li := range infos[1:] {
		if li.numIndents != li0.numIndents || li.bracesBalance != 0 {
			return infos[:n]
		}
		n++
	}
	return infos[:n]
}

// alignPositions takes in a set of lines (and their associated lineInfos)
// and aligns each line for the idx-th (-1 for end comment) align marker.
// The limit controls the maximum number of spaces added to achieve alignment.
func alignPositions(lines [][]byte, infos []lineInfo, idx int, limit int) {
	for len(infos) > 1 {
		li0 := infos[0]
		pos0 := li0.alignOffset(idx)
		if pos0 == 0 {
			infos, lines = infos[1:], lines[1:]
			continue
		}

		// Cluster infos into a subgroup to align together.
		n := 1
		maxPos := pos0
		for _, li := range infos[1:] {
			pos := li.alignOffset(idx)
			if pos == 0 || pos > pos0+limit || pos < pos0-limit {
				break
			}
			if maxPos < pos {
				maxPos = pos
			}
			n++
		}

		// Align the subgroup.
		for i, li := range infos[:n] {
			pos := li.alignOffset(idx)
			lines[i] = infos[i].insertSpaces(lines[i], pos, maxPos-pos)
		}

		infos, lines = infos[n:], lines[n:]
	}
}
