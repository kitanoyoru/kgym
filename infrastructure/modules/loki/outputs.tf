output "namespace" {
  description = "Kubernetes namespace where Loki is deployed"
  value       = var.namespace
}

output "release_name" {
  description = "Helm release name"
  value       = var.release_name
}

output "service_name" {
  description = "Name of the Loki service"
  value       = "${var.release_name}"
}

output "service_endpoint" {
  description = "Endpoint for accessing Loki (default port 3100)"
  value       = "${var.release_name}.${var.namespace}.svc.cluster.local:3100"
}
