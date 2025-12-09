package service

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("url not found")

func (s *service) ShortenURL(original string) (string, error) {
	var id int64
	err := s.DB.QueryRow(`
        INSERT INTO short_urls (long_url)
        VALUES ($1)
        RETURNING id
    `, original).Scan(&id)
	if err != nil {
		return "", err
	}

	code := encodeBase62(id)
	return code, nil
}

// ResolveURL returns the original URL from a Base62 code
func (s *service) ResolveURL(code string) (string, error) {
	id := decodeBase62(code)
	if id == 0 {
		return "", ErrNotFound
	}

	var longURL string
	err := s.DB.QueryRow(`
        SELECT long_url FROM short_urls WHERE id = $1
    `, id).Scan(&longURL)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return longURL, nil
}

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base = int64(len(alphabet))

func encodeBase62(id int64) string {
	if id == 0 {
		return string(alphabet[0])
	}

	var result []byte
	for id > 0 {
		result = append(result, alphabet[id%base])
		id /= base
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

func decodeBase62(s string) int64 {
	var id int64 = 0
	for _, char := range s {
		index := -1
		for i, alphabetLetter := range alphabet {
			if alphabetLetter == char {
				index = i
				break
			}
		}
		if index == -1 {
			return 0
		}
		id = id*base + int64(index)
	}
	return id
}
