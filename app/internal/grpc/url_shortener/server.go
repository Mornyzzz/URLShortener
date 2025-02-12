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
	Post(
		ctx context.Context,
		url string,
	) (string, error)
	Get(
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

func (s *serverAPI) Post(
	ctx context.Context,
	in *urlShortenerv1.PostRequest,
) (*urlShortenerv1.PostResponse, error) {
	if in.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "missing original url")
	}
	if ok := govalidator.IsURL(in.Url); !ok {
		return nil, status.Error(codes.InvalidArgument, "not valid url")
	}

	alias, err := s.urlShortener.Post(ctx, in.Url)
	if errors.Is(err, storage.ErrTimeout) {
		return nil, status.Error(codes.DeadlineExceeded, "failed to shorten url")
	} else if err != nil {
		return nil, status.Error(codes.Aborted, "failed to shorten url")
	}

	return &urlShortenerv1.PostResponse{Alias: alias}, nil
}

func (s *serverAPI) Get(
	ctx context.Context,
	in *urlShortenerv1.GetRequest,
) (*urlShortenerv1.GetResponse, error) {

	if in.Alias == "" {
		return nil, status.Error(codes.InvalidArgument, "missing alias")
	}

	url, err := s.urlShortener.Get(ctx, in.Alias)
	if errors.Is(err, storage.ErrTimeout) {
		return nil, status.Error(codes.DeadlineExceeded, "failed to get url")
	} else if errors.Is(err, storage.ErrURLNotFound) {
		return nil, status.Error(codes.InvalidArgument, "invalid alias")
	} else if err != nil {
		return nil, status.Error(codes.Canceled, "failed to get url")
	}

	return &urlShortenerv1.GetResponse{Url: url}, nil
}
