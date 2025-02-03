package services

import (
	"context"
	"mallbots/modules/user/application/dto"
	"mallbots/modules/user/domain/entities"
	"mallbots/modules/user/domain/interfaces"
	"mallbots/shared/errorx"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	VerifyMessage(message string, signedMessage string) (string, error)
}

type UserService struct {
	repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) interfaces.UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if user already exists
	if existingUser, _ := s.repo.GetByEmail(ctx, req.Email); existingUser != nil {
		return nil, errorx.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
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
		return nil, errorx.ErrCreateUser
	}

	return &dto.UserResponse{
		ID:        newUser.ID,
		Email:     newUser.Email,
		FullName:  newUser.FullName,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (string, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", errorx.ErrCannotGetUser
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errorx.ErrPasswordNotMatch
	}

	// TODO: Generate JWT token
	token := "dummy-token"

	return token, nil
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
