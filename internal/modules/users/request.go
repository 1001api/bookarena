package users

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUserRequest struct {
	Username     string `json:"username" form:"username" validate:"required"`
	Email        string `json:"email" form:"email" validate:"required"`
	PasswordHash string `json:"-"`
	Role         string `json:"role" form:"role" validate:"required"`
}

type UpdateUserRequest struct {
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
}

type UserProfileResponse struct {
	ID          uuid.UUID          `json:"id"`
	Email       string             `json:"email"`
	Username    string             `json:"username"`
	Role        string             `json:"role"`
	LastLoginAt pgtype.Timestamptz `json:"last_login_at"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}
