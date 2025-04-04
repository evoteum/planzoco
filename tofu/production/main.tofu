terraform {
  required_providers {
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

provider "aws" {
  region = var.aws_region
}

module "common" {
  source = "github.com:evoteum/tofu-modules.git//common/vars?ref=main"
}

