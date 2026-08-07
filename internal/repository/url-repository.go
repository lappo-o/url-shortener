package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/helper"
	"url-shortener/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	pool *pgxpool.Pool
}

func NewURLRepository(pool *pgxpool.Pool) *URLRepository {
	return &URLRepository{pool: pool}
}

func (r *URLRepository) SaveRepo(ctx context.Context, originalURL, shortURL string, userID int) error {

	_, err := r.pool.Exec(ctx,
		"INSERT INTO urls (original_url, short_url, user_id) VALUES ($1, $2, $3)",
		originalURL, shortURL, userID,
	)

	return err

}

func (r *URLRepository) CheckLongUrl(ctx context.Context, originalURL string, userID int) (string, error) {

	var oUrl string

	err := r.pool.QueryRow(ctx,
		"SELECT original_url FROM urls WHERE original_url = $1 AND user_id = $2",
		originalURL, userID,
	).Scan(&oUrl)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return oUrl, nil

}

func (r *URLRepository) TakeShortUrlRepo(ctx context.Context, originalURl string, userID int) (string, error) {

	var sUrl string

	err := r.pool.QueryRow(ctx,
		"SELECT short_url FROM urls WHERE original_url = $1 AND user_id = $2",
		originalURl, userID,
	).Scan(&sUrl)

	if err != nil {
		return "", err
	}

	return sUrl, nil

}

func (r *URLRepository) TakeOriginalUrlRepo(ctx context.Context, code string) (string, error) {

	var oUrl string

	err := r.pool.QueryRow(ctx,
		"SELECT original_url FROM urls WHERE short_url = $1",
		code,
	).Scan(&oUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", helper.ErrNotFound
		}
		return "", err
	}

	return oUrl, nil

}

func (r *URLRepository) CheckCode(ctx context.Context, code string) (bool, error) {

	var str string

	err := r.pool.QueryRow(ctx,
		"SELECT short_url FROM urls WHERE short_url = $1",
		code,
	).Scan(&str)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil

}

func (r *URLRepository) GetAllURLs(ctx context.Context, userID int) ([]model.Url, error) {

	query := "SELECT original_url, short_url, created_at FROM urls WHERE user_id = $1 ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query urls: %w", err)
	}
	defer rows.Close()

	var urls []model.Url

	for rows.Next() {

		var u model.Url

		if err := rows.Scan(&u.OriginalURL, &u.ShortCode, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan url: %w", err)
		}

		urls = append(urls, u)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return urls, nil

}
