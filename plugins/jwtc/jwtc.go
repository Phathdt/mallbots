package jwtc

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	sctx "github.com/phathdt/service-context"
)

type JWTComp interface {
	Generate(data TokenPayload, expiry int) (Token, error)
	Validate(token string) (TokenPayload, error)
	SecretKey() string
}

type TokenPayload interface {
}

type Token interface {
	GetToken() string
}

type jwtComp struct {
	id     string
	secret string
	logger sctx.Logger
}

func New(id string) *jwtComp {
	return &jwtComp{id: id}
}

func (j *jwtComp) ID() string {
	return j.id
}

func (j *jwtComp) InitFlags() {
	flag.StringVar(
		&j.secret,
		"jwt-secret",
		"secret-token",
		"Secret key for generating JWT",
	)
}

func (j *jwtComp) Activate(sc sctx.ServiceContext) error {
	j.logger = sctx.GlobalLogger().GetLogger(j.id)
	return nil
}

func (j *jwtComp) Stop() error {
	return nil
}

func (j *jwtComp) SecretKey() string {
	return j.secret
}

type claims struct {
	Payload TokenPayload `json:"payload"`
	jwt.RegisteredClaims
}

type token struct {
	token   string
	created time.Time
	expiry  int
}

func (t *token) GetToken() string {
	return t.token
}

func (j *jwtComp) Generate(data TokenPayload, expiry int) (Token, error) {
	now := time.Now()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		Payload: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(expiry))),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%d", now.UnixNano()),
		},
	})

	signedToken, err := t.SignedString([]byte(j.secret))
	if err != nil {
		j.logger.Error("Error signing token", err.Error())
		return nil, fmt.Errorf("error signing token: %w", err)
	}

	return &token{
		token:   signedToken,
		created: now,
		expiry:  expiry,
	}, nil
}

func (j *jwtComp) Validate(tokenString string) (TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		j.logger.Error("Error parsing token", err.Error())
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims.Payload, nil
}
