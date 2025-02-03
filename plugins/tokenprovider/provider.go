package tokenprovider

import (
	"errors"
)

type Provider interface {
	Generate(data TokenPayload, expiry int) (Token, error)
	Validate(token string) (TokenPayload, error)
	SecretKey() string
}

type TokenPayload interface {
	GetUserId() int32
	GetSubToken() string
	GetEmail() string
}

type Token interface {
	GetToken() string
}

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrEncodingToken = errors.New("error encoding the token")
	ErrInvalidToken  = errors.New("invalid token provided")
)
