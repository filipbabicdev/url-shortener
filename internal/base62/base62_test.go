package base62

import "testing"

func TestEncodeDecode(t *testing.T) {
	testCases := []struct {
		name     string
		input    uint64
		expected string
	}{
		{"zero", 0, "0"},
		{"one", 1, "1"},
		{"sixty-one", 61, "z"},
		{"sixty-two", 62, "10"},
		{"large number", 1234567890, "1LY7VK"},
		{"max uint64", ^uint64(0), "LygHa16AHYF"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded := Encode(tc.input)
			if encoded != tc.expected {
				t.Errorf("Encode(%d) = %s; want %s", tc.input, encoded, tc.expected)
			}

			decoded, err := Decode(encoded)
			if err != nil {
				t.Errorf("Decode(%s) returned error: %v", encoded, err)
			}

			if decoded != tc.input {
				t.Errorf("Decode(%s) = %d; want %d", encoded, decoded, tc.input)
			}
		})
	}
}