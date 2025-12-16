package handlers

import (
	"context"
	"net/http"

	"github.com/Myrtilli/link-shortener-svc/internal/config"
	"github.com/Myrtilli/link-shortener-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	ctxKeyDB   ctxKey = iota
	logCtxKey  ctxKey = iota
	baseCtxKey ctxKey = iota
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxDB(entry data.MasterQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxKeyDB, entry)
	}
}

func DB(r *http.Request) data.MasterQ {
	return r.Context().Value(ctxKeyDB).(data.MasterQ).New()
}

func CtxBase(cfg config.Config) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, baseCtxKey, cfg)
	}
}

func Base(r *http.Request) config.Config {
	return r.Context().Value(baseCtxKey).(config.Config)
}
