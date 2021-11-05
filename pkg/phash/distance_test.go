package phash

import "testing"

func TestCountBitsInByte(t *testing.T) {
	x := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}
	actualNArr := []int{0, 1, 1, 2, 8}
	for ind, c := range x {
		if actualN := CountBitsInByte(c); actualN != actualNArr[ind] {
			t.Errorf("%x -> (actual) %d != %d (expected)", c, actualN, actualNArr[ind])
		}
	}
}

func TestHamming(t *testing.T) {
	x1 := []byte{0x11, 0xFF}
	x2 := []byte{0x01, 0xF3}
	expected := 3
	actual := Hamming(x1, x2)

	if actual != expected {
		t.Errorf("Hamming for x1:%x , x2:%x , (actual) %d != %d (expected)", x1, x2, actual, expected)
	}

}
