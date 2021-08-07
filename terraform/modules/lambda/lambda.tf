data "archive_file" "zip" {
  output_path = "${path.module}/${var.archive}"
  type = "zip"

  source_dir = var.source_dir

  excludes = []
}

resource "aws_iam_role" "lambda_exec_role" {
  name = "lambda_${var.lambda_name}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "test_logging" {
  name = "lambda_${var.lambda_name}"
  path = "/"
  description = "IAM policy for logging from ${var.lambda_name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_apigatewayv2_api" "lambda" {
  name = var.lambda_name
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "lambda" {
  name = var.lambda_name
  api_id = aws_apigatewayv2_api.lambda.id

  auto_deploy = true
}

resource "aws_iam_role_policy_attachment" "test_logs" {
  role = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.test_logging.arn
}

resource "aws_lambda_function" "item_lambdas" {
  for_each = var.lambdas
  function_name = "${var.lambda_name}_${each.key}"
  handler = each.key
  runtime = "go1.x"
  role = aws_iam_role.lambda_exec_role.arn

  filename = data.archive_file.zip.output_path
  source_code_hash = data.archive_file.zip.output_base64sha256
}

resource "aws_cloudwatch_log_group" "lambda" {
  for_each = aws_lambda_function.item_lambdas
  name = "/aws/lambda/${each.value.function_name}"
}

resource "aws_lambda_permission" "lambda" {
  for_each = aws_lambda_function.item_lambdas
  action = "lambda:InvokeFunction"
  function_name = each.value.function_name
  principal = "apigateway.amazonaws.com"
  statement_id = "AllowAPIGatewayInvoke"

  source_arn = "${aws_apigatewayv2_api.lambda.execution_arn}/*/*"
}

resource "aws_apigatewayv2_integration" "lambda" {
  for_each = aws_lambda_function.item_lambdas
  api_id = aws_apigatewayv2_api.lambda.id
  integration_type = "AWS_PROXY"

  connection_type = "INTERNET"
  description = "${var.lambda_name} HTTP integration"
  integration_method = "POST"
  integration_uri = each.value.invoke_arn
  passthrough_behavior = "WHEN_NO_MATCH"

  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "lambda" {
  for_each = var.lambdas
  api_id = aws_apigatewayv2_api.lambda.id
  route_key = each.value.route

  target = "integrations/${aws_apigatewayv2_integration.lambda[each.key].id}"
}
