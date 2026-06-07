package services

import (
	"context"
	"errors"

	"github.com/alexander/go-api-hex/internal/application/ports"
	"github.com/alexander/go-api-hex/internal/domain"
)

type authService struct {
	repo   ports.UserRepository
	hasher ports.PasswordHasher
	tokens ports.TokenService
	users  ports.UserService
}

func NewAuthService(
	repo ports.UserRepository,
	hasher ports.PasswordHasher,
	tokens ports.TokenService,
	users ports.UserService,
) ports.AuthService {
	return &authService{repo: repo, hasher: hasher, tokens: tokens, users: users}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}
	if err := s.hasher.Compare(user.Password, password); err != nil {
		return "", domain.ErrInvalidCredentials
	}
	return s.tokens.Generate(user.ID, user.Email)
}

func (s *authService) Register(ctx context.Context, nombre, email, password string) (*domain.User, error) {
	return s.users.Create(ctx, nombre, email, password)
}
