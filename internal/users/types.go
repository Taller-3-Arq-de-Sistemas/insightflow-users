package users

import (
	"context"

	repository "github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/adapters/postgres/sqlc"
)

type Service interface {
	CreateUser(ctx context.Context, params CreateUserParams) (repository.FindUserByIdRow, error)
	ListUsers(ctx context.Context) ([]repository.ListUsersRow, error)
	FindUserById(ctx context.Context, id string) (repository.FindUserByIdRow, error)
	UpdateUser(ctx context.Context, params repository.UpdateUserParams) (repository.FindUserByIdRow, error)
	DeleteUser(ctx context.Context, id string) error
}

type CreateUserParams struct {
	Name      string `json:"name" validate:"required"`
	LastNames string `json:"last_names" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
	BirthDate string `json:"birth_date" validate:"required"`
	Address   string `json:"address" validate:"required"`
	Phone     string `json:"phone" validate:"required,min=9,max=15"`
	Status    string `json:"status"`
	Role      string `json:"role"`
}

type UpdateUserParams struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
}
