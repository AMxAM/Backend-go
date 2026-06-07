package ports

import "github.com/google/uuid"

// TokenService es el puerto de salida para generar/validar tokens JWT.
type TokenService interface {
	Generate(userID uuid.UUID, email string) (string, error)
	Validate(token string) (*TokenClaims, error)
}

// TokenClaims es el payload que viaja en el JWT.
type TokenClaims struct {
	UserID uuid.UUID
	Email  string
}
