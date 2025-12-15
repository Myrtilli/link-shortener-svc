package dblogic

import (
	"github.com/Myrtilli/link-shortener-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

func NewMasterQ(db *pgdb.DB) data.MasterQ {
	return &masterQ{
		db: db.Clone(),
	}
}

type masterQ struct {
	db *pgdb.DB
}

func (m *masterQ) New() data.MasterQ {
	return NewMasterQ(m.db)
}

func (m *masterQ) URL() data.URLdb {
	return newURLdb(m.db)
}
