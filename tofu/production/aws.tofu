locals {
  container_image = "${var.quay_repository_path}:latest"
  environment     = basename(path.cwd)
}

module "aws_networking" {
  source = "github.com/evoteum/tofu-modules//aws/networking"

  aws_region           = var.aws_region
  container_port       = var.container_port
  private_subnet_cidrs = var.private_subnet_cidrs
  project_name         = var.repository_name
  public_subnet_cidrs  = var.public_subnet_cidrs
  vpc_cidr             = var.vpc_cidr
}

module "aws_dynamodb" {
  source = "github.com/evoteum/tofu-modules//aws/dynamodb"

  project_name = var.repository_name
  environment  = local.environment
}

module "aws_app_runner" {
  source = "github.com/evoteum/tofu-modules//aws/app_runner"

  project_name    = var.repository_name
  container_image = local.container_image
  custom_domain   = var.homepage_url
  aws_region      = var.aws_region

  runtime_environment_variables = {
    AWS_REGION     = var.aws_region
    DYNAMODB_TABLE = module.dynamodb.table_name
  }
}
