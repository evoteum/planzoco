locals {
  container_image = "${var.quay_repository_path}:latest"
}

module "aws_dynamodb" {
  source = "github.com/evoteum/tofu-modules//aws/dynamodb"

  project_name = var.repository_name
  environment  = local.environment
  table_name   = "users"
}

module "aws_app_runner" {
  source = "github.com/evoteum/tofu-modules//aws/app_runner"

  project_name        = var.repository_name
  container_image     = local.container_image
  port                = var.container_port
  custom_domain       = var.homepage_url
  repository_html_url = var.repository_url
  aws_region          = var.aws_region
  organisation = var.organisation
  database_arns = [
    module.aws_dynamodb.table_arn
  ]

  runtime_environment_variables = {
    AWS_REGION     = var.aws_region
    DYNAMODB_TABLE = module.aws_dynamodb.table_name
  }
}
