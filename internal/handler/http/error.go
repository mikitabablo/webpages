package http

import "errors"

var (
	ErrCannotParseRequest = errors.New("cannot parse request")
	ErrInvalidURL         = errors.New("invalid url")
)
