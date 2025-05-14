module "platform_service" {
  source = "github.com/evoteum/tofu-modules//platform/service"

  alb_dns_name          = var.production_alb_dns_name
  security_group_id     = var.production_alb_security_group_id
  cloudflare_account_id = var.cloudflare_account_id
  health_check_path     = "/health"

  cluster_arn = var.production_ecs_cluster_arn
  dynamodb = {
    table_name = "all"
    hash_key   = "pk"
    range_key  = "sk"
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
  }
  environment           = local.environment
  homepage_url          = var.homepage_url
  launch_type           = var.launch_type
  listener_arn          = var.production_alb_https_listener_arn
  log_retention_in_days = var.log_retention_in_days
  memory                = var.memory
  network_mode          = var.network_mode
  organisation          = var.organisation
  priority              = var.priority

  region                 = var.aws_region
  repository_archived    = var.repository_archived
  repository_description = var.repository_description
  repository_id          = var.repository_id
  repository_name        = var.repository_name
  repository_url         = var.repository_url
  repository_visibility  = var.repository_visibility
  vpc_id                 = var.production_vpc_id
}