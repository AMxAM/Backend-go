package services

import (
	"context"
	"errors"

	"github.com/alexander/go-api-hex/internal/application/ports"
	"github.com/alexander/go-api-hex/internal/domain"
	"github.com/google/uuid"
)

// userService implementa ports.UserService.
// Depende de los puertos de salida, no de implementaciones concretas.
type userService struct {
	repo   ports.UserRepository
	hasher ports.PasswordHasher
}

func NewUserService(repo ports.UserRepository, hasher ports.PasswordHasher) ports.UserService {
	return &userService{repo: repo, hasher: hasher}
}

func (s *userService) Create(ctx context.Context, nombre, email, password string) (*domain.User, error) {
	// Verificar que el email no exista
	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(nombre, email, hash)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) List(ctx context.Context) ([]*domain.User, error) {
	return s.repo.List(ctx)
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, nombre, email string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if nombre != "" {
		user.Nombre = nombre
	}
	if email != "" && email != user.Email {
		// Verificar que el nuevo email no esté en uso
		other, err := s.repo.GetByEmail(ctx, email)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		if other != nil && other.ID != id {
			return nil, domain.ErrUserAlreadyExists
		}
		user.Email = email
	}
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	// Asegurar existencia antes de borrar
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
