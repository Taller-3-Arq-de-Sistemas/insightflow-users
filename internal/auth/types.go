package auth

import "context"

type Service interface {
	Login(ctx context.Context, params LoginParams) (string, error)
	Register(ctx context.Context, params RegisterParams) (string, error)
	Logout(ctx context.Context) error
	ValidateToken(ctx context.Context, params ValidateTokenParams) (string, error)
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

type LoginParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterParams struct {
	Name      string `json:"name" validate:"required"`
	LastNames string `json:"last_names" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required"`
	BirthDate string `json:"birth_date" validate:"required"`
	Address   string `json:"address" validate:"required"`
	Phone     string `json:"phone" validate:"required,min=9,max=15"`
	Password  string `json:"password" validate:"required,min=8"`
}

type ValidateTokenParams struct {
	Token string `json:"token" validate:"required"`
}
