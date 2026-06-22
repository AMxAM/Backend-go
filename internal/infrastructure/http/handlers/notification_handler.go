package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alexander/go-api-hex/internal/infrastructure/http/dto"
	"github.com/alexander/go-api-hex/internal/infrastructure/notifications"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	sns *notifications.SNSService
}

func NewNotificationHandler(
	sns *notifications.SNSService,
) *NotificationHandler {

	return &NotificationHandler{
		sns: sns,
	}
}

func (h *NotificationHandler) Send(
	c *gin.Context,
) {

	var req dto.NotificationRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	payload, _ := json.Marshal(req)

	if err := h.sns.Publish(
		string(payload),
	); err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "Mensaje enviado a SNS",
		},
	)
}