package dto

import (
	"time"

	"github.com/alexander/go-api-hex/internal/domain"
	"github.com/google/uuid"
)

// CreateUserRequest es el payload de creación / registro.
type CreateUserRequest struct {
	Nombre   string `json:"nombre" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateUserRequest permite actualizar nombre y/o email.
type UpdateUserRequest struct {
	Nombre string `json:"nombre" binding:"omitempty,min=2"`
	Email  string `json:"email" binding:"omitempty,email"`
}

// LoginRequest credenciales para autenticar.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse expone los campos seguros del usuario (sin password).
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Nombre    string    `json:"nombre"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse devuelve el JWT generado.
type LoginResponse struct {
	Token string `json:"token"`
}

// ErrorResponse formato uniforme para errores.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ToUserResponse mapea una entidad de dominio a su DTO público.
func ToUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Nombre:    u.Nombre,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
