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
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/uuid"
	"github.com/ahargunyllib/freepass-be-bcc-2025/pkg/validator"
)

type userService struct {
	repo      contracts.UserRepository
	validator validator.ValidatorInterface
	uuid      uuid.CustomUUIDInterface
	bcrypt    bcrypt.CustomBcryptInterface
}

func (u *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) error {
	valErr := u.validator.Validate(req)
	if valErr != nil {
		return valErr
	}

	_, err := u.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return domain.ErrEmailAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	hashedPassword, err := u.bcrypt.Hash(req.Password)
	if err != nil {
		return err
	}

	id, err := u.uuid.NewV7()
	if err != nil {
		return err
	}

	user := &entity.User{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     2, // event coordinator
	}

	err = u.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) UpdateUser(ctx context.Context, req dto.UpdateUserRequest) error {
	valErr := u.validator.Validate(req)
	if valErr != nil {
		return valErr
	}

	user, err := u.repo.FindByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}

		return err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Password != "" {
		hashedPassword, hashErr := u.bcrypt.Hash(req.Password)
		if hashErr != nil {
			return hashErr
		}

		user.Password = hashedPassword
	}

	err = u.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) DeleteUser(ctx context.Context, query dto.DeleteUserQuery) error {
	valErr := u.validator.Validate(query)
	if valErr != nil {
		return valErr
	}

	user, err := u.repo.FindByID(ctx, query.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}

		return err
	}

	if user.Role == 3 {
		return domain.ErrCannotDeleteAdmin
	}

	err = u.repo.Delete(ctx, query.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) GetUser(ctx context.Context, query dto.GetUserQuery) (dto.GetUserResponse, error) {
	valErr := u.validator.Validate(query)
	if valErr != nil {
		return dto.GetUserResponse{}, valErr
	}

	user, err := u.repo.FindByID(ctx, query.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GetUserResponse{}, domain.ErrUserNotFound
		}

		return dto.GetUserResponse{}, err
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

	res := dto.GetUserResponse{
		User: userResponse,
	}

	return res, nil
}

func (u *userService) GetUsers(ctx context.Context, query dto.GetUsersQuery) (dto.GetUsersResponse, error) {
	valErr := u.validator.Validate(query)
	if valErr != nil {
		return dto.GetUsersResponse{}, valErr
	}

	if query.Page < 1 {
		query.Page = 1
	}

	if query.Limit < 1 {
		query.Limit = 10
	}

	if query.SortBy == "" {
		query.SortBy = "id"
	}

	if query.SortOrder == "" {
		query.SortOrder = "asc"
	}

	users, err := u.repo.FindAll(
		ctx,
		query.Limit,
		(query.Page-1)*query.Limit,
		query.SortBy,
		query.SortOrder,
		query.Search,
		query.Role,
	)
	if err != nil {
		return dto.GetUsersResponse{}, err
	}

	totalData, err := u.repo.Count(ctx, query.Search, query.Role)
	if err != nil {
		return dto.GetUsersResponse{}, err
	}

	totalPage := int(totalData) / query.Limit
	if int(totalData)%query.Limit != 0 {
		totalPage++
	}

	meta := dto.PaginationResponse{
		TotalData: totalData,
		TotalPage: totalPage,
		Page:      query.Page,
		Limit:     query.Limit,
	}

	res := dto.GetUsersResponse{
		Users: make([]dto.UserResponse, 0, len(users)),
		Meta:  meta,
	}

	for _, user := range users {
		userResponse := dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		}

		if user.ImageURI.Valid {
			userResponse.ImageURI = &user.ImageURI.String
		}

		res.Users = append(res.Users, userResponse)
	}

	return res, nil
}

func NewUserService(
	repo contracts.UserRepository,
	validator validator.ValidatorInterface,
	uuid uuid.CustomUUIDInterface,
	bcrypt bcrypt.CustomBcryptInterface,
) contracts.UserService {
	return &userService{
		repo:      repo,
		validator: validator,
		uuid:      uuid,
		bcrypt:    bcrypt,
	}
}
