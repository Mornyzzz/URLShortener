package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"grpc-service-ref/internal/config"
	"grpc-service-ref/internal/storage"
	e "grpc-service-ref/pkg"
)

type Storage struct {
	r *redis.Client
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewRedisConfig(cfg config.RedisConfig) (*RedisConfig, error) {
	const op = "storage.redis.NewRedisConfig"

	if cfg.Host == "" || cfg.Port == 0 {
		return nil, e.Err(op, storage.ErrNotSetDBParameter)
	}

	return &RedisConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	}, nil
}

func New(cfg config.RedisConfig) (*Storage, error) {
	const op = "storage.redis.New"

	redisCfg, err := NewRedisConfig(cfg)
	if err != nil {
		return nil, e.Err(op, err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	return &Storage{
		r: client,
	}, nil
}

func (s *Storage) Ping(ctx context.Context) error {
	const op = "storage.redis.Ping"
	err := s.r.Ping(ctx).Err()
	if err != nil {
		return e.Err(op, err)
	}
	return nil
}

func (s *Storage) Stop() error {
	const op = "storage.redis.Close"
	err := s.r.Close()
	if err != nil {
		return e.Err(op, err)
	}
	return nil
}
