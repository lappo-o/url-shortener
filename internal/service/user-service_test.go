package service

import (
	"context"
	"testing"
	"url-shortener/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterService_InvalidEmail(t *testing.T) {

	s := &UserService{
		repo:  nil,
		redis: nil,
	}

	ctx := context.Background()

	err := s.RegisterService(ctx, "not-a-valid-email", "qwerty123456")

	if err == nil {
		t.Error("expected error for invalid email, got nil")
	}

}

func TestRegisterService_InvalidPassword(t *testing.T) {

	s := &UserService{
		repo:  nil,
		redis: nil,
	}

	ctx := context.Background()

	err := s.RegisterService(ctx, "user@email.com", "123")

	if err == nil {
		t.Error("expected error for short password, got nil")
	}

}

func TestRegisterService_ValidRegistration(t *testing.T) {

	mockRepo := new(mocks.UserRepositoryMock)

	mockRepo.On("SaveUser", mock.Anything, "user@email.com", mock.Anything).
		Return(nil)

	s := NewUserService(mockRepo, nil)

	err := s.RegisterService(context.Background(), "user@email.com", "validpass123")

	assert.Nil(t, err)
	mockRepo.AssertCalled(t, "SaveUser", mock.Anything, "user@email.com", mock.Anything)

}
