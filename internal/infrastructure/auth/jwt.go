package auth

import (
	"errors"
	"time"

	"github.com/alexander/go-api-hex/internal/application/ports"
	"github.com/alexander/go-api-hex/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtService struct {
	secret []byte
	expire time.Duration
}

// NewJWTService crea un adaptador de tokens JWT.
func NewJWTService(secret string, expire time.Duration) ports.TokenService {
	return &jwtService{secret: []byte(secret), expire: expire}
}

type customClaims struct {
	UserID string `json:"uid"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *jwtService) Generate(userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	claims := customClaims{
		UserID: userID.String(),
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expire)),
			Subject:   userID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *jwtService) Validate(tokenStr string) (*ports.TokenClaims, error) {
	claims := &customClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, domain.ErrUnauthorized
	}
	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return &ports.TokenClaims{UserID: id, Email: claims.Email}, nil
}
