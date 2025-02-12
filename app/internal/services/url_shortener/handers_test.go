package url_shortener

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"grpc-service-ref/internal/storage"
	"log/slog"
)

type MockURLStorage struct {
	mock.Mock
}

func (m *MockURLStorage) SaveURL(ctx context.Context, url string, alias string) error {
	args := m.Called(ctx, url, alias)
	return args.Error(0)
}

func (m *MockURLStorage) GetURL(ctx context.Context, alias string) (string, error) {
	args := m.Called(ctx, alias)
	return args.String(0), args.Error(1)
}

func (m *MockURLStorage) GetAlias(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

func TestShorten_Success(t *testing.T) {
	mockStorage := new(MockURLStorage)
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	shortener := &UrlShortener{
		UrlStorage: mockStorage,
		Log:        log,
	}

	mockStorage.On("GetAlias", mock.Anything, "https://example.com").Return("", storage.ErrURLNotFound)
	mockStorage.On("SaveURL", mock.Anything, "https://example.com", mock.Anything).Return(nil)

	alias, err := shortener.Post(context.Background(), "https://example.com")

	assert.NoError(t, err)
	assert.NotEmpty(t, alias)
	mockStorage.AssertExpectations(t)
}

func TestShorten_AliasExists(t *testing.T) {
	mockStorage := new(MockURLStorage)
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	shortener := &UrlShortener{
		UrlStorage: mockStorage,
		Log:        log,
	}

	mockStorage.On("GetAlias", mock.Anything, "https://example.com").Return("", storage.ErrURLNotFound)
	mockStorage.On("SaveURL", mock.Anything, "https://example.com", mock.Anything).Return(storage.ErrAliasExists).Once()
	mockStorage.On("SaveURL", mock.Anything, "https://example.com", mock.Anything).Return(nil).Once()

	alias, err := shortener.Post(context.Background(), "https://example.com")

	assert.NoError(t, err)
	assert.NotEmpty(t, alias)
	mockStorage.AssertExpectations(t)
}

func TestShorten_StorageError(t *testing.T) {
	mockStorage := new(MockURLStorage)
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	shortener := &UrlShortener{
		UrlStorage: mockStorage,
		Log:        log,
	}

	mockStorage.On("GetAlias", mock.Anything, "https://example.com").Return("", storage.ErrURLNotFound)
	mockStorage.On("SaveURL", mock.Anything, "https://example.com", mock.Anything).Return(errors.New("storage error"))

	alias, err := shortener.Post(context.Background(), "https://example.com")

	assert.Error(t, err)
	assert.Empty(t, alias)
	mockStorage.AssertExpectations(t)
}

func TestGetOriginalUrl_Success(t *testing.T) {
	mockStorage := new(MockURLStorage)
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	shortener := &UrlShortener{
		UrlStorage: mockStorage,
		Log:        log,
	}

	mockStorage.On("GetURL", mock.Anything, "abc123").Return("https://example.com", nil)

	url, err := shortener.Post(context.Background(), "abc123")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", url)
	mockStorage.AssertExpectations(t)
}

func TestGetOriginalUrl_StorageError(t *testing.T) {
	mockStorage := new(MockURLStorage)
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	shortener := &UrlShortener{
		UrlStorage: mockStorage,
		Log:        log,
	}

	mockStorage.On("GetURL", mock.Anything, "abc123").Return("", errors.New("storage error"))

	url, err := shortener.Post(context.Background(), "abc123")

	assert.Error(t, err)
	assert.Empty(t, url)
	mockStorage.AssertExpectations(t)
}
