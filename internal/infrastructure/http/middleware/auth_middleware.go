package middleware

import (
	"net/http"
	"strings"

	"github.com/alexander/go-api-hex/internal/application/ports"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey = "userID"
	ContextEmailKey  = "userEmail"
)

// AuthMiddleware valida el header Authorization: Bearer <token>.
func AuthMiddleware(tokens ports.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token faltante"})
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			return
		}
		claims, err := tokens.Validate(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			return
		}
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextEmailKey, claims.Email)
		c.Next()
	}
}
