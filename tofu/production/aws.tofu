locals {
  container_image = "${var.quay_repository_path}:latest"
  environment     = basename(path.cwd)
}

module "networking" {
  source = "github.com/evoteum/tofu-modules//aws/networking"

  project_name         = var.repository_name
  vpc_cidr             = var.vpc_cidr
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  container_port       = var.container_port
}

module "dynamodb" {
  source = "github.com/evoteum/tofu-modules//aws/dynamodb"

  project_name = var.repository_name
  environment  = local.environment
}

module "app_runner" {
  source = "github.com/evoteum/tofu-modules//aws/app_runner"

  project_name    = var.repository_name
  container_image = local.container_image

  runtime_environment_variables = {
    AWS_REGION     = var.aws_region
    DYNAMODB_TABLE = module.dynamodb.table_name
  }
}
