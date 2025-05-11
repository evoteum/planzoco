variable "aws_apprunner_connection_arn" {
  description = "Amazon Resource Name of the estate aws_apprunner_connection"
  type        = string
  sensitive   = true
}

variable "aws_region" {
  type        = string
  description = "The AWS region"
  default     = "eu-west-2"
}

variable "cloudflare_account_id" {
  description = "ID of the cloudflare account in use"
  type        = string
}

variable "container_port" {
  type        = number
  description = "The container port"
  default     = 8080
}

variable "cpu" {
  description = "Number of CPU units used by the task (e.g. 256, 512)"
  type        = number
  default     = 256
}

variable "desired_count" {
  description = "Number of desired running tasks"
  type        = number
  default     = 1
}

variable "homepage_url" {
  description = "The web address of publicly available application"
  type        = string
}

variable "language" {
  type        = string
  description = "The language that this project is written in"
}

variable "language_version" {
  type        = string
  description = "The version of the language that this project is written in"
}

variable "launch_type" {
  default     = "FARGATE"
  type        = string
  description = "Launch type on which to run your service. The valid values are EC2, FARGATE, and EXTERNAL."
}

variable "log_retention_in_days" {
  type        = number
  description = "Number of days to keep logs"
  default     = 3
}

variable "memory" {
  description = "Amount of memory (in MiB) used by the task (e.g. 512, 1024)"
  type        = number
  default     = 1024
}

variable "network_mode" {
  default     = "awsvpc"
  description = "Docker networking mode to use for the containers in the task. Valid values are none, bridge, awsvpc, and host."
  type        = string
}

variable "organisation" {
  description = "The name of the organisation"
  type        = string
}

variable "priority" {
  description = "Manual priority of the listener rule (must be unique per listener). If not set, a priority will be automatically assigned."
  type        = number
  default     = null
}

variable "private_subnet_cidrs" {
  type        = list(string)
  description = "CIDR blocks for private subnets"
  default     = ["10.0.3.0/24", "10.0.4.0/24"]
}

variable "production_alb_dns_name" {
  type        = string
  description = "ALB dns name for production"
}

variable "production_alb_https_listener_arn" {
  description = "ARN of the production listener"
  type        = string
}

variable "production_alb_security_group_id" {
  type        = string
  description = "ALB security group ID for production"
}

variable "production_production_ecs_cluster_arn" {
  description = "ARN of the ECS cluster to deploy into"
  type        = string
}

variable "production_vpc_id" {
  description = "ID of the VPC that this service will be installed in"
  type        = string
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

variable "region" {
  description = "AWS region"
  type        = string
  default     = "eu-west-2"
}

variable "repository_archived" {
  description = "Whether the repository is archived"
  type        = bool
}

variable "repository_description" {
  description = "The description of this git repository"
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

variable "repository_visibility" {
  description = "Whether this repository is public or private"
  type        = string
}

variable "source_path" {
  description = "The filepath to the source code from the repo root"
  type        = string
}

variable "vpc_cidr" {
  type        = string
  description = "CIDR block for VPC"
  default     = "10.0.0.0/16"
}
