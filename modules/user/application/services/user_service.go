package services

import (
	"context"
	"mallbots/modules/user/application/dto"
	"mallbots/modules/user/domain/entities"
	"mallbots/modules/user/domain/interfaces"
	"mallbots/plugins/tokenprovider"
	"mallbots/shared/common"
	"mallbots/shared/errorx"
	"time"

	"github.com/jaevor/go-nanoid"
	"github.com/phathdt/service-context/core"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo          interfaces.UserRepository
	tokenProvider tokenprovider.Provider
}

func NewUserService(repo interfaces.UserRepository, tokenProvider tokenprovider.Provider) interfaces.UserService {
	return &UserService{
		repo:          repo,
		tokenProvider: tokenProvider,
	}
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (string, error) {
	// Check if user already exists
	if existingUser, _ := s.repo.GetByEmail(ctx, req.Email); existingUser != nil {
		return "", errorx.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &entities.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return "", errorx.ErrCreateUser
	}

	canonicID, _ := nanoid.Standard(21)
	subToken := canonicID()

	payload := common.TokenPayload{
		UserId:   newUser.ID,
		Email:    newUser.Email,
		SubToken: subToken,
	}

	expiredTime := 3600 * 24 * 30
	accessToken, err := s.tokenProvider.Generate(&payload, expiredTime)
	if err != nil {
		return "", core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	return accessToken.GetToken(), nil
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (string, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", errorx.ErrCannotGetUser
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errorx.ErrPasswordNotMatch
	}

	canonicID, _ := nanoid.Standard(21)
	subToken := canonicID()

	payload := common.TokenPayload{
		UserId:   user.ID,
		Email:    user.Email,
		SubToken: subToken,
	}

	expiredTime := 3600 * 24 * 30
	accessToken, err := s.tokenProvider.Generate(&payload, expiredTime)
	if err != nil {
		return "", core.ErrBadRequest.
			WithError(errorx.ErrGenToken.Error()).
			WithDebug(err.Error())
	}

	return accessToken.GetToken(), nil
}

func (s *UserService) GetProfile(ctx context.Context, userID int32) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, errorx.ErrCannotGetUser
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
