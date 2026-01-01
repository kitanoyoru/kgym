output "cluster_name" {
  description = "Name of the Redis cluster"
  value       = var.cluster_name
}

output "namespace" {
  description = "Kubernetes namespace where Redis is deployed"
  value       = var.namespace
}

output "service_name" {
  description = "Name of the Redis service"
  value       = var.cluster_name
}

output "service_endpoint" {
  description = "Endpoint for connecting to Redis (default port 6379)"
  value       = "${var.cluster_name}.${var.namespace}.svc.cluster.local:6379"
}

output "replicas" {
  description = "Number of Redis replicas"
  value       = var.replicas
}

output "password_secret_name" {
  description = "Name of the secret containing Redis password (if password is set)"
  value       = var.password != null ? "${var.cluster_name}-password" : null
}
