package shortly

import (
	"fmt"
	"strings"
)

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func EncodeBase62(num int64) string {
	if num == 0 {
		return "0"
	}

	base := int64(len(base62Chars))
	buf := make([]byte, 0, 11)

	for num > 0 {
		buf = append(buf, base62Chars[num%base])
		num /= base
	}

	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf)
}

func DecodeBase62(s string) (int64, error) {
	var result int64
	base := int64(len(base62Chars))

	for _, c := range s {
		idx := strings.IndexRune(base62Chars, c)
		if idx == -1 {
			return 0, fmt.Errorf("invalid character: %c", c)
		}
		result = result*base + int64(idx)
	}

	return result, nil
}
