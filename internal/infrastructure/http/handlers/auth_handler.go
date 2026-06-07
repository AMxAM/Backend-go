package handlers

import (
	"errors"
	"net/http"

	"github.com/alexander/go-api-hex/internal/application/ports"
	"github.com/alexander/go-api-hex/internal/domain"
	"github.com/alexander/go-api-hex/internal/infrastructure/http/dto"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth ports.AuthService
}

func NewAuthHandler(auth ports.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	user, err := h.auth.Register(c.Request.Context(), req.Nombre, req.Email, req.Password)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

// Login POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	token, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}

// mapError traduce errores de dominio a códigos HTTP.
func mapError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrUserAlreadyExists):
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrInvalidUserData):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "error interno"})
	}
}
