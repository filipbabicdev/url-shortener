package repository

import (
	"context"
	"database/sql"

	"github.com/filip-hric/url-shortener/internal/base62"
	"github.com/filip-hric/url-shortener/internal/model"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Create(ctx context.Context, originalURL string) (*model.URL, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	
	url := model.URL{
		OriginalURL: originalURL,
	}

	query := `INSERT INTO urls (original_url) VALUES ($1) RETURNING id, created_at`	
	err = tx.QueryRowContext(ctx, query, originalURL).Scan(&url.ID, &url.CreatedAt)
	if err != nil {
		return nil, err
	}

	url.ShortCode = base62.Encode(uint64(url.ID))
	updateQuery := `UPDATE urls SET short_code = $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, updateQuery, url.ShortCode, url.ID)
	if err != nil {
		return nil, err
	}

	return &url, tx.Commit()
}

func (r *URLRepository) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	var url model.URL
	
	query := `SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = $1`
	row := r.db.QueryRowContext(ctx, query, shortCode)
	
	err := row.Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &url, nil
}