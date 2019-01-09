resource "aws_iam_policy" "ecs-s3-cf-bucket" {
  name        = "${var.name}-iam-s3-cf-bucket"
  path        = "/"
  description = "IAM policy for delete, upload from ecs to s3 cf bucket"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:ListBucketMultipartUploads",
        "s3:ListBucket",
        "s3:DeleteObject",
        "s3:ListMultipartUploadParts"
      ],
      "Resource": [
        "arn:aws:s3:::${aws_s3_bucket.cf-bucket.bucket}",
        "arn:aws:s3:::${aws_s3_bucket.cf-bucket.bucket}/*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-s3-cf-bucket" {
  role       = "${aws_iam_role.iam-ecs-task.name}"
  policy_arn = "${aws_iam_policy.ecs-s3-cf-bucket.arn}"
}

resource "aws_iam_policy" "ecs-s3-zip-bucket" {
  name        = "${var.name}-iam-s3-zip-client-bucket"
  path        = "/"
  description = "IAM policy for download zip bucket files from ecs"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::${aws_s3_bucket.zip-bucket.bucket}/*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-s3-zip-bucket" {
  role       = "${aws_iam_role.iam-ecs-task.name}"
  policy_arn = "${aws_iam_policy.ecs-s3-zip-bucket.arn}"
}

resource "aws_iam_policy" "ecs-s3-lambda-zip-bucket" {
  name        = "${var.name}-iam-s3-lambda-zip-bucket"
  path        = "/"
  description = "IAM policy for download zip bucket files from ecs"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "s3:PutObject",
          "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::${aws_s3_bucket.lambda-bucket.bucket}/*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-s3-lambda-zip-bucket" {
  role       = "${aws_iam_role.iam-ecs-task.name}"
  policy_arn = "${aws_iam_policy.ecs-s3-lambda-zip-bucket.arn}"
}

resource "aws_iam_policy" "ecs-sqs" {
  name        = "${var.name}-iam-sqs"
  path        = "/"
  description = "IAM policy for read delete sqs from ecs"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "VisualEditor0",
      "Effect": "Allow",
      "Action": [
          "sqs:DeleteMessage",
          "sqs:ReceiveMessage"
      ],
      "Resource": "arn:aws:sqs:${var.default-region}:${data.aws_caller_identity.current.account_id}:${aws_sqs_queue.queue.name}"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-sqs" {
  role       = "${aws_iam_role.iam-ecs-task.name}"
  policy_arn = "${aws_iam_policy.ecs-sqs.arn}"
}

resource "aws_iam_policy" "ecs-lambda" {
  name        = "${var.name}-iam-lambda"
  path        = "/"
  description = "IAM policy for updating lambda function code from ecs"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "lambda:UpdateFunctionCode",
        "lambda:GetFunction",
        "lambda:PublishVersion",
        "lambda:EnableReplication"
      ],
      "Resource": "arn:aws:lambda:*:*:function:*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-lambda" {
  role       = "${aws_iam_role.iam-ecs-task.name}"
  policy_arn = "${aws_iam_policy.ecs-lambda.arn}"
}

resource "aws_iam_policy" "ecs-cognito" {
  name        = "${var.name}-iam-cognito"
  path        = "/"
  description = "IAM policy for updating cognito from ecs"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "cognito-idp:AdminDeleteUser",
          "cognito-idp:AdminInitiateAuth",
          "cognito-idp:AdminCreateUser",
          "cognito-idp:CreateGroup",
          "cognito-idp:AdminAddUserToGroup",
          "cognito-idp:GetGroup",
          "cognito-idp:AdminRespondToAuthChallenge",
          "cognito-idp:AdminGetUser"
      ],
      "Resource": "arn:aws:cognito-idp:${var.default-region}:${data.aws_caller_identity.current.account_id}:userpool/${aws_cognito_user_pool.pool.id}"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-cognito" {
  role       = "${aws_iam_role.iam-ecs-task.name}"
  policy_arn = "${aws_iam_policy.ecs-cognito.arn}"
}

resource "aws_iam_policy" "ecs-dynamodb-rules" {
  name        = "${var.name}-iam-dynamodb-rules"
  path        = "/"
  description = "IAM policy for updating dynamo db rules from ecs"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "dynamodb:GetItem",
          "dynamodb:Scan",
          "dynamodb:UpdateItem"
      ],
      "Resource": [
          "arn:aws:dynamodb:${var.default-region}:${data.aws_caller_identity.current.account_id}:table/${aws_dynamodb_table.rules-table.name}/index/*",
          "arn:aws:dynamodb:${var.default-region}:${data.aws_caller_identity.current.account_id}:table/${aws_dynamodb_table.rules-table.name}"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-dynamodb-rules" {
  role       = "${aws_iam_role.iam-ecs-task.name}"
  policy_arn = "${aws_iam_policy.ecs-dynamodb-rules.arn}"
}

resource "aws_iam_policy" "ecs-dynamodb-logs" {
  name        = "${var.name}-iam-dynamodb-logs"
  path        = "/"
  description = "IAM policy for updating dynamo db logs from ecs"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "dynamodb:GetItem",
          "dynamodb:UpdateItem"
      ],
      "Resource": [
          "arn:aws:dynamodb:${var.default-region}:${data.aws_caller_identity.current.account_id}:table/${aws_dynamodb_table.logs-table.name}/index/*",
          "arn:aws:dynamodb:${var.default-region}:${data.aws_caller_identity.current.account_id}:table/${aws_dynamodb_table.logs-table.name}"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-dynamodb-logs" {
  role       = "${aws_iam_role.iam-ecs-task.name}"
  policy_arn = "${aws_iam_policy.ecs-dynamodb-logs.arn}"
}

resource "aws_iam_policy" "ecs-cloudfront" {
  name        = "${var.name}-iam-cloudfront"
  path        = "/"
  description = "IAM policy for updating cloudfront from ecs"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "cloudfront:UpdateDistribution",
          "cloudfront:GetDistributionConfig",
          "cloudfront:CreateInvalidation"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-cloudfront" {
  role       = "${aws_iam_role.iam-ecs-task.name}"
  policy_arn = "${aws_iam_policy.ecs-cloudfront.arn}"
}
