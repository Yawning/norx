// +build !noasm,go1.10
// hwaccel_amd64.s - AMD64 optimized routines
//
// To the extent possible under law, Yawning Angel has waived all copyright
// and related or neighboring rights to the software, using the Creative
// Commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

#include "textflag.h"

// func cpuidAmd64(cpuidParams *uint32)
TEXT ·cpuidAmd64(SB), NOSPLIT, $0-8
	MOVQ cpuidParams+0(FP), R15
	MOVL 0(R15), AX
	MOVL 8(R15), CX
	CPUID
	MOVL AX, 0(R15)
	MOVL BX, 4(R15)
	MOVL CX, 8(R15)
	MOVL DX, 12(R15)
	RET

// func xgetbv0Amd64(xcrVec *uint32)
TEXT ·xgetbv0Amd64(SB), NOSPLIT, $0-8
	MOVQ xcrVec+0(FP), BX
	XORL CX, CX
	XGETBV
	MOVL AX, 0(BX)
	MOVL DX, 4(BX)
	RET

// Based heavily on the `ymm` referece implementation, but using assembly
// language instead of using intrinsics like in a sane language.
//
// The TWEAK_LOW_LATENCY variant is used for the permutation.

DATA ·vpshufb_idx_r0<>+0x00(SB)/8, $0x0007060504030201
DATA ·vpshufb_idx_r0<>+0x08(SB)/8, $0x080f0e0d0c0b0a09
DATA ·vpshufb_idx_r0<>+0x10(SB)/8, $0x0007060504030201
DATA ·vpshufb_idx_r0<>+0x18(SB)/8, $0x080f0e0d0c0b0a09
GLOBL ·vpshufb_idx_r0<>(SB), (NOPTR+RODATA), $32

DATA ·vpshufb_idx_r2<>+0x00(SB)/8, $0x0403020100070605
DATA ·vpshufb_idx_r2<>+0x08(SB)/8, $0x0c0b0a09080f0e0d
DATA ·vpshufb_idx_r2<>+0x10(SB)/8, $0x0403020100070605
DATA ·vpshufb_idx_r2<>+0x18(SB)/8, $0x0c0b0a09080f0e0d
GLOBL ·vpshufb_idx_r2<>(SB), (NOPTR+RODATA), $32

DATA ·tag_payload<>+0x00(SB)/8, $0x0000000000000000
DATA ·tag_payload<>+0x08(SB)/8, $0x0000000000000000
DATA ·tag_payload<>+0x10(SB)/8, $0x0000000000000000
DATA ·tag_payload<>+0x18(SB)/8, $0x0000000000000002
GLOBL ·tag_payload<>(SB), (NOPTR+RODATA), $32

DATA ·tag_final<>+0x00(SB)/8, $0x0000000000000000
DATA ·tag_final<>+0x08(SB)/8, $0x0000000000000000
DATA ·tag_final<>+0x10(SB)/8, $0x0000000000000000
DATA ·tag_final<>+0x18(SB)/8, $0x0000000000000008
GLOBL ·tag_final<>(SB), (NOPTR+RODATA), $32

#define G(A, B, C, D, T0, T1, R0, R2) \
	VPXOR   A, B, T0   \
	VPAND   A, B, T1   \
	VPADDQ  T1, T1, T1 \
	VPXOR   T0, T1, A  \
	VPXOR   D, T0, D   \
	VPXOR   D, T1, D   \
	VPSHUFB R0, D, D   \
	                   \
	VPXOR   C, D, T0   \
	VPAND   C, D, T1   \
	VPADDQ  T1, T1, T1 \
	VPXOR   T0, T1, C  \
	VPXOR   B, T0, B   \
	VPXOR   B, T1, B   \
	VPSRLQ  $19, B, T0 \
	VPSLLQ  $45, B, T1 \
	VPOR    T0, T1, B  \
	                   \
	VPXOR   A, B, T0   \
	VPAND   A, B, T1   \
	VPADDQ  T1, T1, T1 \
	VPXOR   T0, T1, A  \
	VPXOR   D, T0, D   \
	VPXOR   D, T1, D   \
	VPSHUFB R2, D, D   \
	                   \
	VPXOR   C, D, T0   \
	VPAND   C, D, T1   \
	VPADDQ  T1, T1, T1 \
	VPXOR   T0, T1, C  \
	VPXOR   B, T0, B   \
	VPXOR   B, T1, B   \
	VPADDQ  B, B, T0   \
	VPSRLQ  $63, B, T1 \
	VPOR    T0, T1, B

