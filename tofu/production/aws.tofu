locals {
  container_image = "${var.quay_repository_path}:latest"
  environment     = basename(path.cwd)
}

module "aws_dynamodb" {
  source = "github.com/evoteum/tofu-modules//aws/dynamodb"

  project_name = var.repository_name
  environment  = local.environment
  attributes = [
    {
      name = "event_id"
      type = "S"
    },
    {
      name = "user_id"
      type = "s"
    }
  ]
}

module "aws_app_runner" {
  source = "github.com/evoteum/tofu-modules//aws/app_runner"

  project_name    = var.repository_name
  container_image = local.container_image
  port            = var.container_port
  custom_domain   = var.homepage_url
  aws_region      = var.aws_region

  runtime_environment_variables = {
    AWS_REGION     = var.aws_region
    DYNAMODB_TABLE = module.dynamodb.table_name
  }
}
