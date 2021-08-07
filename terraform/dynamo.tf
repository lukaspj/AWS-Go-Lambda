resource "aws_dynamodb_table" "test_table" {
  name = "sample_table"
  hash_key = "UserId"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "UserId"
    type = "S"
  }
}

resource "aws_iam_policy" "test_dynamo" {
  name = "sample_dynamo"
  path = "/"
  description = "IAM policy for accessing dynamo from the sample lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
        {
            "Action": "dynamodb:*",
            "Resource": [
                "arn:aws:dynamodb::724375571299:global-table/*",
                "arn:aws:dynamodb:eu-north-1:724375571299:table/sample_dynamo/backup/*",
                "arn:aws:dynamodb:eu-north-1:724375571299:table/sample_dynamo/export/*",
                "arn:aws:dynamodb:eu-north-1:724375571299:table/sample_dynamo/stream/*",
                "arn:aws:dynamodb:eu-north-1:724375571299:table/sample_dynamo"
            ],
            "Effect": "Allow"
        }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "test_dynamo" {
  role = module.testtype_lambda.lambda_iam_role.name
  policy_arn = aws_iam_policy.test_dynamo.arn
}