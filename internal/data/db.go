package data

type URLdb interface {
	Insert(URL) (*URL, error)
	Get(shortURL string) (*URL, error)
}

type URL struct {
	ID       int64  `db:"id"`
	LongURL  string `db:"long_url"`
	ShortURL string `db:"short_url"`
}
