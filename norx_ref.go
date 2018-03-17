// norx_ref.go - Reference (portable) implementation
//
// To the extent possible under law, Yawning Angel has waived all copyright
// and related or neighboring rights to the software, using the Creative
// Commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package norx

import (
	"encoding/binary"
	"math/bits"
)

func permuteRef(s *state, rounds int) {
	// The reference code uses a few macros and has much better
	// readability here, but Go does not have macros.  The "idiomatic"
	// thing to do is to replace the macros with a bunch of functions,
	// but at least as of Go 1.10, the resulting quarter round routine
	// is over the inliner budget.

	for i := 0; i < rounds; i++ {
		// Column step
		// G(S[ 0], S[ 4], S[ 8], S[12]);
		// G(S[ 1], S[ 5], S[ 9], S[13]);
		// G(S[ 2], S[ 6], S[10], S[14]);
		// G(S[ 3], S[ 7], S[11], S[15]);

		s.s[0] = (s.s[0] ^ s.s[4]) ^ ((s.s[0] & s.s[4]) << 1)
		s.s[12] ^= s.s[0]
		s.s[12] = bits.RotateLeft64(s.s[12], -paramR0)
		s.s[8] = (s.s[8] ^ s.s[12]) ^ ((s.s[8] & s.s[12]) << 1)
		s.s[4] ^= s.s[8]
		s.s[4] = bits.RotateLeft64(s.s[4], -paramR1)
		s.s[0] = (s.s[0] ^ s.s[4]) ^ ((s.s[0] & s.s[4]) << 1)
		s.s[12] ^= s.s[0]
		s.s[12] = bits.RotateLeft64(s.s[12], -paramR2)
		s.s[8] = (s.s[8] ^ s.s[12]) ^ ((s.s[8] & s.s[12]) << 1)
		s.s[4] ^= s.s[8]
		s.s[4] = bits.RotateLeft64(s.s[4], -paramR3)

		s.s[1] = (s.s[1] ^ s.s[5]) ^ ((s.s[1] & s.s[5]) << 1)
		s.s[13] ^= s.s[1]
		s.s[13] = bits.RotateLeft64(s.s[13], -paramR0)
		s.s[9] = (s.s[9] ^ s.s[13]) ^ ((s.s[9] & s.s[13]) << 1)
		s.s[5] ^= s.s[9]
		s.s[5] = bits.RotateLeft64(s.s[5], -paramR1)
		s.s[1] = (s.s[1] ^ s.s[5]) ^ ((s.s[1] & s.s[5]) << 1)
		s.s[13] ^= s.s[1]
		s.s[13] = bits.RotateLeft64(s.s[13], -paramR2)
		s.s[9] = (s.s[9] ^ s.s[13]) ^ ((s.s[9] & s.s[13]) << 1)
		s.s[5] ^= s.s[9]
		s.s[5] = bits.RotateLeft64(s.s[5], -paramR3)

		s.s[2] = (s.s[2] ^ s.s[6]) ^ ((s.s[2] & s.s[6]) << 1)
		s.s[14] ^= s.s[2]
		s.s[14] = bits.RotateLeft64(s.s[14], -paramR0)
		s.s[10] = (s.s[10] ^ s.s[14]) ^ ((s.s[10] & s.s[14]) << 1)
		s.s[6] ^= s.s[10]
		s.s[6] = bits.RotateLeft64(s.s[6], -paramR1)
		s.s[2] = (s.s[2] ^ s.s[6]) ^ ((s.s[2] & s.s[6]) << 1)
		s.s[14] ^= s.s[2]
		s.s[14] = bits.RotateLeft64(s.s[14], -paramR2)
		s.s[10] = (s.s[10] ^ s.s[14]) ^ ((s.s[10] & s.s[14]) << 1)
		s.s[6] ^= s.s[10]
		s.s[6] = bits.RotateLeft64(s.s[6], -paramR3)

		s.s[3] = (s.s[3] ^ s.s[7]) ^ ((s.s[3] & s.s[7]) << 1)
		s.s[15] ^= s.s[3]
		s.s[15] = bits.RotateLeft64(s.s[15], -paramR0)
		s.s[11] = (s.s[11] ^ s.s[15]) ^ ((s.s[11] & s.s[15]) << 1)
		s.s[7] ^= s.s[11]
		s.s[7] = bits.RotateLeft64(s.s[7], -paramR1)
		s.s[3] = (s.s[3] ^ s.s[7]) ^ ((s.s[3] & s.s[7]) << 1)
		s.s[15] ^= s.s[3]
		s.s[15] = bits.RotateLeft64(s.s[15], -paramR2)
		s.s[11] = (s.s[11] ^ s.s[15]) ^ ((s.s[11] & s.s[15]) << 1)
		s.s[7] ^= s.s[11]
		s.s[7] = bits.RotateLeft64(s.s[7], -paramR3)

		// Diagonal step
		// G(S[ 0], S[ 5], S[10], S[15]);
		// G(S[ 1], S[ 6], S[11], S[12]);
		// G(S[ 2], S[ 7], S[ 8], S[13]);
		// G(S[ 3], S[ 4], S[ 9], S[14]);

		s.s[0] = (s.s[0] ^ s.s[5]) ^ ((s.s[0] & s.s[5]) << 1)
		s.s[15] ^= s.s[0]
		s.s[15] = bits.RotateLeft64(s.s[15], -paramR0)
		s.s[10] = (s.s[10] ^ s.s[15]) ^ ((s.s[10] & s.s[15]) << 1)
		s.s[5] ^= s.s[10]
		s.s[5] = bits.RotateLeft64(s.s[5], -paramR1)
		s.s[0] = (s.s[0] ^ s.s[5]) ^ ((s.s[0] & s.s[5]) << 1)
		s.s[15] ^= s.s[0]
		s.s[15] = bits.RotateLeft64(s.s[15], -paramR2)
		s.s[10] = (s.s[10] ^ s.s[15]) ^ ((s.s[10] & s.s[15]) << 1)
		s.s[5] ^= s.s[10]
		s.s[5] = bits.RotateLeft64(s.s[5], -paramR3)

		s.s[1] = (s.s[1] ^ s.s[6]) ^ ((s.s[1] & s.s[6]) << 1)
		s.s[12] ^= s.s[1]
		s.s[12] = bits.RotateLeft64(s.s[12], -paramR0)
		s.s[11] = (s.s[11] ^ s.s[12]) ^ ((s.s[11] & s.s[12]) << 1)
		s.s[6] ^= s.s[11]
		s.s[6] = bits.RotateLeft64(s.s[6], -paramR1)
		s.s[1] = (s.s[1] ^ s.s[6]) ^ ((s.s[1] & s.s[6]) << 1)
		s.s[12] ^= s.s[1]
		s.s[12] = bits.RotateLeft64(s.s[12], -paramR2)
		s.s[11] = (s.s[11] ^ s.s[12]) ^ ((s.s[11] & s.s[12]) << 1)
		s.s[6] ^= s.s[11]
		s.s[6] = bits.RotateLeft64(s.s[6], -paramR3)

		s.s[2] = (s.s[2] ^ s.s[7]) ^ ((s.s[2] & s.s[7]) << 1)
		s.s[13] ^= s.s[2]
		s.s[13] = bits.RotateLeft64(s.s[13], -paramR0)
		s.s[8] = (s.s[8] ^ s.s[13]) ^ ((s.s[8] & s.s[13]) << 1)
		s.s[7] ^= s.s[8]
		s.s[7] = bits.RotateLeft64(s.s[7], -paramR1)
		s.s[2] = (s.s[2] ^ s.s[7]) ^ ((s.s[2] & s.s[7]) << 1)
		s.s[13] ^= s.s[2]
		s.s[13] = bits.RotateLeft64(s.s[13], -paramR2)
		s.s[8] = (s.s[8] ^ s.s[13]) ^ ((s.s[8] & s.s[13]) << 1)
		s.s[7] ^= s.s[8]
		s.s[7] = bits.RotateLeft64(s.s[7], -paramR3)

		s.s[3] = (s.s[3] ^ s.s[4]) ^ ((s.s[3] & s.s[4]) << 1)
		s.s[14] ^= s.s[3]
		s.s[14] = bits.RotateLeft64(s.s[14], -paramR0)
		s.s[9] = (s.s[9] ^ s.s[14]) ^ ((s.s[9] & s.s[14]) << 1)
		s.s[4] ^= s.s[9]
		s.s[4] = bits.RotateLeft64(s.s[4], -paramR1)
		s.s[3] = (s.s[3] ^ s.s[4]) ^ ((s.s[3] & s.s[4]) << 1)
		s.s[14] ^= s.s[3]
		s.s[14] = bits.RotateLeft64(s.s[14], -paramR2)
		s.s[9] = (s.s[9] ^ s.s[14]) ^ ((s.s[9] & s.s[14]) << 1)
		s.s[4] ^= s.s[9]
		s.s[4] = bits.RotateLeft64(s.s[4], -paramR3)
	}
}

