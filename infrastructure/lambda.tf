resource "aws_iam_policy" "lambda-logging" {
  name        = "${var.name}-iam-lambda-logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"

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

resource "aws_iam_policy" "lambda-sqs" {
  name        = "${var.name}-iam-lambda-sqs"
  path        = "/"
  description = "IAM policy for sqs from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.queue.arn}"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "lambda-ecs-run" {
  name        = "${var.name}-iam-lambda-ecs-run"
  path        = "/"
  description = "IAM policy for ecs from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "ecs:RunTask",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
          "iam:PassRole"
      ],
      "Resource": "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/*"
    }
  ]
}
EOF
}
