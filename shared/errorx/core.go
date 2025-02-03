package errorx

import "errors"

var (
	ErrCannotGetUser     = errors.New("cannot get user")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCreateUser        = errors.New("create user failed")
	ErrPasswordNotMatch  = errors.New("password not match")
	ErrGenToken          = errors.New("when gen token")
)