func padRef(out *[bytesR]byte, in []byte) {
	// Note: This is only called with a zero initialized `out`.
	copy(out[:], in)
	out[len(in)] = 0x01
	out[bytesR-1] |= 0x80
}

func absorbBlockRef(s *state, in []byte, tag uint64) {
	s.s[15] ^= tag
	permuteRef(s, s.rounds)

	for i := 0; i < wordsR; i++ {
		s.s[i] ^= binary.LittleEndian.Uint64(in[i*bytesW:])
	}
}

func absorbLastBlockRef(s *state, in []byte, tag uint64) {
	var lastBlock [bytesR]byte
	padRef(&lastBlock, in)
	absorbBlockRef(s, lastBlock[:], tag)
}

func encryptBlockRef(s *state, out, in []byte) {
	s.s[15] ^= tagPayload
	permuteRef(s, s.rounds)

	for i := 0; i < wordsR; i++ {
		s.s[i] ^= binary.LittleEndian.Uint64(in[i*bytesW:])
		binary.LittleEndian.PutUint64(out[i*bytesW:], s.s[i])
	}
}

func encryptLastBlockRef(s *state, out, in []byte) {
	var lastBlock [bytesR]byte
	padRef(&lastBlock, in)
	encryptBlockRef(s, lastBlock[:], lastBlock[:])
	copy(out, lastBlock[:len(in)])
}

