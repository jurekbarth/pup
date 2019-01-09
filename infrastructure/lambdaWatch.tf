resource "aws_iam_role" "iam-watch" {
  name = "${var.name}-iam-watch"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.watch.arn}"
  principal     = "s3.amazonaws.com"
  source_arn    = "${aws_s3_bucket.zip-bucket.arn}"
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = "${aws_s3_bucket.zip-bucket.id}"

  lambda_function {
    lambda_function_arn = "${aws_lambda_function.watch.arn}"
    events              = ["s3:ObjectCreated:*"]
    filter_suffix       = ".zip"
  }
}

resource "aws_cloudwatch_log_group" "watch" {
  name              = "/aws/lambda/${aws_lambda_function.watch.function_name}"
  retention_in_days = 3
}

resource "aws_iam_role_policy_attachment" "watch-logs" {
  role       = "${aws_iam_role.iam-watch.name}"
  policy_arn = "${aws_iam_policy.lambda-logging.arn}"
}

resource "aws_iam_role_policy_attachment" "watch-sqs" {
  role       = "${aws_iam_role.iam-watch.name}"
  policy_arn = "${aws_iam_policy.lambda-sqs.arn}"
}

resource "aws_iam_role_policy_attachment" "watch-ecs" {
  role       = "${aws_iam_role.iam-watch.name}"
  policy_arn = "${aws_iam_policy.lambda-ecs-run.arn}"
}

data "archive_file" "watch" {
  type        = "zip"
  source_dir  = "${path.module}/lambdaWatch"
  output_path = "${path.module}/.archive/lambdaWatch.zip"
}

resource "aws_lambda_function" "watch" {
  provider         = "aws.frankfurt"
  filename         = "${data.archive_file.watch.output_path}"
  source_code_hash = "${data.archive_file.watch.output_base64sha256}"
  function_name    = "${var.name}-watch-function"
  role             = "${aws_iam_role.iam-watch.arn}"
  handler          = "index.handler"
  runtime          = "nodejs8.10"
  publish          = true
  memory_size      = 128

  environment {
    variables = {
      SQS        = "${aws_sqs_queue.queue.id}"
      WORKERNAME = "${aws_ecs_task_definition.worker.family}"
      CLUSTER    = "${aws_ecs_cluster.ecs.name}"
      SUBNET     = "${aws_subnet.main.id}"
    }
  }
}
