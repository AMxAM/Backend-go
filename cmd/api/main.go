package main

import (
	"log"

	"github.com/alexander/go-api-hex/internal/application/services"
	"github.com/alexander/go-api-hex/internal/infrastructure/auth"
	"github.com/alexander/go-api-hex/internal/infrastructure/config"
	"github.com/alexander/go-api-hex/internal/infrastructure/hasher"
	httpinfra "github.com/alexander/go-api-hex/internal/infrastructure/http"
	"github.com/alexander/go-api-hex/internal/infrastructure/persistence"
)

func main() {
	// 1. Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error cargando configuración: %v", err)
	}

	// 2. Conectar a PostgreSQL (adaptador de salida)
	db, err := persistence.NewPostgresDB(cfg.DSN())
	if err != nil {
		log.Fatalf("error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	// 3. Construir adaptadores de salida (driven adapters)
	userRepo := persistence.NewUserRepository(db)
	pwdHasher := hasher.NewBcryptHasher()
	tokenSvc := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpire)

	// 4. Construir servicios de aplicación (use cases)
	userSvc := services.NewUserService(userRepo, pwdHasher)
	authSvc := services.NewAuthService(userRepo, pwdHasher, tokenSvc, userSvc)

	// 5. Construir adaptador de entrada (driving adapter: HTTP)
	router := httpinfra.NewRouter(userSvc, authSvc, tokenSvc)

	addr := ":" + cfg.HTTPPort
	log.Printf("API escuchando en %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("error iniciando el servidor: %v", err)
	}
}
