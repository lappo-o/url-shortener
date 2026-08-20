package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"url-shortener/internal/db"
	"url-shortener/internal/handler"
	"url-shortener/internal/middleware"
	"url-shortener/internal/rdb"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	levelStr := os.Getenv("LOG_LEVEL")

	var level slog.Level

	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	if os.Getenv("ENV") == "production" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})))
	}

	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, using system environment variables.")
	}

	ctx := context.Background()

	connString := os.Getenv("DB_URL")

	pool, err := db.NewPostgresPool(ctx, connString)
	if err != nil {
		slog.Error("postgres connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("PostgreSQL connected")

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbNum, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		slog.Error("invalid REDIS_DB", "error", err)
		os.Exit(1)
	}

	rdClient, err := rdb.NewRedis(addr, password, dbNum)
	if err != nil {
		slog.Error("Create a Redis client failed", "error", err)
		os.Exit(1)
	}
	defer rdClient.Close()
	slog.Info("Redis connected")

	r := chi.NewRouter()

	urlRepo := repository.NewURLRepository(pool)
	urlService := service.NewURLService(urlRepo, rdClient)
	urlHandler := handler.NewURLHandler(urlService)

	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo, rdClient)
	userHandler := handler.NewUserHandler(userService)

	r.Use(middleware.LoggingMiddleware)

	r.Post("/register", userHandler.RegisterHandler)
	r.Post("/login", userHandler.LoginHandler)
	r.Get("/{code}", urlHandler.Redirect)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/short", urlHandler.Shorten)
		r.Get("/urls", urlHandler.ShowUrls)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		slog.Info("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server startup failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited")

}
