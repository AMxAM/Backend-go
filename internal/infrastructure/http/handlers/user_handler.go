package handlers

import (
	"net/http"

	"github.com/alexander/go-api-hex/internal/application/ports"
	"github.com/alexander/go-api-hex/internal/infrastructure/http/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	users ports.UserService
}

func NewUserHandler(users ports.UserService) *UserHandler {
	return &UserHandler{users: users}
}

// Create POST /api/v1/users
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	user, err := h.users.Create(c.Request.Context(), req.Nombre, req.Email, req.Password)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

// GetByID GET /api/v1/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}
	user, err := h.users.GetByID(c.Request.Context(), id)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// List GET /api/v1/users
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.users.List(c.Request.Context())
	if err != nil {
		mapError(c, err)
		return
	}
	resp := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, dto.ToUserResponse(u))
	}
	c.JSON(http.StatusOK, resp)
}

// Update PUT /api/v1/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	user, err := h.users.Update(c.Request.Context(), id, req.Nombre, req.Email)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// Delete DELETE /api/v1/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}
	if err := h.users.Delete(c.Request.Context(), id); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
