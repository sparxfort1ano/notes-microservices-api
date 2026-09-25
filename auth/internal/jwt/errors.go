package jwt

import "errors"

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token has expired")
	ErrInvalidTokenType  = errors.New("incorrect token type")
	ErrTokenGeneration   = errors.New("token generation error")
	ErrMissingUserID     = errors.New("token user ID missing")
	ErrInvalidSignature  = errors.New("incorrect token signature")
	ErrMissingAuthHeader = errors.New("header Authorization missing")
	ErrInvalidAuthFormat = errors.New("incorrect token format")
)

const (
	MsgInvalidToken  = "invalid token"
	MsgTokenRequired = "token missing or incorrect format"
)
