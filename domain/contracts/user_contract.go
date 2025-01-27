package contracts

import (
	"context"

	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/dto"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindAll(ctx context.Context, limit, offset int, sortBy, sortOrder, search string, role int16) ([]entity.User, error)
	Count(ctx context.Context, search string, role int16) (int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserService interface {
	GetUsers(ctx context.Context, query dto.GetUsersQuery) (dto.GetUsersResponse, error)
	GetUser(ctx context.Context, query dto.GetUserQuery) (dto.GetUserResponse, error)
	CreateUser(ctx context.Context, req dto.CreateUserRequest) error
	DeleteUser(ctx context.Context, query dto.DeleteUserQuery) error
}
