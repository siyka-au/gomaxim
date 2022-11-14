package maxim

import "testing"

func TestCRC(t *testing.T) {
	testData := []struct {
		input    []byte
		expected byte
	}{
		{[]byte{0xa6, 0x00, 0x00}, 0x02},
		{[]byte{0xff, 0x23, 0x00}, 0x0c},
		{[]byte{0x01, 0x01, 0x00}, 0x13},
		{[]byte{0x01, 0x00, 0x00}, 0x05},
	}

	for _, testCase := range testData {
		crc, _ := CRC(testCase.input)
		if crc != testCase.expected {
			t.Errorf("CRC of (%x) was incorrect, got: %d, want: %d.", testCase.input, crc, testCase.expected)
		}
	}
}
