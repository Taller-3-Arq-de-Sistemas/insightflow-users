package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/config"
	repository "github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/adapters/postgres/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type svc struct {
	repo repository.Querier
}

func NewService(repo repository.Querier) *svc {
	return &svc{repo: repo}
}

func (s *svc) GenerateToken(id string, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  id,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetString("JWT_SECRET", "secret")))
}

func (s *svc) Login(ctx context.Context, params LoginParams) (string, error) {
	if params.Email == "" || params.Password == "" {
		return "", ErrInvalidCredentials
	}

	user, err := s.repo.FindUserByEmail(ctx, params.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	return s.GenerateToken(user.ID, user.Role)
}

func (s *svc) Register(ctx context.Context, params RegisterParams) (string, error) {
	if params.Name == "" || params.LastNames == "" || params.Email == "" || params.Username == "" || params.Password == "" || params.BirthDate == "" || params.Address == "" || params.Phone == "" {
		return "", ErrInvalidCredentials
	}

	_, err := s.repo.FindUserByEmail(ctx, params.Email)
	if err == nil {
		return "", ErrUserAlreadyExists
	}

	_, err = s.repo.FindUserByUsername(ctx, params.Username)
	if err == nil {
		return "", ErrUserAlreadyExists
	}

	if err != sql.ErrNoRows {
		return "", err
	}

	birthDate, err := time.Parse("2006-01-02", params.BirthDate)
	if err != nil {
		return "", err
	}

	if time.Now().Before(birthDate) {
		return "", ErrInvalidBirthDate
	}

	role, err := s.repo.FindRoleByName(ctx, "user")
	if err != nil {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	id, err := s.repo.CreateUser(ctx, repository.CreateUserParams{
		ID:        uuid.New().String(),
		Name:      params.Name,
		LastNames: params.LastNames,
		Email:     params.Email,
		Username:  params.Username,
		Password:  string(hashedPassword),
		BirthDate: birthDate,
		Address:   params.Address,
		Phone:     params.Phone,
		Status:    "active",
		RoleID:    role.ID,
	})
	if err != nil {
		return "", err
	}

	return s.GenerateToken(id, "user")
}

func (s *svc) Logout(ctx context.Context) error {
	tokenString := ctx.Value(tokenContextKey).(string)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetString("JWT_SECRET", "secret")), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return ErrInvalidToken
	}
	_, err = s.repo.CreateTokenBlacklist(ctx, tokenString)

	if err != nil {
		return err
	}
	return nil
}

func (s *svc) ValidateToken(ctx context.Context, params ValidateTokenParams) (map[string]interface{}, error) {
	token, err := jwt.Parse(params.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetString("JWT_SECRET", "secret")), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		_, err = s.repo.CreateTokenBlacklist(ctx, params.Token)
		if err != nil {
			return nil, err
		}
		return nil, ErrInvalidToken
	}
	blacklisted, err := s.IsTokenBlacklisted(ctx, params.Token)
	if err != nil {
		return nil, err
	}
	if blacklisted {
		return nil, ErrInvalidToken
	}
	data, err := s.extractTokenData(token)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.FindUserById(ctx, data["id"].(string))
	if err != nil {
		return nil, ErrInvalidToken
	}
	return data, nil
}

func (s *svc) extractTokenData(token *jwt.Token) (map[string]interface{}, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid sub claim")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid role claim")
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid exp claim")
	}

	return map[string]interface{}{"id": sub, "role": role, "exp": exp}, nil
}

func (s *svc) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	_, err := s.repo.FindTokenBlacklist(ctx, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
