package dblogic

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/Myrtilli/link-shortener-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const tableName = "short_urls"

func newURLdb(db *pgdb.DB) data.URLdb {
	return &urlU{
		db:  db,
		sql: sq.StatementBuilder,
	}
}

type urlU struct {
	db  *pgdb.DB
	sql sq.StatementBuilderType
}

func (u *urlU) Get(shortURL string) (*data.URL, error) {
	var result data.URL

	stmt := u.sql.
		Select("id", "long_url", "short_url").
		From(tableName).
		Where("short_url = ?", shortURL)

	err := u.db.Get(&result, stmt)

	if errors.Cause(err) == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get url from db")
	}

	return &result, nil
}

func (u *urlU) Insert(url data.URL) (*data.URL, error) {
	query := sq.Insert("short_urls").
		Columns("long_url", "short_url").
		Values(url.LongURL, url.ShortURL).
		Suffix("ON CONFLICT (short_url) DO UPDATE SET long_url = EXCLUDED.long_url RETURNING id")

	var id int64
	err := u.db.Get(&id, query)
	if err != nil {
		return nil, err
	}

	url.ID = id
	return &url, nil
}
