package contracts

import (
	"context"

	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/dto"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/entity"
	"github.com/google/uuid"
)

type AuthRepository interface {
	Register(ctx context.Context, user entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	GetSession(ctx context.Context, query dto.GetSessionQuery) (dto.GetSessionResponse, error)
}
