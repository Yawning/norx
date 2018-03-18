// api.go - External interface
//
// To the extent possible under law, Yawning Angel has waived all copyright
// and related or neighboring rights to the software, using the Creative
// Commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

// Package norx implements the NORX Authenticated Encryption Algorithm,
// specifically the NORX64-4-1 and NORX64-6-1 variants, as recommended by the
// designers for software implementations on modern 64-bit CPUs.
//
// This implementation is derived from the Public Domain reference
// implementation by Jean-Philippe Aumasson, Philipp Jovanovic, and Samuel
// Neves.
//
// Warning:
// NORX is a rather new authenticated encryption algorithm. The authors are
// confident that it is secure but nevertheless NORX should be considered
// experimental. Therefore, do not use it in your applications!
package norx

import (
	"crypto/cipher"
	"errors"
)

var (
	// ErrInvalidKeySize is the error thrown via a panic when a key is an
	// invalid size.
	ErrInvalidKeySize = errors.New("norx: invalid key size")

	// ErrInvalidNonceSize is the error thrown via a panic when a nonce is
	// an invalid size.
	ErrInvalidNonceSize = errors.New("norx: invalid nonce size")

	// ErrOpen is the error returned when the message authentication fails
	// during an Open call.
	ErrOpen = errors.New("norx: message authentication failed")
)

// AEAD is a parameterized and keyed NORX instance, in the spirit of
// crypto/cipher.AEAD.
type AEAD struct {
	key    []byte
	rounds int
}

// NonceSize returns the size of the nonce that must be passed to Seal and
// Open.
func (ae *AEAD) NonceSize() int {
	return NonceSize
}

// Overhead returns the maximum difference between the lengths of a plaintext
// and its ciphertext.
func (ae *AEAD) Overhead() int {
	return TagSize
}

// Seal encrypts and authenticates plaintext, authenticates the optional
// header and footer (additional data) and, appends the result to dst,
// returning the updated slice. The nonce must be NonceSize() bytes long and
// unique for all time, for a given key.
//
// The plaintext and dst must overlap exactly or not at all. To reuse
// plaintext's storage for the encrypted output, use plaintext[:0] as dst.
func (ae *AEAD) Seal(dst, nonce, plaintext, header, footer []byte) []byte {
	if len(nonce) != NonceSize {
		panic(ErrInvalidNonceSize)
	}
	dst = aeadEncrypt(ae.rounds, dst, header, plaintext, footer, nonce, ae.key)
	return dst
}

// Open decrypts and authenticates ciphertext, authenticates the optonal
// header and footer (additional data) and, if successful, appends the
// resulting plaintext to dst, returning the updated slice. The nonce must
// be NonceSize() bytes long and both it and the additional data must match the
// value passed to Seal.
//
// The ciphertext and dst must overlap exactly or not at all. To reuse
// ciphertext's storage for the decrypted output, use ciphertext[:0] as dst.
//
// Even if the function fails, the contents of dst, up to its capacity,
// may be overwritten.
func (ae *AEAD) Open(dst, nonce, ciphertext, header, footer []byte) ([]byte, error) {
	var err error
	var ok bool

	if len(nonce) != NonceSize {
		panic(ErrInvalidNonceSize)
	}
	dst, ok = aeadDecrypt(ae.rounds, dst, header, ciphertext, footer, nonce, ae.key)
	if !ok {
		err = ErrOpen
	}
	return dst, err
}

// Reset securely purges stored sensitive data from the AEAD instance.
func (ae *AEAD) Reset() {
	burnBytes(ae.key)
}

// ToRuntime converts an AEAD instance to a crypto/cipher.AEAD instance.
//
// The interfaces are distinct as NORX supports both a header and footer as
// additional data, while the runtime interface only has a singular additonal
// data parameter.  The resulting cipher.AEAD instance will use the header
// for additional data if provided, ignoring the footer.
func (ae *AEAD) ToRuntime() cipher.AEAD {
	return &goAEAD{ae}
}

type goAEAD struct {
	aead *AEAD
}

func (ae *goAEAD) NonceSize() int {
	return ae.aead.NonceSize()
}

func (ae *goAEAD) Overhead() int {
	return ae.aead.Overhead()
}

func (ae *goAEAD) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	return ae.aead.Seal(dst, nonce, plaintext, additionalData, nil)
}

func (ae *goAEAD) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	return ae.aead.Open(dst, nonce, ciphertext, additionalData, nil)
}

// New6441 returns a new keyed NORX64-4-1 instance.
func New6441(key []byte) *AEAD {
	return newAEAD(key, 4)
}

// New6461 returns a new keyed NORX64-6-1 instance.
func New6461(key []byte) *AEAD {
	return newAEAD(key, 6)
}

func newAEAD(key []byte, rounds int) *AEAD {
	if len(key) != KeySize {
		panic(ErrInvalidKeySize)
	}

	return &AEAD{
		key:    append([]byte{}, key...),
		rounds: rounds,
	}
}
