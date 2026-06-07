package ports

// PasswordHasher es el puerto de salida para hashear y verificar contraseñas.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
