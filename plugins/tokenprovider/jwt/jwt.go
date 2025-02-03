package jwt

import (
	"flag"
	"fmt"
	"mallbots/plugins/tokenprovider"
	"mallbots/shared/common"
	"time"

	"github.com/golang-jwt/jwt/v5"

	sctx "github.com/phathdt/service-context"
)

type jwtProvider struct {
	id     string
	secret string
}

func (j *jwtProvider) ID() string {
	return j.id
}

func New(id string) *jwtProvider {
	return &jwtProvider{id: id}
}

func (j *jwtProvider) InitFlags() {
	flag.StringVar(&j.secret, "jwt-secret", "secret-token", "Secret key for generating JWT")
}

func (j *jwtProvider) Activate(context sctx.ServiceContext) error {
	return nil
}

func (j *jwtProvider) Stop() error {
	return nil
}
func (j *jwtProvider) Configure() error {
	return nil
}

func (j *jwtProvider) Run() error {
	return nil
}

func (j *jwtProvider) SecretKey() string {
	return j.secret
}

func (j *jwtProvider) Generate(data tokenprovider.TokenPayload, expiry int) (tokenprovider.Token, error) {
	// generate the JWT
	now := time.Now()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		Payload: common.TokenPayload{
			UserId:   data.GetUserId(),
			Email:    data.GetEmail(),
			SubToken: data.GetSubToken()},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(now.Local().Add(time.Second*time.Duration(expiry)).Unix(), 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(now.Local().Unix(), 0)),
			ID:        fmt.Sprintf("%d", now.UnixNano()),
		},
	})

	myToken, err := t.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	// return the token
	return &token{
		Token:   myToken,
		Expiry:  expiry,
		Created: now,
	}, nil
}

func (j *jwtProvider) Validate(myToken string) (tokenprovider.TokenPayload, error) {
	res, err := jwt.ParseWithClaims(myToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, tokenprovider.ErrInvalidToken
	}

	// validate the token
	if !res.Valid {
		return nil, tokenprovider.ErrInvalidToken
	}

	claims, ok := res.Claims.(*myClaims)

	if !ok {
		return nil, tokenprovider.ErrInvalidToken
	}

	// return the token
	return claims.Payload, nil
}

type myClaims struct {
	Payload common.TokenPayload `json:"payload"`
	jwt.RegisteredClaims
}

type token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`
}

func (t *token) GetToken() string {
	return t.Token
}
