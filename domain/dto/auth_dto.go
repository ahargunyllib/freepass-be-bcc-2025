package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type GetSessionQuery struct {
	UserID uuid.UUID // assign from context
}

type GetSessionResponse struct {
	User UserResponse `json:"user"`
}
