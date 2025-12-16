package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Myrtilli/link-shortener-svc/internal/data"
	"github.com/Myrtilli/link-shortener-svc/internal/service/models"
	"github.com/Myrtilli/link-shortener-svc/internal/service/requests"
	"github.com/Myrtilli/link-shortener-svc/internal/shortening"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Shortcode(w http.ResponseWriter, r *http.Request) {
	var request requests.CreateShortLinkRequest

	logger := Log(r)
	db := DB(r)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.WithError(err).Error("failed to decode request body", "error", err)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := request.Validate(); err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	shortCode := shortening.GenerateShortKey(request.URL)

	urlToInsert := data.URL{
		LongURL:  request.URL,
		ShortURL: shortCode,
	}

	insertedURL, err := db.URL().Insert(urlToInsert)
	if err != nil {
		logger.WithError(err).Error("failed to insert URL")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := models.ShortLinkResponse{
		OriginalURL: insertedURL.LongURL,
		ShortURL:    "http://localhost:8000/integrations/link-shortener-svc/" + shortCode,
	}

	ape.Render(w, response)
}
