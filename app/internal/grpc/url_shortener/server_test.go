package url_shortener

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-service-ref/internal/storage"
	urlShortenerv1 "grpc-service-ref/protos/gen/go/url_shortener"
)

type MockUrlShortener struct {
	mock.Mock
}

func (m *MockUrlShortener) Post(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

func (m *MockUrlShortener) Get(ctx context.Context, alias string) (string, error) {
	args := m.Called(ctx, alias)
	return args.String(0), args.Error(1)
}

func TestShorten_Success(t *testing.T) {
	mockShortener := new(MockUrlShortener)
	server := &serverAPI{urlShortener: mockShortener}

	mockShortener.On("Shorten", mock.Anything, "https://example.com").Return("abc123", nil)

	resp, err := server.Post(context.Background(), &urlShortenerv1.PostRequest{
		Url: "https://example.com",
	})

	assert.NoError(t, err)
	assert.Equal(t, "abc123", resp.Alias)
	mockShortener.AssertExpectations(t)
}

func TestShorten_InvalidURL(t *testing.T) {
	mockShortener := new(MockUrlShortener)
	server := &serverAPI{urlShortener: mockShortener}

	resp, err := server.Post(context.Background(), &urlShortenerv1.PostRequest{
		Url: "invalid-url",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	mockShortener.AssertExpectations(t)
}

func TestShorten_StorageError(t *testing.T) {
	mockShortener := new(MockUrlShortener)
	server := &serverAPI{urlShortener: mockShortener}

	mockShortener.On("Shorten", mock.Anything, "https://example.com").Return("", storage.ErrTimeout)

	resp, err := server.Post(context.Background(), &urlShortenerv1.PostRequest{
		Url: "https://example.com",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.DeadlineExceeded, status.Code(err))
	mockShortener.AssertExpectations(t)
}

func TestGetOriginal_Success(t *testing.T) {
	mockShortener := new(MockUrlShortener)
	server := &serverAPI{urlShortener: mockShortener}

	mockShortener.On("GetOriginalUrl", mock.Anything, "abc123").Return("https://example.com", nil)

	resp, err := server.Get(context.Background(), &urlShortenerv1.GetRequest{
		Alias: "abc123",
	})

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", resp.Url)
	mockShortener.AssertExpectations(t)
}

func TestGetOriginal_InvalidAlias(t *testing.T) {
	mockShortener := new(MockUrlShortener)
	server := &serverAPI{urlShortener: mockShortener}

	mockShortener.On("GetOriginalUrl", mock.Anything, "invalid-alias").Return("", storage.ErrURLNotFound)

	resp, err := server.Get(context.Background(), &urlShortenerv1.GetRequest{
		Alias: "invalid-alias",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	mockShortener.AssertExpectations(t)
}

func TestGetOriginal_StorageError(t *testing.T) {
	mockShortener := new(MockUrlShortener)
	server := &serverAPI{urlShortener: mockShortener}

	mockShortener.On("GetOriginalUrl", mock.Anything, "abc123").Return("", errors.New("internal error"))

	resp, err := server.Get(context.Background(), &urlShortenerv1.GetRequest{
		Alias: "abc123",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Canceled, status.Code(err))
	mockShortener.AssertExpectations(t)
}
