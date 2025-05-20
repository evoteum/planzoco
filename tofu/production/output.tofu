output "environment" {
  value = local.environment
}

output "listener_arn" {
  value = module.platform_service.listener_rule_arn
}

output "service_arn" {
  value = module.platform_service.service_arn
}
