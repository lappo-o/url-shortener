package service

import (
	"context"
	"testing"
	"url-shortener/internal/service/mocks"

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
