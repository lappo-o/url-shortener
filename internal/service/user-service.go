package service

import (
	"context"
	"errors"
	"log"
	"net/mail"
	"url-shortener/internal/auth"
	"url-shortener/internal/helper"
	"url-shortener/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  *repository.UserRepository
	redis *redis.Client
}

func NewUserService(repo *repository.UserRepository, redis *redis.Client) *UserService {
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
		return helper.ErrInvalidEmail
	}

	if len(pw) < 6 || len(pw) > 72 {
		return helper.ErrInvalidPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("bcrypt generate err: %v", err)
		return errors.New("internal server error")
	}

	return s.repo.SaveUser(ctx, email, string(hashedPassword))

}

func (s *UserService) LoginService(
	ctx context.Context,
	email, pw string,
) (string, error) {

	userID, hash, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", helper.ErrUnauthorized
		}
		log.Printf("get password hash failed: %v", err)
		return "", helper.ErrInternalServer
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)); err != nil {
		return "", helper.ErrUnauthorized
	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		log.Printf("generate token failed: %v", err)
		return "", helper.ErrInternalServer
	}

	return token, nil

}
