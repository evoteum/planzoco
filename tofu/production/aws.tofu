locals {
  container_image = "${var.quay_repository_path}:latest"
}

module "aws_app_runner" {
  source = "github.com/evoteum/tofu-modules//aws/app_runner"

  aws_apprunner_connection_arn = var.aws_apprunner_connection_arn
  build_command                = "go build main.go"
  dynamodb_arns                = [module.aws_dynamodb.table_arn]
  github_repository            = var.repository_url
  organisation                 = var.organisation
  port                         = var.container_port
  project_name                 = var.repository_name
  runtime = {
    language = var.language
    version  = var.language_version
  }
  source_directory = var.source_path
  start_command    = "./main"

  runtime_environment_variables = {
    AWS_REGION     = var.aws_region
    DYNAMODB_TABLE = module.aws_dynamodb.table_name
  }
}

module "aws_dynamodb" {
  source = "github.com/evoteum/tofu-modules//aws/dynamodb"

  attributes = {
    pk          = "s"
    sk          = "s"
    entity_type = "s"
    event_id    = "s"
    question_id = "s"
  }
  global_secondary_indexes = {
    EntityTypeIndex = {
      hash_key        = "entity_type"
      range_key       = "pk"
      projection_type = "all"
    }
    EventIDIndex = {
      hash_key        = "event_id"
      range_key       = "pk"
      projection_type = "all"
    }
    QuestionIDIndex = {
      hash_key        = "question_id"
      range_key       = "pk"
      projection_type = "all"
    }
  }
  environment  = local.environment
  hash_key     = "pk"
  project_name = var.repository_name
  range_key    = "sk"
  table_name   = "all"
}
