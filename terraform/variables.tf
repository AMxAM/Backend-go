variable "aws_region" {
  description = "Región de AWS"
  type        = string
}

variable "lambda_function_name" {
  description = "Nombre de la función Lambda"
  type        = string
  default = "go-api-hex"
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