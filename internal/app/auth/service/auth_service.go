package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ahargunyllib/freepass-be-bcc-2025/domain"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/contracts"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/dto"
	"github.com/ahargunyllib/freepass-be-bcc-2025/domain/entity"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/bcrypt"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/jwt"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/uuid"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/validator"
)

type authService struct {
	repo      contracts.AuthRepository
	validator validator.ValidatorInterface
	uuid      uuid.CustomUUIDInterface
	bcrypt    bcrypt.CustomBcryptInterface
	jwt       jwt.CustomJwtInterface
}

func (a *authService) GetSession(ctx context.Context, query dto.GetSessionQuery) (dto.GetSessionResponse, error) {
	user, err := a.repo.FindByID(ctx, query.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GetSessionResponse{}, domain.ErrUserNotFound
		}

		return dto.GetSessionResponse{}, err
	}

	userResponse := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	if user.ImageURI.Valid {
		userResponse.ImageURI = &user.ImageURI.String
	}

	res := dto.GetSessionResponse{
		User: userResponse,
	}

	return res, nil
}

// Login implements contracts.AuthService.
func (a *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	valErr := a.validator.Validate(req)
	if valErr != nil {
		return dto.LoginResponse{}, valErr
	}

	user, err := a.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.LoginResponse{}, domain.ErrUserNotFound
		}

		return dto.LoginResponse{}, err
	}

	isValid := a.bcrypt.Compare(req.Password, user.Password)
	if !isValid {
		return dto.LoginResponse{}, domain.ErrInvalidCredentials
	}

	accessToken, err := a.jwt.Create(user.ID, user.Role)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	res := dto.LoginResponse{
		AccessToken: accessToken,
	}

	return res, nil
}

func (a *authService) Register(ctx context.Context, req dto.RegisterRequest) error {
	valErr := a.validator.Validate(req)
	if valErr != nil {
		return valErr
	}

	_, err := a.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return domain.ErrEmailAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	id, err := a.uuid.NewV7()
	if err != nil {
		return err
	}

	hashedPassword, err := a.bcrypt.Hash(req.Password)
	if err != nil {
		return err
	}

	user := entity.User{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	err = a.repo.Register(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func NewAuthService(
	repo contracts.AuthRepository,
	validator validator.ValidatorInterface,
	uuid uuid.CustomUUIDInterface,
	bcrypt bcrypt.CustomBcryptInterface,
	jwt jwt.CustomJwtInterface,
) contracts.AuthService {
	return &authService{
		repo:      repo,
		validator: validator,
		uuid:      uuid,
		bcrypt:    bcrypt,
		jwt:       jwt,
	}
}
