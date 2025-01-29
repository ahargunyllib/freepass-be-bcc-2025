package repository

import (
	"context"

	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/contracts"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type authRepository struct {
	db *sqlx.DB
}

func (a *authRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := a.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *authRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := a.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *authRepository) Register(ctx context.Context, user entity.User) error {
	_, err := a.db.NamedExecContext(ctx, `
		INSERT INTO users
		(id, name, email, password)
		VALUES (:id, :name, :email, :password)
	`, user)
	if err != nil {
		return err
	}

	return nil
}

func NewAuthRepository(db *sqlx.DB) contracts.AuthRepository {
	return &authRepository{
		db: db,
	}
}
