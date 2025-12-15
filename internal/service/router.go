package service

import (
	"github.com/Myrtilli/link-shortener-svc/internal/config"
	"github.com/Myrtilli/link-shortener-svc/internal/data/dblogic"
	"github.com/Myrtilli/link-shortener-svc/internal/service/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxDB(dblogic.NewMasterQ(cfg.DB())),
		),
	)
	r.Route("/integrations/link-shortener-svc", func(r chi.Router) {
		r.Get("/{code}", handlers.OriginalURL)
		r.Post("/", handlers.Shortcode)
	})

	return r
}
