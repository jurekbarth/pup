resource "aws_sqs_queue" "queue" {
  name                       = "${var.name}-queue"
  visibility_timeout_seconds = 900
}
