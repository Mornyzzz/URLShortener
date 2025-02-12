package url_shortener

import (
	"context"
	"errors"
	"grpc-service-ref/internal/lib/logger/sl"
	"grpc-service-ref/internal/storage"
	e "grpc-service-ref/pkg"
	"log/slog"
	"math/rand"
)

type URLStorage interface {
	SaveURL(
		ctx context.Context,
		url string,
		alias string,
	) (err error)
	GetURL(
		ctx context.Context,
		alias string,
	) (url string, err error)
	GetAlias(
		ctx context.Context,
		url string,
	) (alias string, err error)
}

const (
	allowedChars = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789" +
		"_"
	shortURLLen = 10
	//allowedChars = "ab"
	//shortURLLen  = 2
)

func (u *UrlShortener) Shorten(ctx context.Context, url string) (string, error) {
	const op = "services.url_shortener.Shorten"

	log := u.log.With(
		slog.String("op", op),
		slog.String("url", url),
	)

	log.Info("checking alias in db...")
	alias, err := u.urlStorage.GetAlias(ctx, url)
	if err == nil {
		log.Info("find alias:", alias, err)
		return alias, nil
	}

	log.Info("no alias in db")
	log.Info("attempting to generate url...")

	for {
		log.Info("generate alias")
		result := make([]byte, shortURLLen)
		for i := range result {
			result[i] = allowedChars[rand.Intn(len(allowedChars))]
		}
		alias = string(result)

		log.Info("generated alias: ", alias)

		err = u.urlStorage.SaveURL(ctx, url, alias)

		if errors.Is(err, storage.ErrAliasExists) {
			log.Warn("alias exists", sl.Err(err))
		} else if err != nil {
			log.Error("failed to save alias", sl.Err(err))
			return "", e.Err(op, err)
		} else {
			u.log.Info("successfully save alias", alias)
			return alias, nil
		}
	}
}

func (u *UrlShortener) GetOriginalUrl(ctx context.Context, alias string) (string, error) {
	const op = "services.url_shortener.GetOriginalUrl"

	log := u.log.With(
		slog.String("op", op),
		slog.String("alias", alias),
	)

	log.Info("attempting to get url...")

	url, err := u.urlStorage.GetURL(ctx, alias)
	if errors.Is(err, storage.ErrURLNotFound) {
		log.Error("url not found", sl.Err(err))
		return "", e.Err(op, err)
	} else if errors.Is(err, storage.ErrTimeout) {
		log.Error("timeout", sl.Err(err))
		return "", e.Err(op, err)
	} else if err != nil {
		log.Error("error getting url", sl.Err(err))
		return "", e.Err(op, err)
	}
	log.Info("successfully got url")
	return url, nil
}
