variable "project_name" {
  type        = string
  description = "The name of the project"
  default     = "planzoco"
}

variable "container_image" {
  type        = string
  description = "The container image"
  default     = "evoteum/planzoco:latest"
}

variable "container_port" {
  type        = number
  description = "The container port"
  default     = 8080
}

variable "aws_region" {
  type        = string
  description = "The AWS region"
  default     = "eu-west-2"
}

variable "vpc_cidr" {
  type        = string
  description = "CIDR block for VPC"
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  type        = list(string)
  description = "CIDR blocks for public subnets"
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  type        = list(string)
  description = "CIDR blocks for private subnets"
  default     = ["10.0.3.0/24", "10.0.4.0/24"]
}