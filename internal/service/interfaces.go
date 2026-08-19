package service

import (
	"context"
	"time"
	"url-shortener/internal/model"

	"github.com/redis/go-redis/v9"
)

type URLRepository interface {
	SaveRepo(ctx context.Context, originalURL, shortURL string, userID int) error
	CheckLongUrl(ctx context.Context, originalURL string, userID int) (string, error)
	TakeShortUrlRepo(ctx context.Context, originalURl string, userID int) (string, error)
	TakeOriginalUrlRepo(ctx context.Context, code string) (string, error)
	CheckCode(ctx context.Context, code string) (bool, error)
	GetAllURLs(ctx context.Context, userID int) ([]model.Url, error)
}

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}
