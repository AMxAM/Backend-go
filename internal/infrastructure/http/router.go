package http

import (
	"github.com/alexander/go-api-hex/internal/infrastructure/storage"
	"github.com/alexander/go-api-hex/internal/application/ports"
	"github.com/alexander/go-api-hex/internal/infrastructure/http/handlers"
	"github.com/alexander/go-api-hex/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/alexander/go-api-hex/internal/infrastructure/notifications"
)


func NewRouter(
	userSvc ports.UserService,
	authSvc ports.AuthService,
	tokens ports.TokenService,
	s3Storage *storage.S3Storage,
	snsService *notifications.SNSService,
) *gin.Engine {
	r := gin.Default()

	

	userH := handlers.NewUserHandler(userSvc)

	authH := handlers.NewAuthHandler(authSvc)

	uploadH := handlers.NewUploadHandler(
		s3Storage,
	)

	notificationH := handlers.NewNotificationHandler(
		snsService,
	)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	// Rutas públicas de autenticación
	auth := api.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
	}

	// Rutas protegidas por JWT
	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware(tokens))
	{
		users.POST("", userH.Create)
		users.GET("", userH.List)
		users.GET("/:id", userH.GetByID)
		users.PUT("/:id", userH.Update)
		users.DELETE("/:id", userH.Delete)
	}

	

	notifications := api.Group("/notifications")
	{
		notifications.POST(
			"/send",
			notificationH.Send,
		)
	}	

	api.POST("/upload", uploadH.Upload)
	
	return r
}
