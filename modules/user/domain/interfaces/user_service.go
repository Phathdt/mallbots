package interfaces

import (
	"context"
	"mallbots/modules/user/application/dto"
)

type UserService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userID int32) (*dto.UserResponse, error)
}
