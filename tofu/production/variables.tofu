variable "aws_region" {
  type        = string
  description = "The AWS region"
  default     = "eu-west-2"
}

variable "container_port" {
  type        = number
  description = "The container port"
  default     = 8080
}

variable "homepage_url" {
  description = "The web address of publicly available application"
  type        = string
}

variable "private_subnet_cidrs" {
  type        = list(string)
  description = "CIDR blocks for private subnets"
  default     = ["10.0.3.0/24", "10.0.4.0/24"]
}

variable "public_subnet_cidrs" {
  type        = list(string)
  description = "CIDR blocks for public subnets"
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "quay_repository_path" {
  type        = string
  description = "Path to the quay repository. Provided by CI."
}

variable "repository_name" {
  type        = string
  description = "The name of this repository (The name of the project). Provided by CI."
}

variable "vpc_cidr" {
  type        = string
  description = "CIDR block for VPC"
  default     = "10.0.0.0/16"
}

variable "repository_url" {
  description = "The web address of the git repository"
  type        = string
}

variable "organisation" {
  description = "The name of the organisation"
  type = string
}
