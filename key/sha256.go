package key

import (
	"encoding/binary"
)

const (
	blockSize  = 64
	digestSize = 32
)

type SHA256 struct {
	h    [8]uint32
	data []byte
	len  uint64
}

func NewSHA256() *SHA256 {
	return &SHA256{
		h: [8]uint32{
			0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
			0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
		},
		data: make([]byte, 0),
	}
}

func (s *SHA256) Write(p []byte) (int, error) {
	s.data = append(s.data, p...)
	s.len += uint64(len(p))
	return len(p), nil
}

func (s *SHA256) Sum() []byte {
	bitLen := s.len * 8
	padding := make([]byte, 0)
	if len(s.data)%blockSize < 56 {
		padding = make([]byte, 56-len(s.data)%blockSize)
	} else {
		padding = make([]byte, blockSize+56-len(s.data)%blockSize)
	}
	padding[0] = 0x80
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, bitLen)
	padding = append(padding, b...)

	s.Write(padding)

	for i := 0; i < len(s.data); i += blockSize {
		chunk := s.data[i : i+blockSize]
		w := [64]uint32{}
		for j := 0; j < 16; j++ {
			w[j] = binary.BigEndian.Uint32(chunk[j*4:])
		}
		for j := 16; j < 64; j++ {
			s0 := (w[j-15] >> 7) ^ (w[j-15] >> 18) ^ (w[j-15] << 14)
			s1 := (w[j-2] >> 17) ^ (w[j-2] >> 19) ^ (w[j-2] << 13)
			w[j] = w[j-16] + s0 + w[j-7] + s1
		}

		a, b, c, d, e, f, g, h := s.h[0], s.h[1], s.h[2], s.h[3], s.h[4], s.h[5], s.h[6], s.h[7]

		for j := 0; j < 64; j++ {
			S1 := (e >> 6) ^ (e >> 11) ^ (e >> 25)
			ch := (e & f) ^ ((^e) & g)
			temp1 := h + S1 + ch + k[j] + w[j]
			S0 := (a >> 2) ^ (a >> 13) ^ (a >> 22)
			maj := (a & b) ^ (a & c) ^ (b & c)
			temp2 := S0 + maj

			h, g, f, e = g, f, e, d+temp1
			d, c, b, a = c, b, a, temp1+temp2
		}

		s.h[0] += a
		s.h[1] += b
		s.h[2] += c
		s.h[3] += d
		s.h[4] += e
		s.h[5] += f
		s.h[6] += g
		s.h[7] += h
	}

	out := make([]byte, digestSize)
	for i := 0; i < 8; i++ {
		binary.BigEndian.PutUint32(out[i*4:], s.h[i])
	}
	return out
}

var k = [64]uint32{
	0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5,
	0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
	0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3,
	0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
	0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc,
	0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
	0x9beff9a3, 0xaec2083a, 0xbf597fc7, 0xc6e00bf3,
	0xd5a79147, 0x06ca6351, 0x14292967, 0x27b70a85,
	0x2e1b2138, 0x4d2c6dfc, 0x5ac42aed, 0x654be30e,
	0x766a0abb, 0x81c2c92e, 0x92722c85, 0xa2bfe8a1,
	0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819,
	0xd6990624, 0xf40e3585, 0x106aa070, 0x19a4c116,
	0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3,
	0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3, 0x748f82ee,
	0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa,
	0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
}
