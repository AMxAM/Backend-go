output "lambda_role_arn" {
  value = aws_iam_role.lambda_role.arn
}

output "api_gateway_url" {
  value = aws_apigatewayv2_api.api.api_endpoint
}