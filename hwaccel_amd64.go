// hwaccel_amd64.go - AMD64 optimized routines
//
// To the extent possible under law, Yawning Angel has waived all copyright
// and related or neighboring rights to the software, using the Creative
// Commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

// +build amd64,!gccgo,!noasm,go1.10

package norx

//go:noescape
func cpuidAmd64(cpuidParams *uint32)

//go:noescape
func xgetbv0Amd64(xcrVec *uint32)

//go:noescape
func initAVX2(s *uint64, key, nonce *byte, initConsts, instConsts *uint64)

//go:noescape
func absorbBlocksAVX2(s *uint64, in *byte, rounds, blocks uint64, tag *uint64)

//go:noescape
func encryptBlocksAVX2(s *uint64, out, in *byte, rounds, blocks uint64)

//go:noescape
func decryptBlocksAVX2(s *uint64, out, in *byte, rounds, blocks uint64)

//go:noescape
func decryptLastBlockAVX2(s *uint64, out, in *byte, rounds, inLen uint64)

//go:noescape
func finalizeAVX2(s *uint64, out, key *byte, rounds uint64)

func supportsAVX2() bool {
	// https://software.intel.com/en-us/articles/how-to-detect-new-instruction-support-in-the-4th-generation-intel-core-processor-family
	const (
		osXsaveBit = 1 << 27
		avx2Bit    = 1 << 5
	)

	// Check to see if CPUID actually supports the leaf that indicates AVX2.
	// CPUID.(EAX=0H, ECX=0H) >= 7
	regs := [4]uint32{0x00}
	cpuidAmd64(&regs[0])
	if regs[0] < 7 {
		return false
	}

	// Check to see if the OS knows how to save/restore XMM/YMM state.
	// CPUID.(EAX=01H, ECX=0H):ECX.OSXSAVE[bit 27]==1
	regs = [4]uint32{0x01}
	cpuidAmd64(&regs[0])
	if regs[2]&osXsaveBit == 0 {
		return false
	}
	xcrRegs := [2]uint32{}
	xgetbv0Amd64(&xcrRegs[0])
	if xcrRegs[0]&6 != 6 {
		return false
	}

	// Check for AVX2 support.
	// CPUID.(EAX=07H, ECX=0H):EBX.AVX2[bit 5]==1
	regs = [4]uint32{0x07}
	cpuidAmd64(&regs[0])
	return regs[1]&avx2Bit != 0
}

var implAVX2 = &hwaccelImpl{
	name:          "AVX2",
	initFn:        initYMM,
	absorbDataFn:  absorbDataYMM,
	encryptDataFn: encryptDataYMM,
	decryptDataFn: decryptDataYMM,
	finalizeFn:    finalizeYMM,
}

func initYMM(s *state, key, nonce []byte) {
	var instConsts = [4]uint64{paramW, uint64(s.rounds), paramP, paramT}
	initAVX2(&s.s[0], &key[0], &nonce[0], &initializationConstants[8], &instConsts[0])
}

func absorbDataYMM(s *state, in []byte, tag uint64) {
	inLen := len(in)
	if inLen == 0 {
		return
	}

	var tagVec = [4]uint64{0, 0, 0, tag}
	var off int
	if inBlocks := inLen / bytesR; inBlocks > 0 {
		absorbBlocksAVX2(&s.s[0], &in[0], uint64(s.rounds), uint64(inBlocks), &tagVec[0])
		off += inBlocks * bytesR
	}
	in = in[off:]

	var lastBlock [bytesR]byte
	padRef(&lastBlock, in)
	absorbBlocksAVX2(&s.s[0], &lastBlock[0], uint64(s.rounds), 1, &tagVec[0])
}

func encryptDataYMM(s *state, out, in []byte) {
	inLen := len(in)
	if inLen == 0 {
		return
	}

	var off int
	if inBlocks := inLen / bytesR; inBlocks > 0 {
		encryptBlocksAVX2(&s.s[0], &out[0], &in[0], uint64(s.rounds), uint64(inBlocks))
		off += inBlocks * bytesR
	}
	out, in = out[off:], in[off:]

	var lastBlock [bytesR]byte
	padRef(&lastBlock, in)
	encryptBlocksAVX2(&s.s[0], &lastBlock[0], &lastBlock[0], uint64(s.rounds), 1)
	copy(out, lastBlock[:len(in)])
}

func decryptDataYMM(s *state, out, in []byte) {
	inLen := len(in)
	if inLen == 0 {
		return
	}

	var off int
	if inBlocks := inLen / bytesR; inBlocks > 0 {
		decryptBlocksAVX2(&s.s[0], &out[0], &in[0], uint64(s.rounds), uint64(inBlocks))
		off += inBlocks * bytesR
	}
	out, in = out[off:], in[off:]

	var lastBlock [bytesR]byte
	var inPtr *byte
	if len(in) != 0 {
		inPtr = &in[0]
	}
	decryptLastBlockAVX2(&s.s[0], &lastBlock[0], inPtr, uint64(s.rounds), uint64(len(in)))
	copy(out, lastBlock[:len(in)])
	burnBytes(lastBlock[:])
}

func finalizeYMM(s *state, tag, key []byte) {
	var lastBlock [bytesC]byte

	finalizeAVX2(&s.s[0], &lastBlock[0], &key[0], uint64(s.rounds))
	copy(tag, lastBlock[:bytesT])
	burnBytes(lastBlock[:]) // burn buffer
	burnUint64s(s.s[:])     // at this point we can also burn the state
}

func initHardwareAcceleration() {
	if supportsAVX2() {
		isHardwareAccelerated = true
		hardwareAccelImpl = implAVX2
	}
}
