terraform {
    required_providers {
        quay = {
        source = "enthought/quay"
        version = "0.5.1"
        }
        aws = {
        source  = "hashicorp/aws"
        version = "~> 5.0"
        }
    }
  required_version = "1.9.0"
  backend "s3" {
    bucket = ""
    key    = ""
    region = ""
  }
}

module "common" {
  source = "github.com:evoteum/tofu-modules.git//common/vars?ref=main"
}

module "networking" {
  source = "github.com:evoteum/tofu-modules//aws/networking"

  project_name         = var.project_name
  vpc_cidr            = var.vpc_cidr
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  container_port      = var.container_port
}

module "dynamodb" {
  source = "github.com:evoteum/tofu-modules//aws/dynamodb"

  project_name = var.project_name
  environment  = local.environment
}

module "app_runner" {
  source = "github.com:evoteum/tofu-modules//aws/app_runner"

  project_name    = var.project_name
  environment     = local.environment
  container_image = var.container_image
  
  runtime_environment_variables = {
    AWS_REGION      = var.aws_region
    DYNAMODB_TABLE  = module.dynamodb.table_name
  }
  
}

module "quay" {
  source = "github.com:evoteum/tofu-modules//quay/repo"

  project_name = var.project_name
}

