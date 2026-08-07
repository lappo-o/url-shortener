package repository

import (
	"context"
	"errors"
	"fmt"
	"url-shortener/internal/helper"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) SaveUser(
	ctx context.Context,
	email, pwH string,
) error {

	_, err := r.pool.Exec(ctx,
		"INSERT INTO users (email, password_hash) VALUES ($1, $2)",
		email, pwH,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return helper.ErrEmailExists
		}
		return fmt.Errorf("SaveUser: %w", err)
	}

	return nil

}

func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (id int, passwordHash string, err error) {

	err = r.pool.QueryRow(ctx,
		"SELECT id, password_hash FROM users WHERE email = $1",
		email,
	).Scan(&id, &passwordHash)
	if err != nil {
		return 0, "", err
	}

	return id, passwordHash, nil

}
