locals {
  container_image = "${var.quay_repository_path}:latest"
}

module "aws_app_runner" {
  source = "github.com/evoteum/tofu-modules//aws/app_runner"

  aws_apprunner_connection_arn = var.aws_apprunner_connection_arn
  build_command                = "go build -xo app ."
  custom_domain                = var.homepage_url
  dynamodb_arns                = [module.aws_dynamodb.table_arn]
  github_repository            = var.repository_url
  organisation                 = var.organisation
  port                         = var.container_port
  project_name                 = var.repository_name
  runtime = {
    language = var.language
    version  = var.language_version
  }
  source_directory = var.source_path
  start_command    = "./app"

  runtime_environment_variables = {
    AWS_REGION     = var.aws_region
    DYNAMODB_TABLE = module.aws_dynamodb.table_name
  }
}

module "aws_dynamodb" {
  source = "github.com/evoteum/tofu-modules//aws/dynamodb"

  project_name = var.repository_name
  environment  = local.environment
  table_name   = "users"
}
