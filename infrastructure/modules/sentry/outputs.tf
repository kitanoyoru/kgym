output "namespace" {
  description = "Kubernetes namespace where Sentry is deployed"
  value       = var.namespace
}

output "release_name" {
  description = "Helm release name"
  value       = var.release_name
}

output "service_name" {
  description = "Name of the Sentry web service"
  value       = "${var.release_name}-web"
}

output "service_endpoint" {
  description = "Endpoint for accessing Sentry web UI (default port 9000)"
  value       = "${var.release_name}-web.${var.namespace}.svc.cluster.local:9000"
}

output "user_email" {
  description = "Email for the initial Sentry superuser"
  value       = var.user_email
  sensitive   = true
}
