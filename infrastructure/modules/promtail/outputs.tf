output "namespace" {
  description = "Kubernetes namespace where Promtail is deployed"
  value       = var.namespace
}

output "release_name" {
  description = "Helm release name"
  value       = var.release_name
}

output "service_name" {
  description = "Name of the Promtail service"
  value       = "${var.release_name}"
}

output "service_endpoint" {
  description = "Endpoint for accessing Promtail (default port 3101)"
  value       = "${var.release_name}.${var.namespace}.svc.cluster.local:3101"
}
