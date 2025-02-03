package repositories

import (
	"context"
	"mallbots/modules/user/domain/entities"
	"mallbots/modules/user/domain/interfaces"
	"mallbots/modules/user/infrastructure/query/gen"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	queries := gen.New(r.db)

	dbUser, err := queries.CreateUser(ctx, gen.CreateUserParams{
		Email:     user.Email,
		Password:  user.Password,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		FullName:  dbUser.FullName,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (*entities.User, error) {
	queries := gen.New(r.db)

	dbUser, err := queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		FullName:  dbUser.FullName,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	queries := gen.New(r.db)

	dbUser, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		FullName:  dbUser.FullName,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}, nil
}
