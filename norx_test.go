// norx_test.go - NORX tests
//
// To the extent possible under law, Yawning Angel has waived all copyright
// and related or neighboring rights to the software, using the Creative
// Commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package norx

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	katCipherTexts = map[string][]byte{
		"NORX64-4-1": kat6441,
		"NORX64-6-1": kat6461,
	}
	testVectorCipherTexts = map[string][]byte{
		"NORX64-4-1": vec6441,
		"NORX64-6-1": vec6461,
	}
	canAccelerate bool
)

func mustInitHardwareAcceleration() {
	initHardwareAcceleration()
	if !IsHardwareAccelerated() {
		panic("initHardwareAcceleration() failed")
	}
}

func TestKAT(t *testing.T) {
	forceDisableHardwareAcceleration()
	doTestKAT(t)

	if !canAccelerate {
		t.Log("Hardware acceleration not supported on this host.")
		return
	}
	mustInitHardwareAcceleration()
	doTestKAT(t)
}

func doTestKAT(t *testing.T) {
	testRounds := []int{4, 6}
	impl := "_" + hardwareAccelImpl.name

	for _, l := range testRounds {
		n := fmt.Sprintf("NORX64-%d-1", l)
		t.Run(n+"_Vector"+impl, func(t *testing.T) { doTestVector(t, n, l) })
		t.Run(n+"_KAT"+impl, func(t *testing.T) { doTestKATFull(t, n, l) })
	}
}

func doTestVector(t *testing.T, vn string, l int) {
	require := require.New(t)
	var k, n [32]byte
	var a, m, z [128]byte

	for i := range k {
		k[i] = byte(i)
		n[i] = byte(i + 32)
	}
	for i := range a {
		a[i] = byte(i)
		m[i] = byte(i)
		z[i] = byte(i)
	}

	ctVec := testVectorCipherTexts[vn]
	if ctVec == nil {
		panic("missing test vector cipher text")
	}

	aead := newTestAEAD(k[:], l)
	require.Equal(NonceSize, aead.NonceSize(), "NonceSize()")
	require.Equal(TagSize, aead.Overhead(), "Overhead()")

	ct := aead.Seal(nil, n[:], m[:], a[:], z[:])
	require.Equal(ctVec, ct, "aeadSeal()")

	pt, err := aead.Open(nil, n[:], ct, a[:], z[:])
	require.NoError(err, "Open()")
	require.Equal(m[:], pt, "Open()")
}

func doTestKATFull(t *testing.T, vn string, l int) {
	require := require.New(t)
	var w, h [256]byte
	var k, n [32]byte

	for i := range w {
		w[i] = byte(255 & (i*197 + 123))
	}
	for i := range h {
		h[i] = byte(255 & (i*193 + 123))
	}
	for i := range k {
		k[i] = byte(255 & (i*191 + 123))
	}
	for i := range n {
		n[i] = byte(255 & (i*181 + 123))
	}

	var katAcc []byte
	kat := katCipherTexts[vn]
	katOff := 0
	if kat == nil {
		panic("missing KAT cipher texts")
	}

	aead := newTestAEAD(k[:], l).ToRuntime()
	require.Equal(NonceSize, aead.NonceSize(), "NonceSize()")
	require.Equal(TagSize, aead.Overhead(), "Overhead()")

	for i := range w {
		katAcc = aead.Seal(katAcc, n[:], w[:i], h[:i])
		c := katAcc[katOff:]
		require.Len(c, i+TagSize, "Seal(): len(c) %d", i)
		require.Equal(kat[katOff:katOff+len(c)], c)

		m, err := aead.Open(nil, n[:], c, h[:i])
		require.NoError(err, "Open(): %d", i)
		require.Len(m, i, "Open(): len(m) %d", i)
		if len(m) != 0 {
			require.Equal(m, w[:i], "Open(): m %d", i)
		}

		katOff += len(c)
	}
	require.Equal(kat, katAcc, "Final concatenated cipher texts.")
}

func newTestAEAD(k []byte, l int) *AEAD {
	switch l {
	case 4:
		return New6441(k[:])
	case 6:
		return New6461(k[:])
	default:
		panic("unsupported round parameter")
	}
}

func BenchmarkNORX(b *testing.B) {
	forceDisableHardwareAcceleration()
	doBenchmarkNORX(b)

	if !canAccelerate {
		b.Log("Hardware acceleration not supported on this host.")
		return
	}
	mustInitHardwareAcceleration()
	doBenchmarkNORX(b)
}

func doBenchmarkNORX(b *testing.B) {
	benchRounds := []int{4, 6}
	benchSizes := []int{8, 64, 576, 1536, 4096, 1024768}
	impl := "_" + hardwareAccelImpl.name

	for _, l := range benchRounds {
		n := fmt.Sprintf("NORX64-%d-1", l)
		for _, sz := range benchSizes {
			bn := n + impl + "_"
			sn := fmt.Sprintf("_%d", sz)
			b.Run(bn+"Encrypt"+sn, func(b *testing.B) { doBenchmarkAEADEncrypt(b, l, sz) })
			b.Run(bn+"Decrypt"+sn, func(b *testing.B) { doBenchmarkAEADDecrypt(b, l, sz) })
		}
	}
}

func doBenchmarkAEADEncrypt(b *testing.B, l, sz int) {
	b.StopTimer()
	b.SetBytes(int64(sz))

	nonce, key := make([]byte, NonceSize), make([]byte, KeySize)
	m, c := make([]byte, sz), make([]byte, 0, sz+TagSize)
	rand.Read(nonce)
	rand.Read(key)
	rand.Read(m)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c = c[:0]

		c = aeadEncrypt(l, c, nil, m, nil, nonce, key)
		if len(c) != sz+TagSize {
			b.Fatalf("aeadEncrypt failed")
		}
	}
}

func doBenchmarkAEADDecrypt(b *testing.B, l, sz int) {
	b.StopTimer()
	b.SetBytes(int64(sz))

	nonce, key := make([]byte, NonceSize), make([]byte, KeySize)
	m, c, d := make([]byte, sz), make([]byte, 0, sz+TagSize), make([]byte, 0, sz)
	rand.Read(nonce)
	rand.Read(key)
	rand.Read(m)

	c = aeadEncrypt(l, c, nil, m, nil, nonce, key)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		d = d[:0]

		var ok bool
		d, ok = aeadDecrypt(l, d, nil, c, nil, nonce, key)
		if !ok {
			b.Fatalf("aeadDecrypt failed")
		}
	}
	b.StopTimer()

	if !bytes.Equal(m, d) {
		b.Fatalf("aeadDecrypt output mismatch")
	}
}

func init() {
	canAccelerate = IsHardwareAccelerated()
}
