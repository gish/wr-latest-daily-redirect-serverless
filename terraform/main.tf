data "aws_secretsmanager_secret_version" "this" {
  secret_id = "wr-latest-daily-redirect"
}

locals {
  application_name      = "wr-latest-daily-redirect"
  region                = "eu-north-1"
  account_id            = "852264810958"
  log_retention_in_days = 3
  domain                = "dailytest.stafre.se"
  secrets               = jsondecode(data.aws_secretsmanager_secret_version.this.secret_string)
}

provider "aws" {
  region = local.region
  default_tags {
    tags = {
      Terraform   = "yes"
      Application = "${local.application_name}"
    }
  }
}

terraform {
  backend "s3" {
    bucket = "hedberg-terraform-states"
    key    = "wr-latest-daily-redirect" // local variable not supported here
    region = "eu-north-1"
  }
}
## API Gateway
resource "aws_apigatewayv2_api" "this" {
  name          = "${local.application_name}-http"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_domain_name" "this" {
  domain_name = local.domain

  domain_name_configuration {
    certificate_arn = "arn:aws:acm:eu-north-1:852264810958:certificate/a73c2e80-7cfa-43f0-a57d-6f2557d3c2ae"
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_stage" "this" {
  api_id      = aws_apigatewayv2_api.this.id
  name        = "$default"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway.arn
    format          = "$context.identity.sourceIp - - [$context.requestTime] \"$context.httpMethod $context.routeKey $context.protocol\" $context.status $context.responseLength $context.requestId $context.integrationErrorMessage"
  }
}

resource "aws_apigatewayv2_api_mapping" "this" {
  api_id      = aws_apigatewayv2_api.this.id
  domain_name = aws_apigatewayv2_domain_name.this.id
  stage       = aws_apigatewayv2_stage.this.id
}

resource "aws_apigatewayv2_integration" "this" {
  api_id           = aws_apigatewayv2_api.this.id
  integration_type = "AWS_PROXY"

  connection_type    = "INTERNET"
  description        = "lambda redirect"
  integration_method = "POST"
  integration_uri    = module.lambda_redirect.lambda_function_invoke_arn
}

resource "aws_apigatewayv2_route" "this" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.this.id}"
}

resource "aws_api_gateway_account" "demo" {
  cloudwatch_role_arn = aws_iam_role.cloudwatch.arn
}

resource "aws_iam_role" "cloudwatch" {
  name = "api_gateway_cloudwatch_global"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "cloudwatch" {
  name = "default"
  role = aws_iam_role.cloudwatch.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:DescribeLogGroups",
                "logs:DescribeLogStreams",
                "logs:PutLogEvents",
                "logs:GetLogEvents",
                "logs:FilterLogEvents"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

## Route 53
resource "aws_route53_record" "this" {
  zone_id = "Z0558704PABIEWFOHFEE"
  name    = "dailytest"
  type    = "CNAME"
  ttl     = 5

  records = [aws_apigatewayv2_domain_name.this.domain_name_configuration[0].target_domain_name]
}


## Lambda
module "lambda_redirect" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "${local.application_name}-redirect"
  description   = "redirect"
  handler       = "wr-latest-daily-redirect"
  runtime       = "go1.x"
  source_path   = "../bin/wr-latest-daily-redirect"
  publish       = true


  environment_variables = {
    REDDIT_CLIENT_ID     = local.secrets.redditClientId,
    REDDIT_CLIENT_SECRET = local.secrets.redditClientSecret,
    USERNAME             = local.secrets.redditUsername,
    APP_ID               = local.secrets.appId,
  }

  allowed_triggers = {
    AllowExecutionFromAPIGateway = {
      service    = "apigateway"
      source_arn = "${aws_apigatewayv2_stage.this.execution_arn}/*"
    }
  }

  cloudwatch_logs_retention_in_days = local.log_retention_in_days

  tags = {
    Name = "wr-latest-daily-redirect"
  }
}


## Cloudwatch
resource "aws_cloudwatch_log_group" "api_gateway" {
  name              = "${local.application_name}-api_gateway"
  retention_in_days = local.log_retention_in_days
}

