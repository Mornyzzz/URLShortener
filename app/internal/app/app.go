package app

import (
	grpcapp "grpc-service-ref/internal/app/grpc"
	"grpc-service-ref/internal/config"
	"grpc-service-ref/internal/services/url_shortener"
	p "grpc-service-ref/internal/storage/postgres"
	r "grpc-service-ref/internal/storage/redis"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, cfg config.Config) *App {
	var storage url_shortener.URLStorage
	var err error

	log.Info("start app")
	switch cfg.StorageType {

	case "postgres":
		log.Info("init postgres...")
		storage, err = p.New(cfg.Postgres)
		if err != nil {
			panic("no postgres storage")
		}

	case "redis":
		log.Info("init redis...")
		storage, err = r.New(cfg.Redis)
		if err != nil {
			panic("no redis storage")
		}

	default:
		panic("no/unknown storage")
	}

	log.Info("storage successfully init")

	urlShortener := url_shortener.New(log, storage)

	grpcApp := grpcapp.New(log, urlShortener, cfg.GRPC.Port)

	return &App{
		GRPCServer: grpcApp,
	}
}
