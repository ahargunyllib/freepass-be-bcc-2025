package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     int16     `json:"role"`
	ImageURI *string   `json:"image_uri"`
}

type GetUsersQuery struct {
	Role      int16  `query:"role" validate:"omitempty,numeric,min=0,max=2"`
	Search    string `query:"search"`
	Limit     int    `query:"limit" validate:"omitempty,numeric,min=1,max=100"`
	Page      int    `query:"page" validate:"omitempty,numeric,min=1"`
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
	SortBy    string `query:"sort_by" validate:"omitempty,oneof=id name email role created_at updated_at"`
}

type GetUsersResponse struct {
	Users []UserResponse     `json:"users"`
	Meta  PaginationResponse `json:"meta"`
}

type GetUserQuery struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

type GetUserResponse struct {
	User UserResponse `json:"user"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID // from context
	Name     string    `json:"name" validate:"omitempty,min=3,max=255"`
	Password string    `json:"password" validate:"omitempty,min=8,max=255"`
}

type DeleteUserQuery struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}
