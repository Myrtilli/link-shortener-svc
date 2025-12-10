package dbmanagment

import (
	"database/sql"
	"errors"
)

type Service struct {
	DB *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{DB: db}
}

var ErrNotFound = errors.New("url not found")

func (r *Service) InsertURL(longURL string) (int64, error) {
	var id int64
	err := r.DB.QueryRow(`INSERT INTO short_urls (long_url) VALUES ($1) RETURNING id`, longURL).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Service) GetURL(id int64) (string, error) {
	var longURL string
	err := r.DB.QueryRow(`SELECT long_url FROM short_urls WHERE id = $1`, id).Scan(&longURL)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return longURL, nil
}