// -109 -> 147 (See: https://github.com/golang/go/issues/24378)
#define DIAGONALIZE(A, B, C, D) \
	VPERMQ $-109, D, D \
	VPERMQ $78, C, C   \
	VPERMQ $57, B, B

#define UNDIAGONALIZE(A, B, C, D) \
	VPERMQ $57, D, D   \
	VPERMQ $78, C, C   \
	VPERMQ $-109, B, B

// func initAVX2(s *uint64, key, nonce *byte, initConsts, instConsts *uint64)
TEXT ·initAVX2(SB), NOSPLIT, $0-40
	MOVQ s+0(FP), R8
	MOVQ key+8(FP), R9
	MOVQ nonce+16(FP), R10
	MOVQ initConsts+24(FP), R11
	MOVQ instConsts+32(FP), R12
	MOVQ 8(R12), AX

	VMOVDQU (R10), Y0
	VMOVDQU (R9), Y1
	VMOVDQU (R11), Y2
	VMOVDQU 32(R11), Y3

	VMOVDQU (R12), Y4
	VMOVDQA Y1, Y5

	VPXOR Y3, Y4, Y3

	VMOVDQU ·vpshufb_idx_r0<>(SB), Y13
	VMOVDQU ·vpshufb_idx_r2<>(SB), Y12

looprounds:
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	DIAGONALIZE(Y0, Y1, Y2, Y3)
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	UNDIAGONALIZE(Y0, Y1, Y2, Y3)
	SUBQ $1, AX
	JNZ  looprounds

	VPXOR Y3, Y5, Y3

	VMOVDQU Y0, (R8)
	VMOVDQU Y1, 32(R8)
	VMOVDQU Y2, 64(R8)
	VMOVDQU Y3, 96(R8)

	VZEROUPPER
	RET

// func absorbBlocksAVX2(s *uint64, in *byte, rounds, blocks uint64, tag *uint64)
TEXT ·absorbBlocksAVX2(SB), NOSPLIT, $0-40
	MOVQ s+0(FP), R8
	MOVQ in+8(FP), R10
	MOVQ rounds+16(FP), R11
	MOVQ blocks+24(FP), R12
	MOVQ tag+32(FP), R13

	VMOVDQU (R8), Y0
	VMOVDQU 32(R8), Y1
	VMOVDQU 64(R8), Y2
	VMOVDQU 96(R8), Y3

	VMOVDQU ·vpshufb_idx_r0<>(SB), Y13
	VMOVDQU ·vpshufb_idx_r2<>(SB), Y12
	VMOVDQU (R13), Y11

loopblocks:
	VPXOR Y3, Y11, Y3

	MOVQ R11, AX

looprounds:
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	DIAGONALIZE(Y0, Y1, Y2, Y3)
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	UNDIAGONALIZE(Y0, Y1, Y2, Y3)
	SUBQ $1, AX
	JNZ  looprounds

	VMOVDQU (R10), Y4
	VMOVDQU 32(R10), Y5
	VMOVDQU 64(R10), Y6

	VPXOR Y0, Y4, Y0
	VPXOR Y1, Y5, Y1
	VPXOR Y2, Y6, Y2

	VMOVDQU Y0, (R8)
	VMOVDQU Y1, 32(R8)
	VMOVDQU Y2, 64(R8)

	ADDQ $96, R10

	SUBQ $1, R12
	JNZ  loopblocks

	VMOVDQU Y3, 96(R8)

	VZEROUPPER
	RET

