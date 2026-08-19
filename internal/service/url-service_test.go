package service

import (
	"context"
	"testing"
	"time"
	"url-shortener/internal/service/mocks"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenService_InvalidURL(t *testing.T) {

	s := &URLService{
		repo:  nil,
		redis: nil,
	}

	ctx := context.Background()

	_, err := s.ShortenService(ctx, "not-a-valid-url", 1)

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}

}

func TestShortenService_ExistingURL(t *testing.T) {

	mockRepo := new(mocks.URLRepositoryMock)

	mockRepo.On("CheckLongUrl", mock.Anything, "https://example.com", 1).
		Return("https://example.com", nil)
	mockRepo.On("TakeShortUrlRepo", mock.Anything, "https://example.com", 1).
		Return("abc123", nil)

	s := NewURLService(mockRepo, nil)

	result, err := s.ShortenService(context.Background(), "https://example.com", 1)

	assert.Nil(t, err)
	assert.Contains(t, result, "abc123")

}

func TestShortenService_NewURL(t *testing.T) {

	t.Setenv("BASE_URL", "https://short.com/")

	mockRepo := new(mocks.URLRepositoryMock)

	mockRepo.On("CheckLongUrl", mock.Anything, "https://new-example.com", 1).
		Return("", nil)
	mockRepo.On("CheckCode", mock.Anything, mock.Anything).
		Return(false, nil)
	mockRepo.On("SaveRepo", mock.Anything, "https://new-example.com", mock.Anything, 1).
		Return(nil)

	s := NewURLService(mockRepo, nil)

	result, err := s.ShortenService(context.Background(), "https://new-example.com", 1)

	assert.Nil(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "https://short.com/")

}

func TestRedirectService_CacheHit(t *testing.T) {

	mockRepo := new(mocks.URLRepositoryMock)
	mockRedis := new(mocks.RedisClientMock)

	mockRedis.On("Get", mock.Anything, "url:abc123").Return("https://example.com", nil)

	s := NewURLService(mockRepo, mockRedis)

	result, err := s.RedirestService(context.Background(), "abc123")

	assert.Nil(t, err)
	assert.Equal(t, "https://example.com", result)

	mockRepo.AssertNotCalled(t, "TakeOriginalUrlRepo", mock.Anything, "abc123")

}

func TestRedirectService_CacheMiss(t *testing.T) {

	mockRepo := new(mocks.URLRepositoryMock)
	mockRepo.On("TakeOriginalUrlRepo", mock.Anything, "abc123").
		Return("https://example.com", nil)

	mockRedis := new(mocks.RedisClientMock)
	mockRedis.On("Get", mock.Anything, "url:abc123").
		Return(nil, redis.Nil)
	mockRedis.On("Set", mock.Anything, "url:abc123", "https://example.com", 24*time.Hour).
		Return(nil)

	s := NewURLService(mockRepo, mockRedis)

	result, err := s.RedirestService(context.Background(), "abc123")

	assert.Nil(t, err)
	assert.Equal(t, "https://example.com", result)

	mockRepo.AssertCalled(t, "TakeOriginalUrlRepo", mock.Anything, "abc123")

	mockRedis.AssertCalled(t, "Set", mock.Anything, "url:abc123", "https://example.com", 24*time.Hour)

}
