package base62

import ( 
	"fmt"; 
	"math"; 
	"strings" 
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(n uint64) string {
	if n == 0 {
		return string(base62Chars[0])
	}

	buf := make([]byte, 0, 11) // 62^11 > 2^64
	for n > 0 {
		buf = append(buf, base62Chars[n%62])
		n /= 62
	}

	// Reverse the buffer to get the correct order
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf)
}

func Decode(s string) (uint64, error) {
	if s == "" {
		return 0, fmt.Errorf("invalid input: empty string")
	}

	var result uint64
	for _, c := range s {
		idx := strings.IndexByte(base62Chars, byte(c))
		if idx < 0 {
			return 0, fmt.Errorf("invalid character: %c", c)
		}
		if result > (math.MaxUint64-uint64(idx)) / 62 {
			return 0, fmt.Errorf("invalid input: string represents a number too large for uint64")
		}

		result = result*62 + uint64(idx)
	}

	return result, nil
}

