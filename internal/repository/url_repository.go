package repository

import (
	"context"

	"github.com/filipbabicdev/url-shortener/internal/base62"
	"github.com/filipbabicdev/url-shortener/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Create(ctx context.Context, originalURL string) (*model.URL, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	
	url := model.URL{
		OriginalURL: originalURL,
	}

	query := `INSERT INTO urls (original_url) VALUES ($1) RETURNING id, created_at`	
	err = tx.QueryRow(ctx, query, originalURL).Scan(&url.ID, &url.CreatedAt)
	if err != nil {
		return nil, err
	}

	url.ShortCode = base62.Encode(uint64(url.ID))
	updateQuery := `UPDATE urls SET short_code = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, updateQuery, url.ShortCode, url.ID)
	if err != nil {
		return nil, err
	}

	return &url, tx.Commit(ctx)
}

func (r *URLRepository) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	var url model.URL
	
	query := `SELECT id, original_url, short_code, created_at FROM urls WHERE short_code = $1`
	row := r.db.QueryRow(ctx, query, shortCode)
	
	err := row.Scan(&url.ID, &url.OriginalURL, &url.ShortCode, &url.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *URLRepository) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}