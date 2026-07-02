variable "aws_region" {
  type        = string
  description = "The AWS region"
  default     = "eu-west-2"
}

variable "organisation" {
  description = "The name of the organisation"
  type        = string
}

variable "repository_id" {
  description = "The GitHub API ID for this repository"
  type        = string
}

variable "repository_name" {
  type        = string
  description = "The name of this repository (The name of the project). Provided by CI."
}

variable "repository_url" {
  description = "The web address of the git repository"
  type        = string
}
