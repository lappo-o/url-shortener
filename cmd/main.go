package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"url-shortener/internal/db"
	"url-shortener/internal/handler"
	"url-shortener/internal/middleware"
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
		panic(err)
	}
	defer pool.Close()

	fmt.Println("DB shortener connected")

	r := chi.NewRouter()

	urlRepo := repository.NewURLRepository(pool)
	urlService := service.NewURLService(urlRepo)
	urlHandler := handler.NewURLHandler(urlService)

	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r.Post("/register", userHandler.RegisterHandler)
	r.Post("/login", userHandler.LoginHandler)
	r.Get("/{code}", urlHandler.Redirect)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/short", urlHandler.Shorten)
		r.Get("/urls", urlHandler.ShowUrls)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
