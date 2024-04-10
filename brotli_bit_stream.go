package brotli

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
