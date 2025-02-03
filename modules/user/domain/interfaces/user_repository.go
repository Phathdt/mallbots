package interfaces

import (
	"context"
	"mallbots/modules/user/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
	GetByID(ctx context.Context, id int32) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
}
