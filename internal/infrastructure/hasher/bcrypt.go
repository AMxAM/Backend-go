package hasher

import (
	"github.com/alexander/go-api-hex/internal/application/ports"
	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct {
	cost int
}

// NewBcryptHasher crea un adaptador de hashing con bcrypt.
func NewBcryptHasher() ports.PasswordHasher {
	return &bcryptHasher{cost: bcrypt.DefaultCost}
}

func (b *bcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *bcryptHasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