// func encryptBlocksAVX2(s *uint64, out, in *byte, rounds, blocks uint64)
TEXT ·encryptBlocksAVX2(SB), NOSPLIT, $0-40
	MOVQ s+0(FP), R8
	MOVQ out+8(FP), R9
	MOVQ in+16(FP), R10
	MOVQ rounds+24(FP), R11
	MOVQ blocks+32(FP), R12

	VMOVDQU (R8), Y0
	VMOVDQU 32(R8), Y1
	VMOVDQU 64(R8), Y2
	VMOVDQU 96(R8), Y3

	VMOVDQU ·vpshufb_idx_r0<>(SB), Y13
	VMOVDQU ·vpshufb_idx_r2<>(SB), Y12
	VMOVDQU ·tag_payload<>(SB), Y11

loopblocks:
	VPXOR Y3, Y11, Y3

	MOVQ R11, AX

looprounds:
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	DIAGONALIZE(Y0, Y1, Y2, Y3)
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	UNDIAGONALIZE(Y0, Y1, Y2, Y3)
	SUBQ $1, AX
	JNZ  looprounds

	VMOVDQU (R10), Y4
	VMOVDQU 32(R10), Y5
	VMOVDQU 64(R10), Y6

	VPXOR Y0, Y4, Y0
	VPXOR Y1, Y5, Y1
	VPXOR Y2, Y6, Y2

	VMOVDQU Y0, (R9)
	VMOVDQU Y1, 32(R9)
	VMOVDQU Y2, 64(R9)

	ADDQ $96, R9
	ADDQ $96, R10

	SUBQ $1, R12
	JNZ  loopblocks

	VMOVDQU Y0, (R8)
	VMOVDQU Y1, 32(R8)
	VMOVDQU Y2, 64(R8)
	VMOVDQU Y3, 96(R8)

	VZEROUPPER
	RET

// func decryptBlocksAVX2(s *uint64, out, in *byte, rounds, blocks uint64)
TEXT ·decryptBlocksAVX2(SB), NOSPLIT, $0-40
	MOVQ s+0(FP), R8
	MOVQ out+8(FP), R9
	MOVQ in+16(FP), R10
	MOVQ rounds+24(FP), R11
	MOVQ blocks+32(FP), R12

	VMOVDQU (R8), Y0
	VMOVDQU 32(R8), Y1
	VMOVDQU 64(R8), Y2
	VMOVDQU 96(R8), Y3

	VMOVDQU ·vpshufb_idx_r0<>(SB), Y13
	VMOVDQU ·vpshufb_idx_r2<>(SB), Y12
	VMOVDQU ·tag_payload<>(SB), Y11

loopblocks:
	VPXOR Y3, Y11, Y3

	MOVQ R11, AX

looprounds:
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	DIAGONALIZE(Y0, Y1, Y2, Y3)
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	UNDIAGONALIZE(Y0, Y1, Y2, Y3)
	SUBQ $1, AX
	JNZ  looprounds

	VMOVDQU (R10), Y4
	VMOVDQU 32(R10), Y5
	VMOVDQU 64(R10), Y6

	VPXOR Y0, Y4, Y0
	VPXOR Y1, Y5, Y1
	VPXOR Y2, Y6, Y2

	VMOVDQU Y0, (R9)
	VMOVDQU Y1, 32(R9)
	VMOVDQU Y2, 64(R9)

	VMOVDQA Y4, Y0
	VMOVDQA Y5, Y1
	VMOVDQA Y6, Y2

	ADDQ $96, R9
	ADDQ $96, R10

	SUBQ $1, R12
	JNZ  loopblocks

	VMOVDQU Y0, (R8)
	VMOVDQU Y1, 32(R8)
	VMOVDQU Y2, 64(R8)
	VMOVDQU Y3, 96(R8)

	VZEROUPPER
	RET

