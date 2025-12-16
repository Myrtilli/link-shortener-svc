package shortening

import (
	"crypto/sha256"
	"strings"
)

const (
	alphabet        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base            = uint64(len(alphabet))
	hashPrefixBytes = 6
)

func GenerateShortKey(url string) string {
	hash := sha256.Sum256([]byte(url))

	var id uint64
	for i := 0; i < hashPrefixBytes; i++ {
		id = (id << 8) | uint64(hash[i])
	}

	return EncodeBase62(id)
}

func EncodeBase62(id uint64) string {
	if id == 0 {
		return string(alphabet[0])
	}

	buf := make([]byte, 11)
	i := len(buf)

	for id > 0 {
		i--
		buf[i] = alphabet[id%base]
		id /= base
	}

	return string(buf[i:])
}

func DecodeBase62(s string) uint64 {
	var id uint64
	for i := 0; i < len(s); i++ {
		idx := strings.IndexByte(alphabet, s[i])
		if idx == -1 {
			return 0
		}
		id = id*base + uint64(idx)
	}
	return id
}
