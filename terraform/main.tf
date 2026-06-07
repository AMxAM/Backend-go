resource "aws_iam_role" "lambda_role" {
  name = "go-api-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${var.lambda_function_name}"
  retention_in_days = 7
}

resource "aws_s3_bucket" "uploads" {
  bucket = "alexander-guzman-043025084823-uploads"

  tags = {
    Name = "UploadsBucket"
  }
}

resource "aws_iam_role_policy" "lambda_s3_policy" {
  name = "lambda-s3-policy"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject"
      ]
      Resource = [
        "${aws_s3_bucket.uploads.arn}/*"
      ]
    }]
  })
}

resource "aws_lambda_function" "api" {

  function_name = var.lambda_function_name

  role = aws_iam_role.lambda_role.arn

  runtime = "provided.al2023"

  handler = "bootstrap"

  filename         = "lambda.zip"
  source_code_hash = filebase64sha256("lambda.zip")

  timeout = 30

  environment {
    variables = {
      DATABASE_URL   = var.database_url
      JWT_SECRET     = var.jwt_secret
      AWS_BUCKET_NAME = aws_s3_bucket.uploads.bucket
    }
  }
}

resource "aws_apigatewayv2_api" "api" {

  name          = "go-api-gateway"

  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "lambda" {

  api_id = aws_apigatewayv2_api.api.id

  integration_type = "AWS_PROXY"

  integration_uri = aws_lambda_function.api.invoke_arn

  payload_format_version = "1.0"
}

resource "aws_apigatewayv2_route" "default" {

  api_id = aws_apigatewayv2_api.api.id

  route_key = "$default"

  target = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

resource "aws_apigatewayv2_stage" "prod" {

  api_id = aws_apigatewayv2_api.api.id

  name = "$default"

  auto_deploy = true
}

resource "aws_lambda_permission" "api_gateway" {

  statement_id = "AllowExecutionFromAPIGateway"

  action = "lambda:InvokeFunction"

  function_name = aws_lambda_function.api.function_name

  principal = "apigateway.amazonaws.com"
}

