// norx.go - High-level interface
//
// To the extent possible under law, Yawning Angel has waived all copyright
// and related or neighboring rights to the software, using the Creative
// Commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package norx

import "crypto/subtle"

const (
	// KeySize is the size of a key in bytes.
	KeySize = 32

	// NonceSize is the size of a nonce in bytes.
	NonceSize = 32

	// TagSize is the size of an authentication tag in bytes.
	TagSize = 32

	// Version is the version of the NORX specification implemented.
	Version = "3.0"
)

func aeadEncrypt(l int, c, a, m, z, nonce, key []byte) []byte {
	var k [bytesK]byte
	s := &state{rounds: l}
	mLen := len(m)

	ret, out := sliceForAppend(c, mLen+bytesT)

	copy(k[:], key)
	hardwareAccelImpl.initFn(s, k[:], nonce)
	hardwareAccelImpl.absorbDataFn(s, a, tagHeader)
	hardwareAccelImpl.encryptDataFn(s, out, m)
	hardwareAccelImpl.absorbDataFn(s, z, tagTrailer)
	hardwareAccelImpl.finalizeFn(s, out[mLen:], k[:])

	burnUint64s(s.s[:])
	burnBytes(k[:])

	return ret
}

func aeadDecrypt(l int, m, a, c, z, nonce, key []byte) ([]byte, bool) {
	var k [bytesK]byte
	var tag [bytesT]byte
	s := &state{rounds: l}
	cLen := len(c)

	if cLen < bytesT {
		return nil, false
	}

	mLen := cLen - bytesT
	ret, out := sliceForAppend(m, mLen)

	copy(k[:], key)
	hardwareAccelImpl.initFn(s, k[:], nonce)
	hardwareAccelImpl.absorbDataFn(s, a, tagHeader)
	hardwareAccelImpl.decryptDataFn(s, out, c[:mLen])
	hardwareAccelImpl.absorbDataFn(s, z, tagTrailer)
	hardwareAccelImpl.finalizeFn(s, tag[:], k[:])

	srcTag := c[mLen:]
	ok := subtle.ConstantTimeCompare(srcTag, tag[:]) == 1
	if !ok && mLen > 0 { // burn decrypted plaintext on auth failure
		burnBytes(out[:mLen])
		ret = nil
	}

	burnUint64s(s.s[:])
	burnBytes(k[:])

	return ret, ok
}

// Shamelessly stolen from the Go runtime library.
func sliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

func init() {
	if KeySize != bytesK {
		panic("BUG: KeySize != paramK/8")
	}
	if NonceSize != paramN/8 {
		panic("BUG: NonceSize != paramN/8")
	}
	if TagSize != bytesT {
		panic("BUG: TagSize != bytesT")
	}
}
