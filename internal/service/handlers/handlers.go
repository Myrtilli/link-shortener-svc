package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Myrtilli/link-shortener-svc/internal/service/dbmanagment"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

type URL struct {
	URL string `json:"url"`
}

var shortenService *Shorten

func InitShorten(repo *dbmanagment.Service) {
	shortenService = NewService(repo)
}

func URLHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleShorten(w, r)
		return
	case http.MethodGet:
		handleOriginalURL(w, r)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	var input URL

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Log(r).WithError(err).Error("failed to decode request body", "error", err)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if input.URL == "" {
		http.Error(w, "Missing 'url' field in request body", http.StatusBadRequest)
		return
	}

	ShortCode, err := shortenService.ShortenURL(input.URL)
	if err != nil {
		Log(r).WithError(err).Error("error: failed to shorten URL")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	type Response struct {
		OriginalURL string `json:"original_url"`
		ShortURL    string `json:"short_url"`
	}

	response := Response{
		OriginalURL: input.URL,
		ShortURL:    "http:///" + ShortCode,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

func handleOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")

	if err := json.NewDecoder(r.Body).Decode(&shortCode); err != nil {
		Log(r).WithError(err).Error("failed to decode request body", "error", err)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if shortCode == "" {
		http.Error(w, "Missing 'url' field in request body", http.StatusBadRequest)
		return
	}

	originalURL, err := shortenService.ResolveURL(shortCode)
	if err != nil {
		Log(r).WithError(err).Error("error: failed retrieve URL")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
