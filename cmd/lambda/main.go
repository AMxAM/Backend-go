package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/alexander/go-api-hex/internal/application/services"
	"github.com/alexander/go-api-hex/internal/infrastructure/auth"
	"github.com/alexander/go-api-hex/internal/infrastructure/config"
	"github.com/alexander/go-api-hex/internal/infrastructure/hasher"
	httpinfra "github.com/alexander/go-api-hex/internal/infrastructure/http"
	"github.com/alexander/go-api-hex/internal/infrastructure/persistence"
	"github.com/alexander/go-api-hex/internal/infrastructure/storage"
	"github.com/alexander/go-api-hex/internal/infrastructure/notifications"
	
)

var ginLambda *ginadapter.GinLambda

func init() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := persistence.NewPostgresDB(cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}

	userRepo := persistence.NewUserRepository(db)

	pwdHasher := hasher.NewBcryptHasher()

	tokenSvc := auth.NewJWTService(
		cfg.JWTSecret,
		cfg.JWTExpire,
	)

	s3Storage, err := storage.NewS3Storage(
		cfg.AWSBucket,
	)

	if err != nil {
		log.Fatal(err)
	}

	snsService, err := notifications.NewSNSService()

	if err != nil {
		log.Fatal(err)
	}

	userSvc := services.NewUserService(
		userRepo,
		pwdHasher,
	)

	authSvc := services.NewAuthService(
		userRepo,
		pwdHasher,
		tokenSvc,
		userSvc,
	)

	router := httpinfra.NewRouter(
		userSvc,
		authSvc,
		tokenSvc,
		s3Storage,
		snsService,
	)

	ginLambda = ginadapter.New(router)
}

func Handler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (
	events.APIGatewayProxyResponse,
	error,
) {
	return ginLambda.ProxyWithContext(
		ctx,
		req,
	)
}

func main() {
	lambda.Start(Handler)
}