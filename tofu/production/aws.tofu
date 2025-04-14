locals {
  container_image = "${var.quay_repository_path}:latest"
}

module "aws_app_runner" {
  source = "github.com/evoteum/tofu-modules//aws/app_runner"

  application_source_uri       = var.repository_url
  aws_apprunner_connection_arn = var.aws_apprunner_connection_arn
  aws_region                   = var.aws_region
  build_command                = "go build -o app ./go/planzoco"
  custom_domain                = var.homepage_url
  organisation                 = var.organisation
  port                         = var.container_port
  project_name                 = var.repository_name
  runtime = {
    language = var.language
    version  = var.language_version
  }
  service_arns = {
    databases = [module.aws_dynamodb.table_arn]
  }
  start_command = "./app"

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
