package handler

import (
	"encoding/json"
	"net/http"
	"url-shortener/internal/dto"
	"url-shortener/internal/helper"
	"url-shortener/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var req dto.ReqForAuth

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ErrorHandler(w, err)
		return
	}

	err := h.service.RegisterService(r.Context(), req.Email, req.Password)
	if err != nil {
		helper.ErrorHandler(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusCreated, "successful registration", "")

}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var req dto.ReqForAuth

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.ErrorHandler(w, err)
		return
	}

	token, err := h.service.LoginService(r.Context(), req.Email, req.Password)
	if err != nil {
		helper.ErrorHandler(w, err)
		return
	}

	helper.WriteJSON(w, http.StatusOK, token, "")

}
