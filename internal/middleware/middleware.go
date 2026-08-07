package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"url-shortener/internal/auth"
	"url-shortener/internal/helper"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		if header == "" {
			helper.ErrorHandler(w, helper.ErrUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")

		userID, err := auth.ValidateToken(tokenString)
		if err != nil {
			helper.ErrorHandler(w, helper.ErrUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))

	})

}

func GetUserID(ctx context.Context) (int, error) {

	userID, ok := ctx.Value(UserIDKey).(int)
	if !ok {
		return 0, errors.New("user id not found in context")
	}

	return userID, nil

}
