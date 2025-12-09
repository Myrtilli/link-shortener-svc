package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	ctxKeyDB  ctxKey = iota
	logCtxKey ctxKey = iota
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func DBconnect(ctx context.Context) *sql.DB {
	if v := ctx.Value(ctxKeyDB); v != nil {
		if db, ok := v.(*sql.DB); ok {
			return db
		}
	}
	return nil
}

func DBStore(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, ctxKeyDB, db)
}
