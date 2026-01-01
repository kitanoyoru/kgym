output "cluster_name" {
  description = "Name of the MinIO cluster"
  value       = var.cluster_name
}

output "namespace" {
  description = "Kubernetes namespace where MinIO is deployed"
  value       = var.namespace
}

output "service_name" {
  description = "Name of the MinIO API service"
  value       = "${var.cluster_name}-api"
}

output "service_endpoint" {
  description = "Endpoint for accessing MinIO API (default port 9000)"
  value       = "${var.cluster_name}-api.${var.namespace}.svc.cluster.local:9000"
}

output "console_endpoint" {
  description = "Endpoint for accessing MinIO Console (default port 9001)"
  value       = "${var.cluster_name}-console.${var.namespace}.svc.cluster.local:9001"
}

output "access_key_secret_name" {
  description = "Name of the secret containing MinIO credentials"
  value       = "${var.cluster_name}-credentials"
}

output "replicas" {
  description = "Number of MinIO replicas"
  value       = var.replicas
}
