output "environment" {
  value = local.environment
}

output "runtime" {
  value = module.aws_app_runner.runtime
}