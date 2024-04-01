package customErrors

import "errors"

var (
	ErrUnexpected                 = errors.New("unexpected error")
	ErrUnexpectedSigningMethod    = errors.New("unexpected signing method")
	ErrInvalidAccessToken         = errors.New("invalid access token")
	ErrAccessTokenExpired         = errors.New("access token expired")
	ErrMissingField               = errors.New("missing field")
	ErrNoActiveSymbols            = errors.New("no active symbols yet")
	ErrUnsuccessfulListenRequest  = errors.New("listen request unsuccessful")
	ErrMaximumNumberOfConnections = errors.New("maximum number of connections reached")
	ErrFailedRequest              = errors.New("request failed")
)
