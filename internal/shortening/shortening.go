package shortening

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

const (
	alphabet        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base            = uint64(len(alphabet))
	shortKeyLength  = 8
	hashPrefixBytes = 6
)

func GenerateShortKey(url string) string {
	hash := sha256.Sum256([]byte(url))

	var id uint64
	for i := 0; i < hashPrefixBytes; i++ {
		id = (id << 8) | uint64(hash[i])
	}

	return EncodeBase62(id, shortKeyLength)
}

func EncodeBase62(id uint64, length int) string {
	buf := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		buf[i] = alphabet[id%base]
		id /= base
	}
	return string(buf)
}

func DecodeBase62(s string) (uint64, error) {
	var id uint64
	for i := 0; i < len(s); i++ {
		idx := strings.IndexByte(alphabet, s[i])
		if idx == -1 {
			return 0, fmt.Errorf("invalid base62 character: %q", s[i])
		}
		id = id*base + uint64(idx)
	}
	return id, nil
}
