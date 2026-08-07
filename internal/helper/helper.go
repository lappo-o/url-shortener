//Helper for handle errors

package helper

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"url-shortener/internal/model"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrBadRequest      = errors.New("bad request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInternalServer  = errors.New("internal server error")
	ErrEmailExists     = errors.New("email already registered")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidPassword = errors.New("password must be between 6 and 72 characters")
	ErrInvalidToken    = errors.New("invalid or expired token")
)

func WriteJSON(w http.ResponseWriter, status int, data any, errMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.Response{
		Data:  data,
		Error: errMsg,
	})
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, ErrNotFound.Error()
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest, ErrBadRequest.Error()
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, ErrUnauthorized.Error()
	case errors.Is(err, ErrEmailExists):
		return http.StatusConflict, ErrEmailExists.Error()
	case errors.Is(err, ErrInvalidEmail), errors.Is(err, ErrInvalidPassword):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, ErrInternalServer.Error()
	}
}

func ErrorHandler(w http.ResponseWriter, err error) {
	status, errMsg := mapError(err)
	if status == http.StatusInternalServerError {
		log.Printf("internal error: %v", err)
	}
	WriteJSON(w, status, nil, errMsg)
}
