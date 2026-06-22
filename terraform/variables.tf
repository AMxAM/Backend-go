variable "aws_region" {
  description = "Región de AWS"
  type        = string
}

variable "lambda_function_name" {
  description = "Nombre de la función Lambda"
  type        = string
  default     = "go-api-hex"
}

variable "jwt_secret" {
  description = "JWT Secret"
  type        = string
  sensitive   = true
}

variable "database_url" {
  description = "Neon PostgreSQL URL"
  type        = string
  sensitive   = true
}

variable "sns_topic_name" {
  description = "SNS Topic para notificaciones"
  type        = string
  default     = "notifications-topic"
}

variable "sqs_queue_name" {
  description = "SQS Queue para notificaciones"
  type        = string
  default     = "notifications-queue"
}

variable "smtp_user" {
  description = "Correo SMTP"
  type        = string
  sensitive   = true
}

variable "smtp_password" {
  description = "App Password Gmail"
  type        = string
  sensitive   = true
}