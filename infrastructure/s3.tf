data "aws_iam_policy_document" "cf-bucket-policy" {
  statement {
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.cf-bucket.arn}/*"]

    principals {
      type        = "AWS"
      identifiers = ["${aws_cloudfront_origin_access_identity.origin_access_identity.iam_arn}"]
    }
  }

  statement {
    actions   = ["s3:ListBucket"]
    resources = ["${aws_s3_bucket.cf-bucket.arn}"]

    principals {
      type        = "AWS"
      identifiers = ["${aws_cloudfront_origin_access_identity.origin_access_identity.iam_arn}"]
    }
  }
}

resource "aws_s3_bucket_policy" "cf-bucket" {
  bucket = "${aws_s3_bucket.cf-bucket.id}"
  policy = "${data.aws_iam_policy_document.cf-bucket-policy.json}"
}

resource "aws_s3_bucket" "cf-bucket" {
  bucket = "${var.cf-domain}"
  acl    = "private"
}

resource "aws_s3_bucket_object" "loginIndex" {
  depends_on = [
    "aws_s3_bucket.cf-bucket",
  ]

  bucket       = "${var.cf-domain}"
  key          = "/login/index.html"
  source       = "./loginFiles/index.html"
  content_type = "text/html"
  etag         = "${md5(file("./loginFiles/index.html"))}"
}

resource "aws_s3_bucket_object" "loginUnauthorized" {
  depends_on = [
    "aws_s3_bucket.cf-bucket",
  ]

  bucket       = "${var.cf-domain}"
  key          = "/login/unauthorized.html"
  source       = "./loginFiles/unauthorized.html"
  content_type = "text/html"
  etag         = "${md5(file("./loginFiles/unauthorized.html"))}"
}

resource "aws_s3_bucket" "zip-bucket" {
  bucket = "${var.name}-zip-bucket"
  acl    = "private"

  lifecycle_rule {
    enabled = true

    expiration {
      days = 5
    }
  }
}

resource "aws_s3_bucket" "lambda-bucket" {
  provider = "aws.virginia"
  region   = "us-east-1"
  bucket   = "${var.name}-lambda-zip-bucket"
  acl      = "private"

  lifecycle_rule {
    enabled = true

    expiration {
      days = 30
    }
  }
}
