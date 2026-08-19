package service

import (
	"context"
	"net/url"
	"os"
	"time"
	"url-shortener/internal/helper"
	"url-shortener/internal/model"
)

type URLService struct {
	repo  URLRepository
	redis RedisClient
}

func NewURLService(urlR URLRepository, redis RedisClient) *URLService {
	return &URLService{
		repo:  urlR,
		redis: redis,
	}
}

func (s *URLService) ShortenService(ctx context.Context, URL string, userID int) (string, error) {

	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return "", helper.ErrBadRequest
	}

	url, err := s.repo.CheckLongUrl(ctx, URL, userID)
	if err != nil {
		return "", err
	}
	if url != "" {
		codeForUrl, err := s.repo.TakeShortUrlRepo(ctx, URL, userID)
		if err != nil {
			return "", err
		}
		return buildURL(codeForUrl), nil
	}

	codeForUrl, err := GenerateAndCheck(ctx, s)
	if codeForUrl == "" {
		return "", err
	}

	err = s.repo.SaveRepo(ctx, URL, codeForUrl, userID)
	if err != nil {
		return "", err
	}

	if s.redis != nil {
		s.redis.Set(ctx, "url:"+codeForUrl, URL, 24*time.Hour)
	}

	return buildURL(codeForUrl), nil

}

func buildURL(codeForUrl string) string {

	base := os.Getenv("BASE_URL")
	return base + codeForUrl

}

func (s *URLService) RedirestService(ctx context.Context, code string) (string, error) {

	if s.redis != nil {
		if cached, err := s.redis.Get(ctx, "url:"+code).Result(); err == nil {
			return cached, nil
		}
	}

	oUrl, err := s.repo.TakeOriginalUrlRepo(ctx, code)
	if err != nil {
		return "", err
	}

	if s.redis != nil {
		s.redis.Set(ctx, "url:"+code, oUrl, 24*time.Hour)
	}

	return oUrl, nil

}

func (s *URLService) ShowMyUrls(ctx context.Context, userID int) ([]model.Url, error) {

	urls, err := s.repo.GetAllURLs(ctx, userID)
	if err != nil {
		return nil, err
	}

	return urls, err

}
