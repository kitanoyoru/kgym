output "namespace" {
  description = "Kubernetes namespace where Grafana is deployed"
  value       = var.namespace
}

output "release_name" {
  description = "Helm release name"
  value       = var.release_name
}

output "service_name" {
  description = "Name of the Grafana service"
  value       = "${var.release_name}"
}

output "service_endpoint" {
  description = "Endpoint for accessing Grafana (default port 80)"
  value       = "${var.release_name}.${var.namespace}.svc.cluster.local:80"
}

output "admin_user" {
  description = "Grafana admin username"
  value       = var.admin_user
  sensitive   = false
}
