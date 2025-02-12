package sl

import (
	"fmt"
	"log/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func ErrStr(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
