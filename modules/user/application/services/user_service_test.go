package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"mallbots/modules/user/application/dto"
	"mallbots/modules/user/domain/entities"
	"mallbots/plugins/tokenprovider"
	"mallbots/shared/errorx"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// Mock repository
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id int32) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

// Mock token provider
type MockTokenProvider struct {
	mock.Mock
}

func (m *MockTokenProvider) Generate(data tokenprovider.TokenPayload, expiry int) (tokenprovider.Token, error) {
	args := m.Called(data, expiry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(tokenprovider.Token), args.Error(1)
}

func (m *MockTokenProvider) Validate(token string) (tokenprovider.TokenPayload, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(tokenprovider.TokenPayload), args.Error(1)
}

func (m *MockTokenProvider) SecretKey() string {
	args := m.Called()
	return args.String(0)
}

// Mock token
type MockToken struct {
	tokenString string
}

func NewMockToken(token string) *MockToken {
	return &MockToken{tokenString: token}
}

func (m *MockToken) GetToken() string {
	return m.tokenString
}

func TestUserService_Register(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenProvider := new(MockTokenProvider)
	service := NewUserService(mockRepo, mockTokenProvider)

	testCases := []struct {
		name    string
		req     *dto.RegisterRequest
		setup   func()
		wantErr error
	}{
		{
			name: "Successful registration",
			req: &dto.RegisterRequest{
				Email:    "new@example.com",
				Password: "password123",
				FullName: "New User",
			},
			setup: func() {
				// Expect check for existing user
				mockRepo.On("GetByEmail", mock.Anything, "new@example.com").
					Return(nil, errors.New("user not found")).Once()

				// Expect user creation
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(user *entities.User) bool {
					return user.Email == "new@example.com" && user.FullName == "New User"
				})).Return(&entities.User{
					ID:        1,
					Email:     "new@example.com",
					FullName:  "New User",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil).Once()

				// Expect token generation
				mockTokenProvider.On("Generate", mock.MatchedBy(func(payload tokenprovider.TokenPayload) bool {
					return payload.GetEmail() == "new@example.com" && payload.GetUserId() == int32(1)
				}), 3600*24*30).Return(NewMockToken("valid.token.here"), nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "User already exists",
			req: &dto.RegisterRequest{
				Email:    "existing@example.com",
				Password: "password123",
				FullName: "Existing User",
			},
			setup: func() {
				mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").
					Return(&entities.User{
						ID:       1,
						Email:    "existing@example.com",
						FullName: "Existing User",
					}, nil).Once()
			},
			wantErr: errorx.ErrUserAlreadyExists,
		},
		{
			name: "Database error during creation",
			req: &dto.RegisterRequest{
				Email:    "error@example.com",
				Password: "password123",
				FullName: "Error User",
			},
			setup: func() {
				mockRepo.On("GetByEmail", mock.Anything, "error@example.com").
					Return(nil, errors.New("user not found")).Once()

				mockRepo.On("Create", mock.Anything, mock.Anything).
					Return(nil, errorx.ErrCreateUser).Once()
			},
			wantErr: errorx.ErrCreateUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.setup()

			// Execute
			token, err := service.Register(context.Background(), tc.req)

			// Assert
			if tc.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.wantErr, err)
				require.Empty(t, token)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, token)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
			mockTokenProvider.AssertExpectations(t)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenProvider := new(MockTokenProvider)
	service := NewUserService(mockRepo, mockTokenProvider)

	// Create real bcrypt hash for "password123"
	correctHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		req     *dto.LoginRequest
		setup   func(string)
		wantErr error
	}{
		{
			name: "Successful login",
			req: &dto.LoginRequest{
				Email:    "user@example.com",
				Password: "password123",
			},
			setup: func(hash string) {
				mockRepo.On("GetByEmail", mock.Anything, "user@example.com").
					Return(&entities.User{
						ID:       1,
						Email:    "user@example.com",
						Password: hash,
					}, nil).Once()

				mockToken := NewMockToken("valid.token.here")
				mockTokenProvider.On("Generate", mock.MatchedBy(func(payload tokenprovider.TokenPayload) bool {
					return payload.GetEmail() == "user@example.com" && payload.GetUserId() == int32(1)
				}), 3600*24*30).Return(mockToken, nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "User not found",
			req: &dto.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setup: func(hash string) {
				mockRepo.On("GetByEmail", mock.Anything, "nonexistent@example.com").
					Return(nil, errorx.ErrCannotGetUser).Once()
			},
			wantErr: errorx.ErrCannotGetUser,
		},
		{
			name: "Wrong password",
			req: &dto.LoginRequest{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			setup: func(hash string) {
				mockRepo.On("GetByEmail", mock.Anything, "user@example.com").
					Return(&entities.User{
						ID:       1,
						Email:    "user@example.com",
						Password: hash,
					}, nil).Once()
			},
			wantErr: errorx.ErrPasswordNotMatch,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mocks before each test case
			mockRepo.ExpectedCalls = nil
			mockTokenProvider.ExpectedCalls = nil

			// Setup mock expectations with the real hash
			tc.setup(string(correctHash))

			// Execute
			token, err := service.Login(context.Background(), tc.req)

			// Assert
			if tc.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.wantErr, err)
				require.Empty(t, token)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, token)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
			mockTokenProvider.AssertExpectations(t)
		})
	}
}

func TestUserService_GetProfile(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenProvider := new(MockTokenProvider)
	service := NewUserService(mockRepo, mockTokenProvider)

	testCases := []struct {
		name    string
		userID  int32
		setup   func()
		wantErr error
	}{
		{
			name:   "Get existing profile",
			userID: 1,
			setup: func() {
				mockRepo.On("GetByID", mock.Anything, int32(1)).
					Return(&entities.User{
						ID:        1,
						Email:     "user@example.com",
						FullName:  "Test User",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil).Once()
			},
			wantErr: nil,
		},
		{
			name:   "Profile not found",
			userID: 999,
			setup: func() {
				mockRepo.On("GetByID", mock.Anything, int32(999)).
					Return(nil, errorx.ErrCannotGetUser).Once()
			},
			wantErr: errorx.ErrCannotGetUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			tc.setup()

			// Execute
			profile, err := service.GetProfile(context.Background(), tc.userID)

			// Assert
			if tc.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.wantErr, err)
				require.Nil(t, profile)
			} else {
				require.NoError(t, err)
				require.NotNil(t, profile)
				require.Equal(t, tc.userID, profile.ID)
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
