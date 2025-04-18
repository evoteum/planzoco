output "environment" {
  value = local.environment
}

output "runtime" {
  value = module.aws_app_runner.runtime
}

output "service_url" {
  value = module.aws_app_runner.service_url
}