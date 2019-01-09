data "aws_caller_identity" "current" {}

provider "aws" {
  alias  = "frankfurt"
  region = "eu-central-1"
}

provider "aws" {
  alias  = "virginia"
  region = "us-east-1"
}

variable "name" {
  type    = "string"
  default = "testname"
}

variable "root-domain" {
  type    = "string"
  default = "test.com"
}

variable "email-domain" {
  type    = "string"
  default = "test.com"
}

variable "cf-domain" {
  type    = "string"
  default = "www.test.com"
}

variable "auth-domain" {
  type    = "string"
  default = "auth.test.com"
}

variable "default-region" {
  type    = "string"
  default = "eu-central-1"
}
