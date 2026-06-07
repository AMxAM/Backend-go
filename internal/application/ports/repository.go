package ports

import (
	"context"

	"github.com/alexander/go-api-hex/internal/domain"
	"github.com/google/uuid"
)

// UserRepository es el puerto de salida hacia la persistencia.
// La capa de dominio/aplicación depende de esta interfaz, NO de una BD concreta.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
