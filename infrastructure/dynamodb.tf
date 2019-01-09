resource "aws_dynamodb_table" "rules-table" {
  name         = "${var.name}-rules-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "baseUri"

  attribute {
    name = "baseUri"
    type = "S"
  }
}

resource "aws_dynamodb_table" "logs-table" {
  name         = "${var.name}-logs-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "deployid"

  attribute {
    name = "deployid"
    type = "S"
  }
}
