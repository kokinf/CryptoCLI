package key

import "bytes"

func PBKDF2(password, salt []byte, iterations, dkLen int) []byte {
	var dk bytes.Buffer
	block := 1
	for dk.Len() < dkLen {
		u := HmacSha256(password, append(salt, IntToFourBytes(block)...))
		for i := 1; i < iterations; i++ {
			u = HmacSha256(password, u)
		}
		dk.Write(u)
		block++
	}
	return dk.Bytes()[:dkLen]
}

func IntToFourBytes(i int) []byte {
	b := make([]byte, 4)
	b[0] = byte(i >> 24)
	b[1] = byte(i >> 16)
	b[2] = byte(i >> 8)
	b[3] = byte(i)
	return b
}
