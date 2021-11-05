package phash

func CountBitsInByte(x byte) int {
	x = (x & 0x55) + ((x >> 1) & 0x55)
	x = (x & 0x33) + ((x >> 2) & 0x33)
	return int((x & 0x0f) + ((x >> 4) & 0x0f))
}

func Hamming(a, b []byte) int {
	d := 0
	for i, x := range a {
		d += CountBitsInByte(x ^ b[i])
	}
	return d
}
