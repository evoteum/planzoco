terraform {
  required_version = "1.9.0"
  backend "s3" {
    # Values provided by -backend-config in CI
    bucket = ""
    key    = ""
    region = ""
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.84.0"
    }
    cloudflare = {
      source = "cloudflare/cloudflare"
      # Thow shalt not upgrade to v5 lest ye seeketh chaos, confusion and dragons. Ye be warned.
      version = "4.52.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.5.0"
    }
    netparse = {
      source  = "gmeligio/netparse"
      version = "0.0.3"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.6.0"
    }
    quay = {
      source = "registry.terraform.io/enthought/quay"
    }
  }
}

provider "aws" {
  region = var.aws_region
  default_tags {
    tags = {
      Environment   = local.environment
      Project       = var.repository_name
      ManagedBy     = "OpenTofu"
      RepositoryURL = var.repository_url
      RepositoryID  = var.repository_id
    }
  }
}

provider "cloudflare" {
}

provider "github" {
  owner = var.organisation
  # Using `GITHUB_APP_XXX` environment variables in CI
  app_auth {}
}

provider "quay" {
  url = var.quay_url
}
