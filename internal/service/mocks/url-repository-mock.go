package mocks

import (
	"context"
	"url-shortener/internal/model"

	"github.com/stretchr/testify/mock"
)

type URLRepositoryMock struct {
	mock.Mock
}

func (m *URLRepositoryMock) SaveRepo(
	ctx context.Context,
	originalURL, shortURL string,
	userID int,
) error {
	args := m.Called(ctx, originalURL, shortURL, userID)
	return args.Error(0)
}

func (m *URLRepositoryMock) CheckLongUrl(
	ctx context.Context,
	originalURL string,
	userID int,
) (string, error) {
	args := m.Called(ctx, originalURL, userID)
	return args.String(0), args.Error(1)
}

func (m *URLRepositoryMock) TakeShortUrlRepo(
	ctx context.Context,
	originalURl string,
	userID int,
) (string, error) {
	args := m.Called(ctx, originalURl, userID)
	return args.String(0), args.Error(1)
}

func (m *URLRepositoryMock) TakeOriginalUrlRepo(
	ctx context.Context,
	code string,
) (string, error) {
	args := m.Called(ctx, code)
	return args.String(0), args.Error(1)
}

func (m *URLRepositoryMock) CheckCode(
	ctx context.Context,
	code string,
) (bool, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Error(1)
}

func (m *URLRepositoryMock) GetAllURLs(
	ctx context.Context,
	userID int,
) ([]model.Url, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Url), args.Error(1)
}
