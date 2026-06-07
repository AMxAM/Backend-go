package domain

import (
	"time"

	"github.com/google/uuid"
)

// User es la entidad principal del dominio.
// No contiene tags JSON ni de DB: es independiente de cualquier framework.
type User struct {
	ID        uuid.UUID
	Nombre    string
	Email     string
	Password  string // hash de la contraseña
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser construye un nuevo usuario validando invariantes básicos.
func NewUser(nombre, email, passwordHash string) (*User, error) {
	if nombre == "" {
		return nil, ErrInvalidUserData
	}
	if email == "" {
		return nil, ErrInvalidUserData
	}
	if passwordHash == "" {
		return nil, ErrInvalidUserData
	}
	now := time.Now().UTC()
	return &User{
		ID:        uuid.New(),
		Nombre:    nombre,
		Email:     email,
		Password:  passwordHash,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
