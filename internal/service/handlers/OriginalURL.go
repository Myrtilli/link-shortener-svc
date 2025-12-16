package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func OriginalURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	logger := Log(r)
	db := DB(r)

	if strings.TrimSpace(shortCode) == "" {
		ape.RenderErr(w, problems.BadRequest(errors.New("Missing short code in URL path"))...)
		return
	}

	originalURL, err := db.URL().Get(shortCode)
	if err != nil {
		logger.WithError(err).Error("error: failed retrieve URL")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if originalURL == nil {
		logger.WithField("code", shortCode).Info("requested short code not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	http.Redirect(w, r, originalURL.LongURL, http.StatusTemporaryRedirect)
}
