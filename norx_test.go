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
		"NORX-64-4-1": kat6441,
		"NORX-64-6-1": kat6461,
	}
	testVectorCipherTexts = map[string][]byte{
		"NORX-64-4-1": vec6441,
		"NORX-64-6-1": vec6461,
	}
	canAccelerate bool
)

func mustInitHardwareAcceleration() {
	// initHardwareAcceleration()
	if !IsHardwareAccelerated() {
		panic("initHardwareAcceleration() failed")
	}
}

func TestF(t *testing.T) {
	forceDisableHardwareAcceleration()
	doTestF(t)

	if !canAccelerate {
		t.Log("Hardware acceleration not supported on this host.")
		return
	}
	mustInitHardwareAcceleration()
	doTestF(t)
}

func doTestF(t *testing.T) {
	impl := "_" + hardwareAccelImpl.name
	t.Run("F"+impl, func(t *testing.T) {
		require := require.New(t)
		s := &state{}
		for i := range s.s {
			s.s[i] = uint64(i)
		}
		hardwareAccelImpl.permuteFn(s, 2)
		require.Equal(initializationConstants, s.s, "pre-generated vs calculated")
	})
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
		n := fmt.Sprintf("NORX-64-%d-1", l)
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

	ct := aeadEncrypt(l, nil, a[:], m[:], z[:], n[:], k[:])
	require.Equal(ctVec, ct, "aeadEncrypt()")

	pt, ok := aeadDecrypt(l, nil, a[:], ct, z[:], n[:], k[:])
	require.True(ok, "aeadDecrypt(): ok")
	require.Equal(m[:], pt, "aeadDecrupt()")
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

	for i := range w {
		katAcc = aeadEncrypt(l, katAcc, h[:i], w[:i], nil, n[:], k[:])
		c := katAcc[katOff:]
		require.Len(c, i+TagSize, "aeadEncrypt(): len(c) %d", i)
		require.Equal(kat[katOff:katOff+len(c)], c)

		m, ok := aeadDecrypt(l, nil, h[:i], c, nil, n[:], k[:])
		require.True(ok, "aeadDecrypt(): ok %d", i)
		require.Len(m, i, "aeadDecrypt(): len(m) %d", i)
		if len(m) != 0 {
			require.Equal(m, w[:i], "aeadDecrupt(): m %d", i)
		}

		katOff += len(c)
	}
	require.Equal(kat, katAcc, "Final concatenated cipher texts.")
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
		n := fmt.Sprintf("NORX-64-%d-1", l)
		for _, sz := range benchSizes {
			bn := n + impl + "_"
			sn := fmt.Sprintf("_%d", sz)
			b.Run(bn+"Encrypt"+sn, func(b *testing.B) { doBenchmarkAEAD(b, l, sz, true) })
			b.Run(bn+"Decrypt"+sn, func(b *testing.B) { doBenchmarkAEAD(b, l, sz, false) })
		}
	}
}

func doBenchmarkAEAD(b *testing.B, l, sz int, isEncrypt bool) {
	b.StopTimer()
	b.SetBytes(int64(sz))

	nonce, key := make([]byte, NonceSize), make([]byte, KeySize)
	m, c, d := make([]byte, sz), make([]byte, 0, sz+TagSize), make([]byte, 0, sz)
	rand.Read(nonce)
	rand.Read(key)
	rand.Read(m)

	for i := 0; i < b.N; i++ {
		c = c[:0]
		d = d[:0]

		if isEncrypt {
			b.StartTimer()
		}
		c = aeadEncrypt(l, c, nil, m, nil, nonce, key)
		if isEncrypt {
			b.StopTimer()
		} else {
			b.StartTimer()
		}

		var ok bool
		d, ok = aeadDecrypt(l, d, nil, c, nil, nonce, key)
		if !isEncrypt {
			b.StopTimer()
		}

		if !ok {
			b.Fatalf("aeadDecrypt failed")
		}
		if !bytes.Equal(m, d) {
			b.Fatalf("aeadDecrypt output mismatch")
		}
	}
}

func init() {
	canAccelerate = IsHardwareAccelerated()
}
