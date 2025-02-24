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
}

provider "aws" {
  region = var.aws_region
} 

provider "quay" {
  url = "https://quay.example.com"
}