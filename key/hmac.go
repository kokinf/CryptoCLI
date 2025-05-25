package key

func HmacSha256(key, data []byte) []byte {
	opad := make([]byte, 64)
	ipad := make([]byte, 64)
	if len(key) > 64 {
		hash := NewSHA256()
		hash.Write(key)
		key = hash.Sum()
	}
	copy(opad, key)
	copy(ipad, key)
	for i := range opad {
		opad[i] ^= 0x5c
	}
	for i := range ipad {
		ipad[i] ^= 0x36
	}

	hash1 := NewSHA256()
	hash1.Write(ipad)
	hash1.Write(data)
	inner := hash1.Sum()

	hash2 := NewSHA256()
	hash2.Write(opad)
	hash2.Write(inner)
	return hash2.Sum()
}