// func decryptLastBlockAVX2(s *uint64, out, in *byte, rounds, inLen uint64)
TEXT ·decryptLastBlockAVX2(SB), NOSPLIT, $0-40
	MOVQ s+0(FP), R8
	MOVQ out+8(FP), R9
	MOVQ in+16(FP), R10
	MOVQ rounds+24(FP), AX
	MOVQ inLen+32(FP), R12

	VMOVDQU (R8), Y0
	VMOVDQU 32(R8), Y1
	VMOVDQU 64(R8), Y2
	VMOVDQU 96(R8), Y3

	VMOVDQU ·vpshufb_idx_r0<>(SB), Y13
	VMOVDQU ·vpshufb_idx_r2<>(SB), Y12
	VMOVDQU ·tag_payload<>(SB), Y11

	VPXOR Y3, Y11, Y3

looprounds:
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	DIAGONALIZE(Y0, Y1, Y2, Y3)
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	UNDIAGONALIZE(Y0, Y1, Y2, Y3)
	SUBQ $1, AX
	JNZ  looprounds

	VMOVDQU Y0, (R9)
	VMOVDQU Y1, 32(R9)
	VMOVDQU Y2, 64(R9)

	CMPQ R12, $0
	JEQ  skipcopy
	XORQ AX, AX

loopcopy:
	MOVB (R10)(AX*1), BX
	MOVB BX, (R9)(AX*1)
	ADDQ $1, AX
	CMPQ AX, R12
	JNE  loopcopy

skipcopy:

	XORB $0x01, (R9)(R12*1)
	XORB $0x80, 95(R9)

	VMOVDQU (R9), Y4
	VMOVDQU 32(R9), Y5
	VMOVDQU 64(R9), Y6

	VPXOR Y0, Y4, Y0
	VPXOR Y1, Y5, Y1
	VPXOR Y2, Y6, Y2

	VMOVDQU Y0, (R9)
	VMOVDQU Y1, 32(R9)
	VMOVDQU Y2, 64(R9)

	VMOVDQU Y4, (R8)
	VMOVDQU Y5, 32(R8)
	VMOVDQU Y6, 64(R8)
	VMOVDQU Y3, 96(R8)

	VZEROUPPER
	RET

// func finalizeAVX2(s *uint64, out, key *byte, rounds uint64)
TEXT ·finalizeAVX2(SB), NOSPLIT, $0-32
	MOVQ s+0(FP), R8
	MOVQ out+8(FP), R9
	MOVQ key+16(FP), R10
	MOVQ rounds+24(FP), R11

	VMOVDQU (R8), Y0
	VMOVDQU 32(R8), Y1
	VMOVDQU 64(R8), Y2
	VMOVDQU 96(R8), Y3

	VMOVDQU ·vpshufb_idx_r0<>(SB), Y13
	VMOVDQU ·vpshufb_idx_r2<>(SB), Y12
	VMOVDQU ·tag_final<>(SB), Y11

	VPXOR Y3, Y11, Y3

	MOVQ R11, AX

looprounds:
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	DIAGONALIZE(Y0, Y1, Y2, Y3)
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	UNDIAGONALIZE(Y0, Y1, Y2, Y3)
	SUBQ $1, AX
	JNZ  looprounds

	VMOVDQU (R10), Y11
	VPXOR   Y3, Y11, Y3

looprounds2:
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	DIAGONALIZE(Y0, Y1, Y2, Y3)
	G(Y0, Y1, Y2, Y3, Y15, Y14, Y13, Y12)
	UNDIAGONALIZE(Y0, Y1, Y2, Y3)
	SUBQ $1, R11
	JNZ  looprounds2

	VPXOR   Y3, Y11, Y3
	VMOVDQU Y3, (R9)

	VZEROUPPER
	RET
