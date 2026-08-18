package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	ctx := context.Background()

	connString := os.Getenv("DB_URL")

	pool, err := db.NewPostgresPool(ctx, connString)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()
	fmt.Println("PostgreSQL connected")

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbNum, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	rdClient, err := rdb.NewRedis(addr, password, dbNum)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer rdClient.Close()
	fmt.Println("Redis connected")

	r := chi.NewRouter()

	urlRepo := repository.NewURLRepository(pool)
	urlService := service.NewURLService(urlRepo, rdClient)
	urlHandler := handler.NewURLHandler(urlService)

	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo, rdClient)
	userHandler := handler.NewUserHandler(userService)

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
		fmt.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}

	fmt.Println("Server exited")

}
