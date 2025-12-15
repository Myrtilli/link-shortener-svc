package data

type URLdb interface {
	Insert(value URL) (*URL, error)
	Get(shortURL string) (*URL, error)
	UpdateShortCode(id int64, shortCode string) error
}

type URL struct {
	ID       int64  `db:"id"`
	LongURL  string `db:"long_url"`
	ShortURL string `db:"short_url"`
}
