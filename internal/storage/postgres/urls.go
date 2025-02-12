package postgresql

import (
	"context"
	_ "database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	_ "grpc-service-ref/internal/domain/models"
	"grpc-service-ref/internal/storage"
	e "grpc-service-ref/pkg"
)

const (
	uniqueViolation = "23505"
)

func (s *Storage) SaveURL(ctx context.Context, url string, alias string) error {
	const op = "storage.postgresql.SaveURL"

	err := s.db.WithContext(ctx).Raw(
		"INSERT INTO urls(url, alias) VALUES(?, ?)",
		url,
		alias,
	).Error

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			switch pgErr.ConstraintName {
			case "unique_url":
				return e.Err(op, storage.ErrURLExists)
			case "unique_alias":
				return e.Err(op, storage.ErrAliasExists)
			default:
				return e.Err(op, err)
			}
		}

		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return e.Err(op, storage.ErrTimeout)
		default:
			return e.Err(op, err)
		}
	}

	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	var resURL string

	err := s.db.WithContext(ctx).Raw(
		"SELECT url FROM urls WHERE alias= ?",
		alias,
	).Scan(&resURL).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return "", e.Err(op, storage.ErrURLNotFound)
		case errors.Is(err, context.DeadlineExceeded):
			return "", e.Err(op, storage.ErrTimeout)
		default:
			return "", e.Err(op, err)
		}
	}

	return resURL, nil
}

func (s *Storage) GetAlias(ctx context.Context, url string) (string, error) {
	const op = "storage.sqlite.GetAlias"

	var resAlias string

	err := s.db.WithContext(ctx).Raw(
		"SELECT alias FROM urls WHERE url= ?",
		url,
	).Scan(&resAlias).Error

	if resAlias == "" {
		return "", e.Err(op, storage.ErrAliasNotFound)
	}

	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return "", e.Err(op, storage.ErrTimeout)
		default:
			return "", e.Err(op, err)
		}
	}

	return resAlias, nil
}
