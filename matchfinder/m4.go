package matchfinder

import (
	"encoding/binary"
	"math/bits"
	"runtime"
)

const (
	ssapBits = 17
	ssapMask = (1 << ssapBits) - 1
)

// M4 is an implementation of the MatchFinder
// interface that uses a simple hash table to find matches,
// but the advanced parsing technique from
// https://fastcompression.blogspot.com/2011/12/advanced-parsing-strategies.html,
// except that it looks for matches at every input position.
type M4 struct {
	// MaxDistance is the maximum distance (in bytes) to look back for
	// a match. The default is 65535.
	MaxDistance int

	// MinLength is the length of the shortest match to return.
	// The default is 4.
	MinLength int

	// HashLen is the number of bytes to use to calculate the hashes.
	// The maximum is 8 and the default is 6.
	HashLen int

	table [1 << ssapBits]uint32

	history []byte
}

func (q *M4) Reset() {
	q.table = [1 << ssapBits]uint32{}
	q.history = q.history[:0]
}

func (q *M4) FindMatches(dst []Match, src []byte) []Match {
	if q.MaxDistance == 0 {
		q.MaxDistance = 65535
	}
	if q.MinLength == 0 {
		q.MinLength = 4
	}
	if q.HashLen == 0 {
		q.HashLen = 6
	}
	var nextEmit int

	if len(q.history) > q.MaxDistance*2 {
		// Trim down the history buffer.
		delta := len(q.history) - q.MaxDistance
		copy(q.history, q.history[delta:])
		q.history = q.history[:q.MaxDistance]

		for i, v := range q.table {
			newV := int(v) - delta
			if newV < 0 {
				newV = 0
			}
			q.table[i] = uint32(newV)
		}
	}

	// Append src to the history buffer.
	nextEmit = len(q.history)
	q.history = append(q.history, src...)
	src = q.history

	// matches stores the matches that have been found but not emitted,
	// in reverse order. (matches[0] is the most recent one.)
	var matches [3]absoluteMatch
	for i := nextEmit; i < len(src)-7; i++ {
		if matches[0] != (absoluteMatch{}) && i >= matches[0].End {
			// We have found some matches, and we're far enough along that we probably
			// won't find overlapping matches, so we might as well emit them.
			if matches[1] != (absoluteMatch{}) {
				if matches[1].End > matches[0].Start {
					matches[1].End = matches[0].Start
				}
				if matches[1].End-matches[1].Start >= q.MinLength {
					dst = append(dst, Match{
						Unmatched: matches[1].Start - nextEmit,
						Length:    matches[1].End - matches[1].Start,
						Distance:  matches[1].Start - matches[1].Match,
					})
					nextEmit = matches[1].End
				}
			}
			dst = append(dst, Match{
				Unmatched: matches[0].Start - nextEmit,
				Length:    matches[0].End - matches[0].Start,
				Distance:  matches[0].Start - matches[0].Match,
			})
			nextEmit = matches[0].End
			matches = [3]absoluteMatch{}
		}

		// Now look for a match.
		h := ((binary.LittleEndian.Uint64(src[i:]) & (1<<(8*q.HashLen) - 1)) * hashMul64) >> (64 - ssapBits)
		candidate := int(q.table[h&ssapMask])
		q.table[h&ssapMask] = uint32(i)

		if candidate == 0 || i-candidate > q.MaxDistance || i-candidate == matches[0].Start-matches[0].Match {
			continue
		}
		if binary.LittleEndian.Uint32(src[candidate:]) != binary.LittleEndian.Uint32(src[i:]) {
			continue
		}

		// We have a 4-byte match now.

		start := i
		match := candidate
		end := extendMatch(src, match+4, start+4)
		for start > nextEmit && match > 0 && src[start-1] == src[match-1] {
			start--
			match--
		}
		if end-start <= matches[0].End-matches[0].Start {
			continue
		}

		matches = [3]absoluteMatch{
			absoluteMatch{
				Start: start,
				End:   end,
				Match: match,
			},
			matches[0],
			matches[1],
		}

		if matches[2] == (absoluteMatch{}) {
			continue
		}

		// We have three matches, so it's time to emit one and/or eliminate one.
		switch {
		case matches[0].Start < matches[2].End:
			// The first and third matches overlap; discard the one in between.
			matches = [3]absoluteMatch{
				matches[0],
				matches[2],
				absoluteMatch{},
			}

		case matches[0].Start < matches[2].End+q.MinLength:
			// The first and third matches don't overlap, but there's no room for
			// another match between them. Emit the first match and discard the second.
			dst = append(dst, Match{
				Unmatched: matches[2].Start - nextEmit,
				Length:    matches[2].End - matches[2].Start,
				Distance:  matches[2].Start - matches[2].Match,
			})
			nextEmit = matches[2].End
			matches = [3]absoluteMatch{
				matches[0],
				absoluteMatch{},
				absoluteMatch{},
			}

		default:
			// Emit the first match, shortening it if necessary to avoid overlap with the second.
			if matches[2].End > matches[1].Start {
				matches[2].End = matches[1].Start
			}
			if matches[2].End-matches[2].Start >= q.MinLength {
				dst = append(dst, Match{
					Unmatched: matches[2].Start - nextEmit,
					Length:    matches[2].End - matches[2].Start,
					Distance:  matches[2].Start - matches[2].Match,
				})
				nextEmit = matches[2].End
			}
			matches[2] = absoluteMatch{}
		}
	}

	// We've found all the matches now; emit the remaining ones.
	if matches[1] != (absoluteMatch{}) {
		if matches[1].End > matches[0].Start {
			matches[1].End = matches[0].Start
		}
		if matches[1].End-matches[1].Start >= q.MinLength {
			dst = append(dst, Match{
				Unmatched: matches[1].Start - nextEmit,
				Length:    matches[1].End - matches[1].Start,
				Distance:  matches[1].Start - matches[1].Match,
			})
			nextEmit = matches[1].End
		}
	}
	if matches[0] != (absoluteMatch{}) {
		dst = append(dst, Match{
			Unmatched: matches[0].Start - nextEmit,
			Length:    matches[0].End - matches[0].Start,
			Distance:  matches[0].Start - matches[0].Match,
		})
		nextEmit = matches[0].End
	}

	if nextEmit < len(src) {
		dst = append(dst, Match{
			Unmatched: len(src) - nextEmit,
		})
	}

	return dst
}

