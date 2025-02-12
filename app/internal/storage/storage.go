package storage

import "errors"

var (
	ErrNotSetDBParameter = errors.New("db parameter not set")
	ErrURLExists         = errors.New("URL already exists")
	ErrAliasExists       = errors.New("alias already exists")
	ErrURLNotFound       = errors.New("URL not found")
	ErrAliasNotFound     = errors.New("alias not found")
	ErrTimeout           = errors.New("timeout")
)
