package service

import (
	"context"
	"url-shortener/internal/model"
)

type URLRepository interface {
	SaveRepo(ctx context.Context, originalURL, shortURL string, userID int) error
	CheckLongUrl(ctx context.Context, originalURL string, userID int) (string, error)
	TakeShortUrlRepo(ctx context.Context, originalURl string, userID int) (string, error)
	TakeOriginalUrlRepo(ctx context.Context, code string) (string, error)
	CheckCode(ctx context.Context, code string) (bool, error)
	GetAllURLs(ctx context.Context, userID int) ([]model.Url, error)
}
