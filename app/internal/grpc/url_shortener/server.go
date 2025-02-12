package url_shortener

import (
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	_ "google.golang.org/grpc/status"
	"grpc-service-ref/internal/storage"
	urlShortenerv1 "grpc-service-ref/protos/gen/go/url_shortener"
)

type UrlShortener interface {
	Shorten(
		ctx context.Context,
		url string,
	) (string, error)
	GetOriginalUrl(
		ctx context.Context,
		alias string,
	) (string, error)
}

type serverAPI struct {
	urlShortenerv1.UnimplementedUrlShortenerServer
	urlShortener UrlShortener
}

func Register(gRPCServer *grpc.Server, urlShortener UrlShortener) {
	urlShortenerv1.RegisterUrlShortenerServer(gRPCServer, &serverAPI{urlShortener: urlShortener})
}

func (s *serverAPI) Shorten(
	ctx context.Context,
	in *urlShortenerv1.ShortenRequest,
) (*urlShortenerv1.ShortenResponse, error) {
	if in.OriginalUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "missing original url")
	}
	if ok := govalidator.IsURL(in.OriginalUrl); !ok {
		return nil, status.Error(codes.InvalidArgument, "not valid url")
	}

	short, err := s.urlShortener.Shorten(ctx, in.OriginalUrl)
	if errors.Is(err, storage.ErrTimeout) {
		return nil, status.Error(codes.DeadlineExceeded, "failed to shorten url")
	} else if err != nil {
		return nil, status.Error(codes.Aborted, "failed to shorten url")
	}

	return &urlShortenerv1.ShortenResponse{ShortUrl: short}, nil
}

func (s *serverAPI) GetOriginal(
	ctx context.Context,
	in *urlShortenerv1.GetOriginalRequest,
) (*urlShortenerv1.GetOriginalResponse, error) {

	if in.ShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "missing alias")
	}

	url, err := s.urlShortener.GetOriginalUrl(ctx, in.ShortUrl)
	if errors.Is(err, storage.ErrTimeout) {
		return nil, status.Error(codes.DeadlineExceeded, "failed to get url")
	} else if errors.Is(err, storage.ErrURLNotFound) {
		return nil, status.Error(codes.InvalidArgument, "invalid alias")
	} else if err != nil {
		return nil, status.Error(codes.Canceled, "failed to get url")
	}

	return &urlShortenerv1.GetOriginalResponse{OriginalUrl: url}, nil
}
