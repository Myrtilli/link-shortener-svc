package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Myrtilli/link-shortener-svc/internal/data"
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

	urlToInsert := data.URL{
		LongURL:  request.URL,
		ShortURL: "",
	}

	urlRepo := db.URL()

	insertedURL, err := urlRepo.Insert(urlToInsert)
	if err != nil {
		logger.WithError(err).Error("failed to insert URL into database")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	shortCode := shortening.EncodeBase62(insertedURL.ID)

	if err := urlRepo.UpdateShortCode(insertedURL.ID, shortCode); err != nil {
		logger.WithError(err).Error("failed to update short code in database")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	type Response struct {
		OriginalURL string `json:"original_url"`
		ShortURL    string `json:"short_url"`
	}

	response := Response{
		OriginalURL: insertedURL.LongURL,
		ShortURL:    "http://localhost:8000/integrations/link-shortener-svc/" + shortCode,
	}

	ape.Render(w, response)
}
