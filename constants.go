package brotli

/* Copyright 2016 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Specification: 3.3. Alphabet sizes: insert-and-copy length */
const numLiteralSymbols = 256

const numCommandSymbols = 704

const numBlockLenSymbols = 26

/* Specification: 3.5. Complex prefix codes */
const repeatPreviousCodeLength = 16

const repeatZeroCodeLength = 17

const codeLengthCodes = (repeatZeroCodeLength + 1)

/* "code length of 8 is repeated" */
const initialRepeatedCodeLength = 8

/* "Large Window Brotli" */
const largeMaxDistanceBits = 62

const largeMinWbits = 10

const largeMaxWbits = 30

/* Specification: 4. Encoding of distances */
const numDistanceShortCodes = 16

const maxNpostfix = 3

const maxDistanceBits = 24

func distanceAlphabetSize(NPOSTFIX uint, NDIRECT uint, MAXNBITS uint) uint {
	return numDistanceShortCodes + NDIRECT + uint(MAXNBITS<<(NPOSTFIX+1))
}

const maxAllowedDistance = 0x7FFFFFFC

/* 7.1. Context modes and context ID lookup for literals */
/* "context IDs for literals are in the range of 0..63" */
const literalContextBits = 6

/* 7.2. Context ID for distances */
const distanceContextBits = 2

/* 9.1. Format of the Stream Header */
/* Number of slack bytes for window size. Don't confuse
   with BROTLI_NUM_DISTANCE_SHORT_CODES. */
const windowGap = 16

/*
Represents the range of values belonging to a prefix code:

	[offset, offset + 2^nbits)
*/
type prefixCodeRange struct {
	offset uint32
	nbits  uint32
}

var kBlockLengthPrefixCode = [numBlockLenSymbols]prefixCodeRange{
	prefixCodeRange{1, 2},
	prefixCodeRange{5, 2},
	prefixCodeRange{9, 2},
	prefixCodeRange{13, 2},
	prefixCodeRange{17, 3},
	prefixCodeRange{25, 3},
	prefixCodeRange{33, 3},
	prefixCodeRange{41, 3},
	prefixCodeRange{49, 4},
	prefixCodeRange{65, 4},
	prefixCodeRange{81, 4},
	prefixCodeRange{97, 4},
	prefixCodeRange{113, 5},
	prefixCodeRange{145, 5},
	prefixCodeRange{177, 5},
	prefixCodeRange{209, 5},
	prefixCodeRange{241, 6},
	prefixCodeRange{305, 6},
	prefixCodeRange{369, 7},
	prefixCodeRange{497, 8},
	prefixCodeRange{753, 9},
	prefixCodeRange{1265, 10},
	prefixCodeRange{2289, 11},
	prefixCodeRange{4337, 12},
	prefixCodeRange{8433, 13},
	prefixCodeRange{16625, 24},
}
