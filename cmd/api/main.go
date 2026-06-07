package main

import (
	"log"

	"github.com/alexander/go-api-hex/internal/application/services"
	"github.com/alexander/go-api-hex/internal/infrastructure/auth"
	"github.com/alexander/go-api-hex/internal/infrastructure/config"
	"github.com/alexander/go-api-hex/internal/infrastructure/hasher"
	httpinfra "github.com/alexander/go-api-hex/internal/infrastructure/http"
	"github.com/alexander/go-api-hex/internal/infrastructure/persistence"
	"github.com/alexander/go-api-hex/internal/infrastructure/storage"
)

func main() {
	// 1. Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error cargando configuración: %v", err)
	}

	// 2. Conectar a PostgreSQL
	db, err := persistence.NewPostgresDB(cfg.DSN())
	if err != nil {
		log.Fatalf("error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	// 3. Construir adaptadores de salida
	userRepo := persistence.NewUserRepository(db)
	pwdHasher := hasher.NewBcryptHasher()
	tokenSvc := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpire)

	// 4. Crear cliente S3
	s3Storage, err := storage.NewS3Storage(cfg.AWSBucket)
	if err != nil {
		log.Fatalf(
			"error creando cliente S3: %v",
			err,
		)
	}

	// 5. Construir servicios de aplicación
	userSvc := services.NewUserService(userRepo, pwdHasher)
	authSvc := services.NewAuthService(
		userRepo,
		pwdHasher,
		tokenSvc,
		userSvc,
	)

	// 6. Construir router
	router := httpinfra.NewRouter(
		userSvc,
		authSvc,
		tokenSvc,
		s3Storage,
	)

	addr := ":" + cfg.HTTPPort

	log.Printf(
		"API escuchando en %s",
		addr,
	)

	if err := router.Run(addr); err != nil {
		log.Fatalf(
			"error iniciando el servidor: %v",
			err,
		)
	}
}