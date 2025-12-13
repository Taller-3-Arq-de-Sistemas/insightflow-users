package users

import (
	"context"
	"database/sql"
	"strings"
	"time"

	repository "github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/adapters/sqlite/sqlc"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type svc struct {
	repo repository.Querier
}

func NewService(repo repository.Querier) *svc {
	return &svc{repo: repo}
}

func (s *svc) CreateUser(ctx context.Context, params CreateUserParams) (repository.FindUserByIdRow, error) {
	id := uuid.New().String()
	birthDate, err := time.Parse("2006-01-02", params.BirthDate)
	if err != nil {
		return repository.FindUserByIdRow{}, err
	}
	if params.Status == "" {
		params.Status = "active"
	}
	if params.Role == "" {
		params.Role = "user"
	}
	if params.Status != "active" && params.Status != "inactive" {
		return repository.FindUserByIdRow{}, ErrInvalidStatus
	}
	if params.Role != "user" && params.Role != "admin" && params.Role != "editor" {
		return repository.FindUserByIdRow{}, ErrInvalidRole
	}
	_, err = s.repo.FindUserByEmail(ctx, params.Email)
	if err == nil {
		return repository.FindUserByIdRow{}, ErrUserAlreadyExists
	}
	_, err = s.repo.FindUserByUsername(ctx, params.Username)
	if err == nil {
		return repository.FindUserByIdRow{}, ErrUserAlreadyExists
	}

	role, err := s.repo.FindRoleByName(ctx, params.Role)
	if err != nil {
		return repository.FindUserByIdRow{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return repository.FindUserByIdRow{}, err
	}

	var newID string
	newID, err = s.repo.CreateUser(ctx, repository.CreateUserParams{
		ID:        id,
		Name:      params.Name,
		LastNames: params.LastNames,
		Email:     params.Email,
		Username:  params.Username,
		Password:  string(hashedPassword),
		BirthDate: birthDate,
		Address:   params.Address,
		Phone:     params.Phone,
		Status:    params.Status,
		RoleID:    role.ID,
	})
	if err != nil {
		return repository.FindUserByIdRow{}, err
	}

	return s.repo.FindUserById(ctx, newID)
}

func (s *svc) ListUsers(ctx context.Context) ([]repository.ListUsersRow, error) {
	return s.repo.ListUsers(ctx)
}

func (s *svc) FindUserById(ctx context.Context, id string) (repository.FindUserByIdRow, error) {
	user, err := s.repo.FindUserById(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		return repository.FindUserByIdRow{}, err
	}
	if err == sql.ErrNoRows {
		return repository.FindUserByIdRow{}, ErrUserNotFound
	}
	return user, nil
}

func (s *svc) UpdateUser(ctx context.Context, id string, params UpdateUserParams) (repository.FindUserByIdRow, error) {
	currentUser, err := s.repo.FindUserById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.FindUserByIdRow{}, ErrUserNotFound
		}
		return repository.FindUserByIdRow{}, err
	}

	name := currentUser.Name
	lastNames := currentUser.LastNames
	username := currentUser.Username

	if params.Username != "" && params.Username != currentUser.Username {
		existing, err := s.repo.FindUserByUsername(ctx, params.Username)
		if err == nil {
			_ = existing
			return repository.FindUserByIdRow{}, ErrUserAlreadyExists
		}
		if err != sql.ErrNoRows {
			return repository.FindUserByIdRow{}, err
		}
		username = params.Username
	}

	if params.FullName != "" {
		parts := strings.Fields(params.FullName)
		if len(parts) < 3 {
			return repository.FindUserByIdRow{}, ErrInvalidFullname
		}
		name = parts[0]
		lastNames = strings.Join(parts[1:], " ")
	}

	id, err = s.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:        id,
		Name:      name,
		LastNames: lastNames,
		Username:  username,
	})
	if err != nil {
		return repository.FindUserByIdRow{}, err
	}

	return s.repo.FindUserById(ctx, id)
}

func (s *svc) DeleteUser(ctx context.Context, id string) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return ErrUserNotFound
	}
	return nil
}
