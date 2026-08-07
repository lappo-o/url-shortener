package handler

import (
	"encoding/json"
	"net/http"
	"url-shortener/internal/dto"
	"url-shortener/internal/helper"
	"url-shortener/internal/middleware"
	"url-shortener/internal/service"

	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(urlS *service.URLService) *URLHandler {
	return &URLHandler{service: urlS}
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {

	var req dto.ReqForShortenHandler

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helper.ErrorHandler(w, err)
		return
	}

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		helper.ErrorHandler(w, err)
		return
	}

	shortURL, err := h.service.ShortenService(r.Context(), req.URL, userID)
	if err != nil {
		helper.ErrorHandler(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, shortURL, "")

}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {

	code := chi.URLParam(r, "code")

	oUrl, err := h.service.RedirestService(r.Context(), code)
	if err != nil {
		helper.ErrorHandler(w, err)
		return
	}

	http.Redirect(w, r, oUrl, http.StatusFound)

}

func (h *URLHandler) ShowUrls(w http.ResponseWriter, r *http.Request) {

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		helper.ErrorHandler(w, err)
		return
	}

	urls, err := h.service.ShowMyUrls(r.Context(), userID)
	if err != nil {
		helper.ErrorHandler(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, urls, "")

}