const hashMul64 = 0x1E35A7BD1E35A7BD

// An absoluteMatch is like a Match, but it stores indexes into the byte
// stream instead of lengths.
type absoluteMatch struct {
	// Start is the index of the first byte.
	Start int

	// End is the index of the byte after the last byte
	// (so that End - Start = Length).
	End int

	// Match is the index of the previous data that matches
	// (Start - Match = Distance).
	Match int
}

// extendMatch returns the largest k such that k <= len(src) and that
// src[i:i+k-j] and src[j:k] have the same contents.
//
// It assumes that:
//
//	0 <= i && i < j && j <= len(src)
func extendMatch(src []byte, i, j int) int {
	switch runtime.GOARCH {
	case "amd64":
		// As long as we are 8 or more bytes before the end of src, we can load and
		// compare 8 bytes at a time. If those 8 bytes are equal, repeat.
		for j+8 < len(src) {
			iBytes := binary.LittleEndian.Uint64(src[i:])
			jBytes := binary.LittleEndian.Uint64(src[j:])
			if iBytes != jBytes {
				// If those 8 bytes were not equal, XOR the two 8 byte values, and return
				// the index of the first byte that differs. The BSF instruction finds the
				// least significant 1 bit, the amd64 architecture is little-endian, and
				// the shift by 3 converts a bit index to a byte index.
				return j + bits.TrailingZeros64(iBytes^jBytes)>>3
			}
			i, j = i+8, j+8
		}
	case "386":
		// On a 32-bit CPU, we do it 4 bytes at a time.
		for j+4 < len(src) {
			iBytes := binary.LittleEndian.Uint32(src[i:])
			jBytes := binary.LittleEndian.Uint32(src[j:])
			if iBytes != jBytes {
				return j + bits.TrailingZeros32(iBytes^jBytes)>>3
			}
			i, j = i+4, j+4
		}
	}
	for ; j < len(src) && src[i] == src[j]; i, j = i+1, j+1 {
	}
	return j
}