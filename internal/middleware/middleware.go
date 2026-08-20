package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"url-shortener/internal/auth"
	"url-shortener/internal/helper"
)

type contextKey string

const UserIDKey contextKey = "userID"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		slog.Info("request started",
			"method", r.Method,
			"path", r.URL.Path,
		)

		next.ServeHTTP(w, r)

		slog.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)

	})
}

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
