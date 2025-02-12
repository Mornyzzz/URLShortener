package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"grpc-service-ref/internal/storage"
	e "grpc-service-ref/pkg"
)

const (
	aliasPrefix = "alias:"
	URLPrefix   = "url:"
)

func (s *Storage) buildKey(prefix string, key string) string {
	return prefix + key
}

func (s *Storage) SaveURL(ctx context.Context, url string, alias string) error {
	const op = "storage.redis.SaveURL"

	aliasKey := s.buildKey(aliasPrefix, url)
	urlKey := s.buildKey(URLPrefix, alias)

	err := s.r.Watch(ctx, func(tx *redis.Tx) error {
		existsAlias, err := tx.Exists(ctx, aliasKey).Result()
		if err != nil {
			return err
		}
		if existsAlias == 1 {
			return storage.ErrAliasExists
		}

		existsUrl, err := tx.Exists(ctx, urlKey).Result()
		if err != nil {
			return err
		}
		if existsUrl == 1 {
			return storage.ErrURLExists
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, aliasKey, alias, 0)
			pipe.Set(ctx, urlKey, url, 0)
			return nil
		})
		return err
	}, aliasKey, urlKey)

	if err != nil {
		return e.Err(op, err)
	}

	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.redis.GetURL"

	val, err := s.r.Get(ctx, s.buildKey(URLPrefix, alias)).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return "", storage.ErrURLNotFound
	case err != nil:
		return "", e.Err(op, err)
	}

	return val, nil
}

func (s *Storage) GetAlias(ctx context.Context, url string) (string, error) {
	const op = "storage.redis.GetAlias"

	val, err := s.r.Get(ctx, s.buildKey(aliasPrefix, url)).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return "", storage.ErrAliasNotFound
	case err != nil:
		return "", e.Err(op, err)
	}

	return val, nil
}
