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

func (u *urlU) Insert(value data.URL) (*data.URL, error) {
	var result data.URL

	stmt := u.sql.
		Insert(tableName).
		Columns("long_url", "short_url").
		Values(value.LongURL, value.ShortURL).
		Suffix("RETURNING id, long_url, short_url")

	err := u.db.Get(&result, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert url to db")
	}
	return &result, nil
}

func (u *urlU) UpdateShortCode(id int64, shortCode string) error {
	stmt := u.sql.
		Update(tableName).
		Set("short_url", shortCode).
		Where("id = ?", id)

	err := u.db.Exec(stmt)

	if err != nil {
		return errors.Wrap(err, "failed to update short_url in db")
	}
	return nil
}
