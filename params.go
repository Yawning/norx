// params.go - Parameters and constants
//
// To the extent possible under law, Yawning Angel has waived all copyright
// and related or neighboring rights to the software, using the Creative
// Commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package norx

const (
	paramW = 64              // Word size
	paramP = 1               // Parallelism degree
	paramT = paramW * 4      // Tag size
	paramN = paramW * 4      // Nonce size
	paramK = paramW * 4      // Key size
	paramB = paramW * 16     // Permutation width
	paramC = paramW * 4      // Capacity
	paramR = paramB - paramC // Rate

	// Rotation constants
	paramR0 = 8
	paramR1 = 19
	paramR2 = 40
	paramR3 = 63

	// Tags
	tagHeader  = 0x01
	tagPayload = 0x02
	tagTrailer = 0x04
	tagFinal   = 0x08
	// tagBranch  = 0x10
	// tagMerge   = 0x20

	bytesW = paramW / 8
	bytesT = paramT / 8
	bytesK = paramK / 8
	bytesR = paramR / 8
	bytesC = paramC / 8
	wordsR = paramR / paramW
)

type state struct {
	s      [16]uint64
	rounds int // Round number (NORX_L)
}

// Taken from "Table 3.4: Initialisation constants".
var initializationConstants = [16]uint64{
	0xE4D324772B91DF79,
	0x3AEC9ABAAEB02CCB,
	0x9DFBA13DB4289311,
	0xEF9EB4BF5A97F2C8,
	0x3F466E92C1532034,
	0xE6E986626CC405C1,
	0xACE40F3B549184E1,
	0xD9CFD35762614477,
	0xB15E641748DE5E6B,
	0xAA95E955E10F8410,
	0x28D1034441A9DD40,
	0x7F31BBF964E93BF5,
	0xB5E9E22493DFFB96,
	0xB980C852479FAFBD,
	0xDA24516BF55EAFD4,
	0x86026AE8536F1501,
}
