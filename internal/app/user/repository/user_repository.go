package repository

import (
	"context"
	"fmt"

	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/contracts"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func (u *userRepository) Create(ctx context.Context, user entity.User) error {
	_, err := u.db.NamedExecContext(ctx, `
		INSERT INTO users
		(id, name, email, password, role)
		VALUES (:id, :name, :email, :password, :role)
		`, user,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := u.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) FindAll(
	ctx context.Context,
	limit int,
	offset int,
	sortBy string,
	sortOrder string,
	search string,
	role int16,
) ([]entity.User, error) {
	users := []entity.User{}
	query := "SELECT * FROM users WHERE 1=1"
	args := []interface{}{}

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d)", len(args)+1, len(args)+1)
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	if role != 0 {
		query += fmt.Sprintf(" AND role = $%d", len(args)+1)
		args = append(args, role)
	}

	query += fmt.Sprintf(" ORDER BY $%d $%d LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2, len(args)+3, len(args)+4)
	args = append(args, sortBy, sortOrder, limit, offset)

	err := u.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userRepository) Count(ctx context.Context, search string, role int16) (int64, error) {
	var count int64
	query := "SELECT COUNT(*) FROM users WHERE 1=1"
	args := []interface{}{}

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d)", len(args)+1, len(args)+1)
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	if role != 0 {
		query += fmt.Sprintf(" AND role = $%d", len(args)+1)
		args = append(args, role)
	}

	err := u.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func NewUserRepository(db *sqlx.DB) contracts.UserRepository {
	return &userRepository{
		db: db,
	}
}
