package url_shortener

import (
	"log/slog"
)

type UrlShortener struct {
	Log        *slog.Logger `json:"log,omitempty"`
	UrlStorage URLStorage   `json:"usr-saver,omitempty"`
}

func New(
	log *slog.Logger,
	urlStorage URLStorage,
) *UrlShortener {
	return &UrlShortener{
		Log:        log,
		UrlStorage: urlStorage,
	}
}
