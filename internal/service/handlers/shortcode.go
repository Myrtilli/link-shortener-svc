package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Myrtilli/link-shortener-svc/internal/data"
	"github.com/Myrtilli/link-shortener-svc/internal/shortening"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func Shortcode(w http.ResponseWriter, r *http.Request) {

	type URL struct {
		URL string `json:"url"`
	}
	var input URL

	logger := Log(r)
	db := DB(r)

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.WithError(err).Error("failed to decode request body", "error", err)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if input.URL == "" {
		ape.RenderErr(w, problems.BadRequest(errors.New("Missing 'url' field in request body"))...)
		return
	}

	urlToInsert := data.URL{
		LongURL:  input.URL,
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
