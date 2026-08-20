package service

import (
	"context"
	"errors"
	"log/slog"
	"net/mail"
	"url-shortener/internal/auth"
	"url-shortener/internal/helper"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  UserRepository
	redis RedisClient
}

func NewUserService(repo UserRepository, redis RedisClient) *UserService {
	return &UserService{
		repo:  repo,
		redis: redis,
	}
}

func (s *UserService) RegisterService(
	ctx context.Context,
	email, pw string,
) error {

	_, err := mail.ParseAddress(email)
	if err != nil {
		slog.Warn("invalid email format", "email", email)
		return helper.ErrInvalidEmail
	}

	if len(pw) < 6 || len(pw) > 72 {
		slog.Warn("invalid password length", "email", email)
		return helper.ErrInvalidPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("bcrypt generate failed", "error", err)
		return errors.New("internal server error")
	}

	slog.Info("user registered", "email", email)
	return s.repo.SaveUser(ctx, email, string(hashedPassword))

}

func (s *UserService) LoginService(
	ctx context.Context,
	email, pw string,
) (string, error) {

	userID, hash, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("login failed: user not found", "email", email)
			return "", helper.ErrUnauthorized
		}
		slog.Error("GetUserByEmail failed", "error", err)
		return "", helper.ErrInternalServer
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)); err != nil {
		slog.Warn("login failed: wrong password", "email", email)
		return "", helper.ErrUnauthorized
	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		slog.Error("generate token failed", "error", err, "userID", userID)
		return "", helper.ErrInternalServer
	}

	slog.Info("user logged in", "userID", userID, "email", email)
	return token, nil

}
