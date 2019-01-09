resource "aws_cognito_user_pool" "pool" {
  provider = "aws.frankfurt"
  name     = "${var.name}-pool"

  password_policy {
    minimum_length    = 8
    require_lowercase = false
    require_numbers   = false
    require_symbols   = false
    require_uppercase = false
  }

  admin_create_user_config {
    allow_admin_create_user_only = true
  }

  provisioner "local-exec" {
    command = "node scripts/setSettingsLambdaEdge.js ${var.default-region} ${aws_cognito_user_pool.pool.id} ${var.name} ${var.cf-domain}"
  }
}

resource "aws_cognito_user_pool_domain" "main" {
  domain = "${var.name}"

  # domain          = "${var.auth-domain}"
  # certificate_arn = "${data.aws_acm_certificate.cert.arn}"
  user_pool_id = "${aws_cognito_user_pool.pool.id}"
}

resource "aws_cognito_user_pool_client" "spa-client" {
  name                                 = "${var.name}-spa-client"
  user_pool_id                         = "${aws_cognito_user_pool.pool.id}"
  supported_identity_providers         = ["COGNITO"]
  explicit_auth_flows                  = ["USER_PASSWORD_AUTH"]
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["implicit"]
  allowed_oauth_scopes                 = ["phone", "email", "openid", "profile", "aws.cognito.signin.user.admin"]
  callback_urls                        = ["https://${var.cf-domain}/login/index.html"]
  logout_urls                          = ["https://${var.cf-domain}/login/index.html"]

  provisioner "local-exec" {
    command = "node scripts/setClientIdLambdaEdge.js ${aws_cognito_user_pool_client.spa-client.id}"
  }
}

resource "aws_cognito_user_pool_client" "backend-client" {
  name                = "${var.name}-backend-client"
  user_pool_id        = "${aws_cognito_user_pool.pool.id}"
  explicit_auth_flows = ["ADMIN_NO_SRP_AUTH"]
}
