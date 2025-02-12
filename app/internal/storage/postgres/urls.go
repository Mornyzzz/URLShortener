package postgres

import (
	"context"
	_ "database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	_ "grpc-service-ref/internal/domain/models"
	"grpc-service-ref/internal/lib/logger/sl"
	"grpc-service-ref/internal/storage"
)

const (
	uniqueViolation = "23505"
)

func (s *Storage) SaveURL(ctx context.Context, url string, alias string) error {
	const op = "storage.postgres.SaveURL"
	err := s.db.WithContext(ctx).Exec(
		"INSERT INTO urls(url, alias) VALUES(?, ?) RETURNING id",
		url,
		alias,
	).Error

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			switch pgErr.ConstraintName {
			case "unique_url":
				return sl.ErrStr(op, storage.ErrURLExists)
			case "unique_alias":
				return sl.ErrStr(op, storage.ErrAliasExists)
			default:
				return sl.ErrStr(op, err)
			}
		}

		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return sl.ErrStr(op, storage.ErrTimeout)
		default:
			return sl.ErrStr(op, err)
		}
	}

	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	var resURL string

	err := s.db.WithContext(ctx).Raw(
		"SELECT url FROM urls WHERE alias= ?",
		alias,
	).Scan(&resURL).Error

	if resURL == "" {
		return "", sl.ErrStr(op, storage.ErrURLNotFound)
	}

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return "", sl.ErrStr(op, storage.ErrURLNotFound)
		case errors.Is(err, context.DeadlineExceeded):
			return "", sl.ErrStr(op, storage.ErrTimeout)
		default:
			return "", sl.ErrStr(op, err)
		}
	}

	return resURL, nil
}

func (s *Storage) GetAlias(ctx context.Context, url string) (string, error) {
	const op = "storage.postgres.GetAlias"

	var resAlias string

	err := s.db.WithContext(ctx).Raw(
		"SELECT alias FROM urls WHERE url= ?",
		url,
	).Scan(&resAlias).Error

	if resAlias == "" {
		return "", sl.ErrStr(op, storage.ErrAliasNotFound)
	}

	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return "", sl.ErrStr(op, storage.ErrTimeout)
		default:
			return "", sl.ErrStr(op, err)
		}
	}

	return resAlias, nil
}
