resource "aws_iam_role" "iam-edge" {
  name = "${var.name}-iam-edge"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": [
          "edgelambda.amazonaws.com",
          "lambda.amazonaws.com"
        ]
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_cloudwatch_log_group" "edge" {
  provider          = "aws.virginia"
  name              = "/aws/lambda/${aws_lambda_function.edge.function_name}"
  retention_in_days = 3
}

resource "aws_iam_role_policy_attachment" "edge-logs" {
  role       = "${aws_iam_role.iam-edge.name}"
  policy_arn = "${aws_iam_policy.lambda-logging.arn}"
}

data "archive_file" "edge" {
  type        = "zip"
  source_dir  = "${path.module}/lambdaEdge"
  output_path = "${path.module}/.archive/lambdaEdge.zip"
}

resource "aws_lambda_function" "edge" {
  provider         = "aws.virginia"
  filename         = "${data.archive_file.edge.output_path}"
  source_code_hash = "${data.archive_file.edge.output_base64sha256}"
  function_name    = "${var.name}-edge-function"
  role             = "${aws_iam_role.iam-edge.arn}"
  handler          = "index.handler"
  runtime          = "nodejs8.10"
  publish          = true
  memory_size      = 128
}
