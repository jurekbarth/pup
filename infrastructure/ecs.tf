resource "aws_iam_role" "iam-ecs-task" {
  name = "${var.name}-iam-ecs-task"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "ecs-logging" {
  name        = "${var.name}-iam-ecs-logging"
  path        = "/"
  description = "IAM policy for logging from a ecs fargate"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "logs:CreateLogStream",
          "logs:DescribeLogStreams",
          "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    },
    {
      "Effect": "Allow",
      "Action": "logs:CreateLogGroup",
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-logs" {
  role       = "${aws_iam_role.iam-ecs-task.name}"
  policy_arn = "${aws_iam_policy.ecs-logging.arn}"
}

resource "aws_cloudwatch_log_group" "worker" {
  name              = "/aws/ecs/${var.name}-worker"
  retention_in_days = 3
}

resource "aws_ecs_cluster" "ecs" {
  name = "${var.name}-cluster"
}

resource "aws_ecs_task_definition" "worker" {
  depends_on = [
    "aws_cloudwatch_log_group.worker",
  ]

  family = "${var.name}-worker"

  network_mode = "awsvpc"

  task_role_arn      = "${aws_iam_role.iam-ecs-task.arn}"
  execution_role_arn = "${aws_iam_role.iam-ecs-task.arn}"

  container_definitions = <<EOF
[
  {
    "name": "worker",
    "memory": 300,
    "image": "docker.io/jurekbarth/worker:v0.0.5",
    "environment": [
      {
        "name": "PUP_DOWNLOAD_DIR",
        "value": "."
      },
      {
        "name": "PUP_UNZIP_DIR",
        "value": "./output"
      },
      {
        "name": "PUP_COGNITO_POOL_ID",
        "value": "${aws_cognito_user_pool.pool.id}"
      },
      {
        "name": "PUP_COGNITO_CLIENT_SPA_SUBDOMAIN",
        "value": "${var.name}"
      },
      {
        "name": "PUP_COGNITO_CLIENT_SPA_CLIENT_ID",
        "value": "${aws_cognito_user_pool_client.spa-client.id}"
      },
      {
        "name": "PUP_COGNITO_CLIENT_BACKEND_CLIENT_ID",
        "value": "${aws_cognito_user_pool_client.backend-client.id}"
      },
      {
        "name": "PUP_CF_DOMAIN",
        "value": "${var.cf-domain}"
      },
      {
        "name": "PUP_CF_ID",
        "value": "${aws_cloudfront_distribution.s3_distribution.id}"
      },
      {
        "name": "PUP_DEFAULT_REGION",
        "value": "${var.default-region}"
      },
      {
        "name": "PUP_EMAIL_DOMAIN",
        "value": "${var.email-domain}"
      },
      {
        "name": "PUP_S3_FROM_BUCKET",
        "value": "${aws_s3_bucket.zip-bucket.bucket}"
      },
      {
        "name": "PUP_S3_DESTINATION_BUCKET",
        "value": "${aws_s3_bucket.cf-bucket.bucket}"
      },
      {
        "name": "PUP_S3_LAMBDA_BUCKET",
        "value": "${aws_s3_bucket.lambda-bucket.bucket}"
      },
      {
        "name": "PUP_LAMBDA_EDGE_ARN",
        "value": "${aws_lambda_function.edge.arn}"
      },
      {
        "name": "PUP_SQS_URI",
        "value": "${aws_sqs_queue.queue.id}"
      },
      {
        "name": "PUP_DDB_RULES_TABLE",
        "value": "${aws_dynamodb_table.rules-table.name}"
      },
      {
        "name": "PUP_DDB_LOGS_TABLE",
        "value": "${aws_dynamodb_table.logs-table.name}"
      }
    ],
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
        "awslogs-group": "/aws/ecs/${var.name}-worker",
        "awslogs-region": "${var.default-region}",
        "awslogs-stream-prefix": "worker"
      }
    }
  }
]
EOF

  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512

  volume {
    name = "storage-volume"
  }
}
