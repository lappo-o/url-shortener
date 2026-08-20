package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) SaveUser(
	ctx context.Context,
	email, pwH string,
) error {
	args := m.Called(ctx, email, pwH)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetUserByEmail(
	ctx context.Context,
	email string,
) (id int, passwordHash string, err error) {
	args := m.Called(ctx, email)
	return args.Int(0), args.String(1), args.Error(2)
}
