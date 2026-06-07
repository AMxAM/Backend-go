package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config agrupa todas las variables de entorno necesarias.
type Config struct {
	HTTPPort    string
	DatabaseURL string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string
	JWTExpire time.Duration
}

// Load lee las variables de entorno (usa .env si está presente).
func Load() (*Config, error) {
	// Si existe un .env, lo carga (ignora error si no existe).
	_ = godotenv.Load()

	cfg := &Config{
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "231298"),
		DBName:     getEnv("DB_NAME", "users_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", "cambia-este-secreto"),
	}

	expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	cfg.JWTExpire = time.Duration(expireHours) * time.Hour

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET es obligatorio")
	}
	return cfg, nil
}

// DSN construye la cadena de conexión para PostgreSQL.
func (c *Config) DSN() string {

	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
