package ports

import (
	"context"

	"github.com/alexander/go-api-hex/internal/domain"
	"github.com/google/uuid"
)

// UserService es el puerto de entrada (use case) para el CRUD de usuarios.
type UserService interface {
	Create(ctx context.Context, nombre, email, password string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, id uuid.UUID, nombre, email string) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// AuthService es el puerto de entrada para la autenticación.
type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, nombre, email, password string) (*domain.User, error)
}
