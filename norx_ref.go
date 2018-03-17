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

	// Performance: Explicitly load the state into temp vars, and write
	// it back on completion since the compiler will do all of the
	// loads/stores otherwise.
	s0, s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11, s12, s13, s14, s15 := s.s[0], s.s[1], s.s[2], s.s[3], s.s[4], s.s[5], s.s[6], s.s[7], s.s[8], s.s[9], s.s[10], s.s[11], s.s[12], s.s[13], s.s[14], s.s[15]

	for i := 0; i < rounds; i++ {
		// Column step
		// G(S[ 0], S[ 4], S[ 8], S[12]);
		// G(S[ 1], S[ 5], S[ 9], S[13]);
		// G(S[ 2], S[ 6], S[10], S[14]);
		// G(S[ 3], S[ 7], S[11], S[15]);

		s0 = (s0 ^ s4) ^ ((s0 & s4) << 1)
		s12 ^= s0
		s12 = bits.RotateLeft64(s12, -paramR0)
		s8 = (s8 ^ s12) ^ ((s8 & s12) << 1)
		s4 ^= s8
		s4 = bits.RotateLeft64(s4, -paramR1)
		s0 = (s0 ^ s4) ^ ((s0 & s4) << 1)
		s12 ^= s0
		s12 = bits.RotateLeft64(s12, -paramR2)
		s8 = (s8 ^ s12) ^ ((s8 & s12) << 1)
		s4 ^= s8
		s4 = bits.RotateLeft64(s4, -paramR3)

		s1 = (s1 ^ s5) ^ ((s1 & s5) << 1)
		s13 ^= s1
		s13 = bits.RotateLeft64(s13, -paramR0)
		s9 = (s9 ^ s13) ^ ((s9 & s13) << 1)
		s5 ^= s9
		s5 = bits.RotateLeft64(s5, -paramR1)
		s1 = (s1 ^ s5) ^ ((s1 & s5) << 1)
		s13 ^= s1
		s13 = bits.RotateLeft64(s13, -paramR2)
		s9 = (s9 ^ s13) ^ ((s9 & s13) << 1)
		s5 ^= s9
		s5 = bits.RotateLeft64(s5, -paramR3)

		s2 = (s2 ^ s6) ^ ((s2 & s6) << 1)
		s14 ^= s2
		s14 = bits.RotateLeft64(s14, -paramR0)
		s10 = (s10 ^ s14) ^ ((s10 & s14) << 1)
		s6 ^= s10
		s6 = bits.RotateLeft64(s6, -paramR1)
		s2 = (s2 ^ s6) ^ ((s2 & s6) << 1)
		s14 ^= s2
		s14 = bits.RotateLeft64(s14, -paramR2)
		s10 = (s10 ^ s14) ^ ((s10 & s14) << 1)
		s6 ^= s10
		s6 = bits.RotateLeft64(s6, -paramR3)

		s3 = (s3 ^ s7) ^ ((s3 & s7) << 1)
		s15 ^= s3
		s15 = bits.RotateLeft64(s15, -paramR0)
		s11 = (s11 ^ s15) ^ ((s11 & s15) << 1)
		s7 ^= s11
		s7 = bits.RotateLeft64(s7, -paramR1)
		s3 = (s3 ^ s7) ^ ((s3 & s7) << 1)
		s15 ^= s3
		s15 = bits.RotateLeft64(s15, -paramR2)
		s11 = (s11 ^ s15) ^ ((s11 & s15) << 1)
		s7 ^= s11
		s7 = bits.RotateLeft64(s7, -paramR3)

		// Diagonal step
		// G(S[ 0], S[ 5], S[10], S[15]);
		// G(S[ 1], S[ 6], S[11], S[12]);
		// G(S[ 2], S[ 7], S[ 8], S[13]);
		// G(S[ 3], S[ 4], S[ 9], S[14]);

		s0 = (s0 ^ s5) ^ ((s0 & s5) << 1)
		s15 ^= s0
		s15 = bits.RotateLeft64(s15, -paramR0)
		s10 = (s10 ^ s15) ^ ((s10 & s15) << 1)
		s5 ^= s10
		s5 = bits.RotateLeft64(s5, -paramR1)
		s0 = (s0 ^ s5) ^ ((s0 & s5) << 1)
		s15 ^= s0
		s15 = bits.RotateLeft64(s15, -paramR2)
		s10 = (s10 ^ s15) ^ ((s10 & s15) << 1)
		s5 ^= s10
		s5 = bits.RotateLeft64(s5, -paramR3)

		s1 = (s1 ^ s6) ^ ((s1 & s6) << 1)
		s12 ^= s1
		s12 = bits.RotateLeft64(s12, -paramR0)
		s11 = (s11 ^ s12) ^ ((s11 & s12) << 1)
		s6 ^= s11
		s6 = bits.RotateLeft64(s6, -paramR1)
		s1 = (s1 ^ s6) ^ ((s1 & s6) << 1)
		s12 ^= s1
		s12 = bits.RotateLeft64(s12, -paramR2)
		s11 = (s11 ^ s12) ^ ((s11 & s12) << 1)
		s6 ^= s11
		s6 = bits.RotateLeft64(s6, -paramR3)

		s2 = (s2 ^ s7) ^ ((s2 & s7) << 1)
		s13 ^= s2
		s13 = bits.RotateLeft64(s13, -paramR0)
		s8 = (s8 ^ s13) ^ ((s8 & s13) << 1)
		s7 ^= s8
		s7 = bits.RotateLeft64(s7, -paramR1)
		s2 = (s2 ^ s7) ^ ((s2 & s7) << 1)
		s13 ^= s2
		s13 = bits.RotateLeft64(s13, -paramR2)
		s8 = (s8 ^ s13) ^ ((s8 & s13) << 1)
		s7 ^= s8
		s7 = bits.RotateLeft64(s7, -paramR3)

		s3 = (s3 ^ s4) ^ ((s3 & s4) << 1)
		s14 ^= s3
		s14 = bits.RotateLeft64(s14, -paramR0)
		s9 = (s9 ^ s14) ^ ((s9 & s14) << 1)
		s4 ^= s9
		s4 = bits.RotateLeft64(s4, -paramR1)
		s3 = (s3 ^ s4) ^ ((s3 & s4) << 1)
		s14 ^= s3
		s14 = bits.RotateLeft64(s14, -paramR2)
		s9 = (s9 ^ s14) ^ ((s9 & s14) << 1)
		s4 ^= s9
		s4 = bits.RotateLeft64(s4, -paramR3)
	}

	s.s[0], s.s[1], s.s[2], s.s[3], s.s[4], s.s[5], s.s[6], s.s[7], s.s[8], s.s[9], s.s[10], s.s[11], s.s[12], s.s[13], s.s[14], s.s[15] = s0, s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11, s12, s13, s14, s15
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