func decryptBlockRef(s *state, out, in []byte) {
	s.s[15] ^= tagPayload
	permuteRef(s, s.rounds)

	for i := 0; i < wordsR; i++ {
		c := binary.LittleEndian.Uint64(in[i*bytesW:])
		binary.LittleEndian.PutUint64(out[i*bytesW:], s.s[i]^c)
		s.s[i] = c
	}
}

func decryptLastBlockRef(s *state, out, in []byte) {
	s.s[15] ^= tagPayload
	permuteRef(s, s.rounds)

	var lastBlock [bytesR]byte
	for i := 0; i < wordsR; i++ {
		binary.LittleEndian.PutUint64(lastBlock[i*bytesW:], s.s[i])
	}

	copy(lastBlock[:], in)
	lastBlock[len(in)] ^= 0x01
	lastBlock[bytesR-1] ^= 0x80

	for i := 0; i < wordsR; i++ {
		c := binary.LittleEndian.Uint64(lastBlock[i*bytesW:])
		binary.LittleEndian.PutUint64(lastBlock[i*bytesW:], s.s[i]^c)
		s.s[i] = c
	}

	copy(out, lastBlock[:len(in)])
	burnBytes(lastBlock[:])
}

func initRef(s *state, key, nonce []byte) {
	// Note: Ensuring a correctly sized key/nonce is the caller's
	// responsibility.

	for i := 0; i < 4; i++ {
		s.s[i] = binary.LittleEndian.Uint64(nonce[i*bytesW:])
		s.s[i+4] = binary.LittleEndian.Uint64(key[i*bytesW:])
	}
	copy(s.s[8:], initializationConstants[8:])

	s.s[12] ^= paramW
	s.s[13] ^= uint64(s.rounds)
	s.s[14] ^= paramP
	s.s[15] ^= paramT

	permuteRef(s, s.rounds)

	for i := 0; i < 4; i++ {
		s.s[i+12] ^= binary.LittleEndian.Uint64(key[i*bytesW:])
	}
}

func absorbDataRef(s *state, in []byte, tag uint64) {
	inLen, off := len(in), 0
	if inLen == 0 {
		return
	}

	for inLen >= bytesR {
		absorbBlockRef(s, in[off:off+bytesR], tag)
		inLen, off = inLen-bytesR, off+bytesR
	}
	absorbLastBlockRef(s, in[off:], tag)
}

func encryptDataRef(s *state, out, in []byte) {
	inLen, off := len(in), 0
	if inLen == 0 {
		return
	}

	for inLen >= bytesR {
		encryptBlockRef(s, out[off:off+bytesR], in[off:off+bytesR])
		inLen, off = inLen-bytesR, off+bytesR
	}
	encryptLastBlockRef(s, out[off:], in[off:])
}

func decryptDataRef(s *state, out, in []byte) {
	inLen, off := len(in), 0
	if inLen == 0 {
		return
	}

	for inLen >= bytesR {
		decryptBlockRef(s, out[off:off+bytesR], in[off:off+bytesR])
		inLen, off = inLen-bytesR, off+bytesR
	}
	decryptLastBlockRef(s, out[off:], in[off:])
}

func finalizeRef(s *state, tag []byte, key []byte) {
	var lastBlock [bytesC]byte

	s.s[15] ^= tagFinal
	permuteRef(s, s.rounds)

	for i := 0; i < 4; i++ {
		s.s[i+12] ^= binary.LittleEndian.Uint64(key[i*bytesW:])
	}

	permuteRef(s, s.rounds)

	for i := 0; i < 4; i++ {
		s.s[i+12] ^= binary.LittleEndian.Uint64(key[i*bytesW:])
		binary.LittleEndian.PutUint64(lastBlock[i*bytesW:], s.s[i+12])
	}

	copy(tag, lastBlock[:bytesT])

	burnBytes(lastBlock[:]) // burn buffer
	burnUint64s(s.s[:])     // at this point we can also burn the state
}
