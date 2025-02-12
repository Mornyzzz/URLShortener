package url_shortener

import (
	"errors"
	"log/slog"
)

type UrlShortener struct {
	log        *slog.Logger `json:"log,omitempty"`
	urlStorage URLStorage   `json:"usr-saver,omitempty"`
}

var (
	ErrInvalidURL = errors.New("invalid URL")
)

func New(
	log *slog.Logger,
	urlStorage URLStorage,
) *UrlShortener {
	return &UrlShortener{
		log:        log,
		urlStorage: urlStorage,
	}
}
