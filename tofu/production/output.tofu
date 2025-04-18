output "environment" {
  value = local.environment
}

output "runtime" {
  value = "${var.language}:${var.language_version}"
}