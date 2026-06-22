package dto

type NotificationRequest struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}