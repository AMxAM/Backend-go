package domain

import "errors"

// Errores de dominio. Las capas externas los mapean a códigos HTTP.
var (
	ErrUserNotFound      = errors.New("usuario no encontrado")
	ErrUserAlreadyExists = errors.New("el usuario ya existe")
	ErrInvalidUserData   = errors.New("datos de usuario inválidos")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUnauthorized      = errors.New("no autorizado")
)
