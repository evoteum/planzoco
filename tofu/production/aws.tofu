locals {
  container_image = "${var.quay_repository_path}:latest"
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
  project_name = var.repository_name
  environment  = local.environment
  table_name   = "all"
}
